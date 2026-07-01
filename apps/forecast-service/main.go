package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Forecast struct {
	Location         string  `json:"location"`
	Region           string  `json:"region"`
	Date             string  `json:"date"`
	Condition        string  `json:"condition"`
	TemperatureC     int     `json:"temperatureC"`
	WindKph          int     `json:"windKph"`
	RainChance       int     `json:"rainChance"`
	MarineHeatRisk   string  `json:"marineHeatRisk"`
	CoastalFloodRisk string  `json:"coastalFloodRisk"`
	RiskScore        float64 `json:"riskScore"`
	Summary          string  `json:"summary"`
	PremiumAdvisory  string  `json:"premiumAdvisory,omitempty"`
}

type PageData struct {
	Title    string
	Forecast *Forecast
	Error    string
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/forecast", handleForecastAPI)
	mux.HandleFunc("/api/premium/forecast", handlePremiumForecastAPI)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("AMOC Sentinel forecast service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, logRequest(mux)); err != nil {
		log.Fatal(err)
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
	writeJSON(w, forecast)
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
