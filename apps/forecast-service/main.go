package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Forecast struct {
	Location             string                 `json:"location"`
	Region               string                 `json:"region"`
	Date                 string                 `json:"date"`
	Condition            string                 `json:"condition"`
	TemperatureC         int                    `json:"temperatureC"`
	WindKph              int                    `json:"windKph"`
	RainChance           int                    `json:"rainChance"`
	MarineHeatRisk       string                 `json:"marineHeatRisk"`
	CoastalFloodRisk     string                 `json:"coastalFloodRisk"`
	RiskScore            float64                `json:"riskScore"`
	Summary              string                 `json:"summary"`
	PremiumAdvisory      string                 `json:"premiumAdvisory,omitempty"`
	RememberedContext    *OperatorContextMemory `json:"rememberedContext,omitempty"`
	MemoryImplementation string                 `json:"memoryImplementation,omitempty"`
	Latitude             float64                `json:"latitude,omitempty"`
	Longitude            float64                `json:"longitude,omitempty"`
	CoordinatesAvailable bool                   `json:"coordinatesAvailable"`
	CoordinateLabel      string                 `json:"coordinateLabel,omitempty"`
}

type PageData struct {
	Title    string
	Forecast *Forecast
	Error    string
}

type TernaryFlag string

const (
	FlagNormal   TernaryFlag = "normal"
	FlagInfinity TernaryFlag = "infinity"
	FlagNaN      TernaryFlag = "nan"
	setunMax     int         = 19682
	setunBias    int         = 9841
	setunScale   float64     = 19683.0
)

type TritEngine9 struct {
	Bits uint16 `json:"bits"`
}

type SetunPayload struct {
	Input     float64     `json:"input,omitempty"`
	Bits      uint16      `json:"bits"`
	Hex       string      `json:"hex"`
	Flag      TernaryFlag `json:"flag"`
	RawValue  uint16      `json:"rawValue"`
	Value     float64     `json:"value"`
	Trits     []int       `json:"trits"`
	TritLabel string      `json:"tritLabel"`
}

type SetunAddResponse struct {
	A      SetunPayload `json:"a"`
	B      SetunPayload `json:"b"`
	Result SetunPayload `json:"result"`
}

var (
	tmpl                          = template.Must(template.ParseFiles("templates/index.html"))
	aetherBioHawkTmpl             = template.Must(template.ParseFiles("templates/aetherbiohawk.html"))
	memoryStore       MemoryStore = NewInMemoryStore()
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/aetherbiohawk", handleAetherBioHawk)
	mux.HandleFunc("/api/forecast", handleForecastAPI)
	mux.HandleFunc("/api/premium/forecast", handlePremiumForecastAPI)
	mux.HandleFunc("/api/memory/operator", handleOperatorMemoryAPI)
	mux.HandleFunc("/api/setun/encode", handleSetunEncodeAPI)
	mux.HandleFunc("/api/setun/add", handleSetunAddAPI)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("AMOC Sentinel forecast service listening on :%s", port)
	log.Printf("Memory interface enabled with placeholder implementation: %T", memoryStore)
	if err := http.ListenAndServe(":"+port, logRequest(mux)); err != nil {
		log.Fatal(err)
	}
}

func handleAetherBioHawk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := aetherBioHawkTmpl.Execute(w, PageData{Title: "AetherBioHawk Concept Demo"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderPage(w, PageData{Title: "AMOC Sentinel Forecast Service"})
		return
	}

	if err := r.ParseForm(); err != nil {
		renderPage(w, PageData{Title: "AMOC Sentinel Forecast Service", Error: "Could not parse form."})
		return
	}

	location := strings.TrimSpace(r.FormValue("location"))
	region := strings.TrimSpace(r.FormValue("region"))
	date := strings.TrimSpace(r.FormValue("date"))
	if location == "" || date == "" {
		renderPage(w, PageData{Title: "AMOC Sentinel Forecast Service", Error: "Location and date are required."})
		return
	}

	forecast := generateForecast(location, region, date, false)
	hydrateForecastCoordinates(&forecast)
	hydrateForecastMemory(r, &forecast)
	renderPage(w, PageData{Title: "AMOC Sentinel Forecast Service", Forecast: &forecast})
}

func handleForecastAPI(w http.ResponseWriter, r *http.Request) {
	location := strings.TrimSpace(r.URL.Query().Get("location"))
	region := strings.TrimSpace(r.URL.Query().Get("region"))
	date := strings.TrimSpace(r.URL.Query().Get("date"))

	if location == "" || date == "" {
		http.Error(w, `{"error":"location and date are required"}`, http.StatusBadRequest)
		return
	}

	forecast := generateForecast(location, region, date, false)
	hydrateForecastCoordinates(&forecast)
	hydrateForecastMemory(r, &forecast)
	writeJSON(w, forecast)
}

func handlePremiumForecastAPI(w http.ResponseWriter, r *http.Request) {
	paymentHeader := strings.TrimSpace(r.Header.Get("X-X402-Payment"))
	if paymentHeader == "" {
		w.WriteHeader(http.StatusPaymentRequired)
		writeJSON(w, map[string]any{
			"error": "payment required",
			"hint":  "Send X-X402-Payment header to simulate premium access.",
			"price": "0.01 premium credit",
		})
		return
	}

	location := strings.TrimSpace(r.URL.Query().Get("location"))
	region := strings.TrimSpace(r.URL.Query().Get("region"))
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	if location == "" || date == "" {
		http.Error(w, `{"error":"location and date are required"}`, http.StatusBadRequest)
		return
	}

	forecast := generateForecast(location, region, date, true)
	hydrateForecastCoordinates(&forecast)
	hydrateForecastMemory(r, &forecast)
	writeJSON(w, forecast)
}

func handleOperatorMemoryAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		location := strings.TrimSpace(r.URL.Query().Get("location"))
		region := strings.TrimSpace(r.URL.Query().Get("region"))
		if location == "" {
			http.Error(w, `{"error":"location is required"}`, http.StatusBadRequest)
			return
		}

		mem, err := memoryStore.GetOperatorContext(r.Context(), memoryKey(location, region))
		if err != nil {
			http.Error(w, `{"error":"could not load memory"}`, http.StatusInternalServerError)
			return
		}
		if mem == nil {
			writeJSON(w, map[string]any{"memory": nil, "implementation": fmt.Sprintf("%T", memoryStore)})
			return
		}
		writeJSON(w, map[string]any{"memory": mem, "implementation": fmt.Sprintf("%T", memoryStore)})
	case http.MethodPost:
		var payload OperatorContextMemory
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}
		payload.Location = strings.TrimSpace(payload.Location)
		payload.Region = strings.TrimSpace(payload.Region)
		if payload.Key == "" {
			payload.Key = memoryKey(payload.Location, payload.Region)
		}
		if payload.Key == "" {
			http.Error(w, `{"error":"location or key is required"}`, http.StatusBadRequest)
			return
		}
		if err := memoryStore.UpsertOperatorContext(r.Context(), payload); err != nil {
			http.Error(w, `{"error":"could not save memory"}`, http.StatusInternalServerError)
			return
		}
		mem, _ := memoryStore.GetOperatorContext(r.Context(), payload.Key)
		writeJSON(w, map[string]any{"memory": mem, "implementation": fmt.Sprintf("%T", memoryStore)})
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSetunEncodeAPI(w http.ResponseWriter, r *http.Request) {
	value, err := parseFloatParam(r, "value")
	if err != nil {
		http.Error(w, `{"error":"value query parameter is required and must be numeric"}`, http.StatusBadRequest)
		return
	}
	engine := setunFromFloat(value)
	writeJSON(w, buildSetunPayload(engine, value))
}

func handleSetunAddAPI(w http.ResponseWriter, r *http.Request) {
	aValue, err := parseFloatParam(r, "a")
	if err != nil {
		http.Error(w, `{"error":"a query parameter is required and must be numeric"}`, http.StatusBadRequest)
		return
	}
	bValue, err := parseFloatParam(r, "b")
	if err != nil {
		http.Error(w, `{"error":"b query parameter is required and must be numeric"}`, http.StatusBadRequest)
		return
	}
	a := setunFromFloat(aValue)
	b := setunFromFloat(bValue)
	result := a.Add(b)
	writeJSON(w, SetunAddResponse{
		A:      buildSetunPayload(a, aValue),
		B:      buildSetunPayload(b, bValue),
		Result: buildSetunPayload(result, result.ToFloat()),
	})
}

func parseFloatParam(r *http.Request, key string) (float64, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return 0, fmt.Errorf("missing value")
	}
	return strconv.ParseFloat(value, 64)
}

func buildSetunPayload(engine TritEngine9, input float64) SetunPayload {
	trits := engine.ToTrits()
	tritInts := make([]int, len(trits))
	parts := make([]string, len(trits))
	for i, t := range trits {
		tritInts[i] = t
		parts[i] = strconv.Itoa(t)
	}
	return SetunPayload{
		Input:     input,
		Bits:      engine.Bits,
		Hex:       fmt.Sprintf("0x%04X", engine.Bits),
		Flag:      engine.Flag(),
		RawValue:  engine.RawValue(),
		Value:     engine.ToFloat(),
		Trits:     tritInts,
		TritLabel: strings.Join(parts, " "),
	}
}

func generateForecast(location, region, date string, premium bool) Forecast {
	seed := int64(len(location) + len(region) + len(date))
	for _, c := range location + region + date {
		seed += int64(c)
	}
	r := rand.New(rand.NewSource(seed))

	conditions := []string{"Partly cloudy", "Scattered showers", "Sunny intervals", "Windy with showers", "Humid with thunderstorms"}
	marineRisks := []string{"Low", "Moderate", "Elevated"}
	coastalRisks := []string{"Low", "Moderate", "Elevated"}

	temp := 27 + r.Intn(7)
	wind := 12 + r.Intn(28)
	rain := 10 + r.Intn(80)
	marine := marineRisks[r.Intn(len(marineRisks))]
	coastal := coastalRisks[r.Intn(len(coastalRisks))]
	riskScore := float64((rain+wind)/2) / 10.0
	if marine == "Elevated" {
		riskScore += 1.2
	}
	if coastal == "Elevated" {
		riskScore += 1.2
	}

	summary := fmt.Sprintf("%s on %s is expected to see %s conditions with %d%% rain chance and winds near %d kph.", location, date, strings.ToLower(conditions[r.Intn(len(conditions))]), rain, wind)

	forecast := Forecast{
		Location:         location,
		Region:           region,
		Date:             date,
		Condition:        conditions[r.Intn(len(conditions))],
		TemperatureC:     temp,
		WindKph:          wind,
		RainChance:       rain,
		MarineHeatRisk:   marine,
		CoastalFloodRisk: coastal,
		RiskScore:        round(riskScore, 1),
		Summary:          summary,
	}

	if premium {
		forecast.PremiumAdvisory = premiumAdvice(forecast)
	}

	return forecast
}

func hydrateForecastMemory(r *http.Request, forecast *Forecast) {
	key := memoryKey(forecast.Location, forecast.Region)
	forecast.MemoryImplementation = fmt.Sprintf("%T", memoryStore)

	if mem, err := memoryStore.GetOperatorContext(r.Context(), key); err == nil && mem != nil {
		forecast.RememberedContext = mem
	}

	summary := ForecastMemorySummary{
		Location: forecast.Location,
		Region:   forecast.Region,
		Date:     forecast.Date,
		Summary:  forecast.Summary,
	}
	if err := memoryStore.RecordForecastSummary(r.Context(), key, summary); err != nil {
		log.Printf("memory record failed for %s: %v", key, err)
		return
	}

	if forecast.RememberedContext == nil {
		if mem, err := memoryStore.GetOperatorContext(r.Context(), key); err == nil && mem != nil {
			forecast.RememberedContext = mem
		}
	}
}

func premiumAdvice(f Forecast) string {
	var notes []string
	if f.CoastalFloodRisk == "Elevated" {
		notes = append(notes, "Review coastal drainage, move exposed materials higher, and prepare flood barriers.")
	}
	if f.MarineHeatRisk == "Elevated" {
		notes = append(notes, "Monitor marine heat stress impacts and temperature-sensitive coastal operations.")
	}
	if f.WindKph >= 30 {
		notes = append(notes, "Secure loose exterior assets and check backup power readiness.")
	}
	if len(notes) == 0 {
		notes = append(notes, "Maintain normal monitoring cadence and review the next forecast cycle.")
	}
	return strings.Join(notes, " ")
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func renderPage(w http.ResponseWriter, data PageData) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func round(v float64, places int) float64 {
	factor := mathPow10(places)
	return float64(int(v*factor+0.5)) / factor
}

func mathPow10(p int) float64 {
	return float64(pow10Int(p))
}

func pow10Int(p int) int {
	result := 1
	for i := 0; i < p; i++ {
		result *= 10
	}
	return result
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, strconv.FormatInt(time.Since(start).Milliseconds(), 10)+"ms")
	})
}

func (t TritEngine9) Flag() TernaryFlag {
	switch (t.Bits >> 14) & 0x03 {
	case 0:
		return FlagNormal
	case 1:
		return FlagInfinity
	case 2:
		return FlagNaN
	default:
		return FlagNaN
	}
}

func (t TritEngine9) RawValue() uint16 {
	return t.Bits & 0x3FFF
}

func setunFromFloat(val float64) TritEngine9 {
	if math.IsNaN(val) {
		return TritEngine9{Bits: 2 << 14}
	}
	if math.IsInf(val, 0) {
		return TritEngine9{Bits: 1 << 14}
	}
	bounded := math.Max(-0.5, math.Min(0.5, val))
	scaled := bounded * setunScale
	balancedInt := int(math.Round(scaled))
	unsignedVal := balancedInt + setunBias
	if unsignedVal < 0 {
		unsignedVal = 0
	}
	if unsignedVal > setunMax {
		unsignedVal = setunMax
	}
	return TritEngine9{Bits: uint16(unsignedVal)}
}

func (t TritEngine9) ToFloat() float64 {
	switch t.Flag() {
	case FlagNaN:
		return math.NaN()
	case FlagInfinity:
		return math.Inf(1)
	default:
		balancedInt := int(t.RawValue()) - setunBias
		return float64(balancedInt) / setunScale
	}
}

func (t TritEngine9) Add(other TritEngine9) TritEngine9 {
	if t.Flag() == FlagNaN || other.Flag() == FlagNaN {
		return TritEngine9{Bits: 2 << 14}
	}
	if t.Flag() == FlagInfinity || other.Flag() == FlagInfinity {
		return TritEngine9{Bits: 1 << 14}
	}
	val1 := int(t.RawValue()) - setunBias
	val2 := int(other.RawValue()) - setunBias
	result := val1 + val2
	if result < -setunBias {
		result = -setunBias
	}
	if result > setunBias {
		result = setunBias
	}
	return TritEngine9{Bits: uint16(result + setunBias)}
}

func (t TritEngine9) ToTrits() []int {
	trits := make([]int, 9)
	if t.Flag() != FlagNormal {
		return trits
	}
	rem := int(t.RawValue())
	for i := 8; i >= 0; i-- {
		remainder := rem % 3
		rem /= 3
		switch remainder {
		case 0:
			trits[i] = 0
		case 1:
			trits[i] = 1
		case 2:
			trits[i] = -1
			rem += 1
		}
	}
	return trits
}
