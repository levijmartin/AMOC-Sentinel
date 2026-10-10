package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultArchiveWeatherEndpoint  = "https://archive-api.open-meteo.com/v1/archive"
	defaultForecastWeatherEndpoint = "https://api.open-meteo.com/v1/forecast"
)

var (
	archiveWeatherEndpoint  = defaultArchiveWeatherEndpoint
	forecastWeatherEndpoint = defaultForecastWeatherEndpoint
	weatherClient           = &http.Client{Timeout: 12 * time.Second}
)

type openMeteoDaily struct {
	Time                        []string  `json:"time"`
	WeatherCode                 []int     `json:"weather_code"`
	Temperature2MMax            []float64 `json:"temperature_2m_max"`
	Temperature2MMin            []float64 `json:"temperature_2m_min"`
	PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
	PrecipitationSum            []float64 `json:"precipitation_sum"`
	WindSpeed10MMax             []float64 `json:"wind_speed_10m_max"`
}

type openMeteoResponse struct {
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	Timezone  string         `json:"timezone"`
	Daily     openMeteoDaily `json:"daily"`
	Reason    string         `json:"reason"`
	Error     bool           `json:"error"`
}

func buildWeatherForecast(ctx context.Context, location, region, date string, premium bool) (Forecast, error) {
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	if err != nil {
		return Forecast{}, fmt.Errorf("date must use YYYY-MM-DD")
	}

	lookup, ok := geocodeLocation(location, region)
	if !ok {
		return Forecast{}, fmt.Errorf("location could not be resolved")
	}

	weather, kind, err := fetchOpenMeteoDaily(ctx, lookup.Latitude, lookup.Longitude, parsedDate)
	if err != nil {
		return Forecast{}, err
	}

	condition := weatherCodeLabel(weather.WeatherCode)
	temperature := int((weather.TemperatureMax + weather.TemperatureMin) / 2.0)
	coastalRisk := coastalScreeningRisk(weather.PrecipitationMM, weather.WindKph)
	forecast := Forecast{
		Location:             lookup.Label,
		Region:               "",
		Date:                 parsedDate.Format("2006-01-02"),
		Condition:            condition,
		TemperatureC:         temperature,
		TemperatureMinC:      weather.TemperatureMin,
		TemperatureMaxC:      weather.TemperatureMax,
		WindKph:              int(weather.WindKph + 0.5),
		RainChance:           int(weather.RainChance + 0.5),
		PrecipitationMM:      weather.PrecipitationMM,
		MarineHeatRisk:       "Not assessed",
		CoastalFloodRisk:     coastalRisk,
		RiskScore:            weatherRiskScore(weather.PrecipitationMM, weather.WindKph),
		Summary:              weatherSummary(lookup.Label, parsedDate, kind, condition, weather),
		Latitude:             lookup.Latitude,
		Longitude:            lookup.Longitude,
		CoordinatesAvailable: true,
		CoordinateLabel:      lookup.Label,
		WeatherSource:        "Open-Meteo",
		WeatherSourceURL:     "https://open-meteo.com/",
		WeatherDataKind:      kind,
		WeatherTimezone:      weather.Timezone,
	}
	if premium {
		forecast.PremiumAdvisory = premiumAdvice(forecast)
	}
	return forecast, nil
}

type normalizedDailyWeather struct {
	WeatherCode     int
	TemperatureMax  float64
	TemperatureMin  float64
	PrecipitationMM float64
	RainChance      float64
	WindKph         float64
	Timezone        string
}

func fetchOpenMeteoDaily(ctx context.Context, latitude, longitude float64, date time.Time) (normalizedDailyWeather, string, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -5)
	endpoint := forecastWeatherEndpoint
	kind := "forecast"
	dailyVariables := "weather_code,temperature_2m_max,temperature_2m_min,precipitation_probability_max,precipitation_sum,wind_speed_10m_max"
	if date.Before(time.Date(cutoff.Year(), cutoff.Month(), cutoff.Day(), 0, 0, 0, 0, time.UTC)) {
		endpoint = archiveWeatherEndpoint
		kind = "historical"
		dailyVariables = "weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum,wind_speed_10m_max"
	}

	query := url.Values{}
	query.Set("latitude", strconv.FormatFloat(latitude, 'f', 5, 64))
	query.Set("longitude", strconv.FormatFloat(longitude, 'f', 5, 64))
	query.Set("start_date", date.Format("2006-01-02"))
	query.Set("end_date", date.Format("2006-01-02"))
	query.Set("daily", dailyVariables)
	query.Set("timezone", "auto")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+query.Encode(), nil)
	if err != nil {
		return normalizedDailyWeather{}, kind, fmt.Errorf("weather request could not be created")
	}
	req.Header.Set("User-Agent", "AMOC-Sentinel/1.0 (canonical weather query)")
	resp, err := weatherClient.Do(req)
	if err != nil {
		return normalizedDailyWeather{}, kind, fmt.Errorf("Open-Meteo request failed: %w", err)
	}
	defer resp.Body.Close()

	var payload openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return normalizedDailyWeather{}, kind, fmt.Errorf("Open-Meteo response could not be decoded")
	}
	if resp.StatusCode != http.StatusOK || payload.Error {
		if payload.Reason == "" {
			payload.Reason = resp.Status
		}
		return normalizedDailyWeather{}, kind, fmt.Errorf("Open-Meteo rejected the query: %s", payload.Reason)
	}
	if len(payload.Daily.Time) != 1 || len(payload.Daily.WeatherCode) != 1 || len(payload.Daily.Temperature2MMax) != 1 || len(payload.Daily.Temperature2MMin) != 1 || len(payload.Daily.PrecipitationSum) != 1 || len(payload.Daily.WindSpeed10MMax) != 1 {
		return normalizedDailyWeather{}, kind, fmt.Errorf("Open-Meteo returned no daily record for %s", date.Format("2006-01-02"))
	}

	result := normalizedDailyWeather{
		WeatherCode:     payload.Daily.WeatherCode[0],
		TemperatureMax:  payload.Daily.Temperature2MMax[0],
		TemperatureMin:  payload.Daily.Temperature2MMin[0],
		PrecipitationMM: payload.Daily.PrecipitationSum[0],
		WindKph:         payload.Daily.WindSpeed10MMax[0],
		Timezone:        payload.Timezone,
	}
	if len(payload.Daily.PrecipitationProbabilityMax) == 1 {
		result.RainChance = payload.Daily.PrecipitationProbabilityMax[0]
	}
	return result, kind, nil
}

func weatherCodeLabel(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code == 1:
		return "Mainly clear"
	case code == 2:
		return "Partly cloudy"
	case code == 3:
		return "Overcast"
	case code == 45 || code == 48:
		return "Fog"
	case code >= 51 && code <= 57:
		return "Drizzle"
	case code >= 61 && code <= 67:
		return "Rain"
	case code >= 71 && code <= 77:
		return "Snow"
	case code >= 80 && code <= 82:
		return "Rain showers"
	case code >= 85 && code <= 86:
		return "Snow showers"
	case code >= 95:
		return "Thunderstorm"
	default:
		return "Unclassified weather"
	}
}

func coastalScreeningRisk(precipitationMM, windKph float64) string {
	switch {
	case precipitationMM >= 50 || windKph >= 65:
		return "Elevated"
	case precipitationMM >= 20 || windKph >= 40:
		return "Moderate"
	default:
		return "Low"
	}
}

func weatherRiskScore(precipitationMM, windKph float64) float64 {
	score := precipitationMM/20 + windKph/25
	if score > 10 {
		score = 10
	}
	return round(score, 1)
}

func weatherSummary(location string, date time.Time, kind, condition string, weather normalizedDailyWeather) string {
	label := "forecast"
	if kind == "historical" {
		label = "historical record"
	}
	return fmt.Sprintf("Open-Meteo %s for %s on %s: %s, %.1f–%.1f°C, %.1f mm precipitation, and maximum wind %.1f kph.", label, location, date.Format("2006-01-02"), strings.ToLower(condition), weather.TemperatureMin, weather.TemperatureMax, weather.PrecipitationMM, weather.WindKph)
}
