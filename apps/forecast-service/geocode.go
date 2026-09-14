package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const defaultGeocodingEndpoint = "https://geocoding-api.open-meteo.com/v1/search"

type geoResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Admin1    string  `json:"admin1"`
	Admin2    string  `json:"admin2"`
}

type geoResponse struct {
	Results []geoResult `json:"results"`
}

type coordinateLookup struct {
	Latitude  float64
	Longitude float64
	Label     string
}

var (
	geocodingEndpoint = defaultGeocodingEndpoint
	geocodingClient   = &http.Client{Timeout: 4 * time.Second}
	coordinateCache   = struct {
		sync.RWMutex
		values map[string]coordinateLookup
	}{values: make(map[string]coordinateLookup)}
)

func hydrateForecastCoordinates(forecast *Forecast) {
	lookup, ok := geocodeLocation(forecast.Location, forecast.Region)
	if !ok {
		return
	}
	forecast.Latitude = lookup.Latitude
	forecast.Longitude = lookup.Longitude
	forecast.CoordinatesAvailable = true
	forecast.CoordinateLabel = lookup.Label
}

func geocodeLocation(location, region string) (coordinateLookup, bool) {
	location = strings.TrimSpace(location)
	region = strings.TrimSpace(region)
	if location == "" {
		return coordinateLookup{}, false
	}

	cacheKey := strings.ToLower(location + "|" + region)
	coordinateCache.RLock()
	cached, found := coordinateCache.values[cacheKey]
	coordinateCache.RUnlock()
	if found {
		return cached, true
	}

	query := url.Values{}
	query.Set("name", location)
	query.Set("count", "10")
	query.Set("language", "en")
	query.Set("format", "json")

	req, err := http.NewRequest(http.MethodGet, geocodingEndpoint+"?"+query.Encode(), nil)
	if err != nil {
		return coordinateLookup{}, false
	}
	req.Header.Set("User-Agent", "AMOC-Sentinel/1.0 (location globe)")

	resp, err := geocodingClient.Do(req)
	if err != nil {
		return coordinateLookup{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return coordinateLookup{}, false
	}

	var payload geoResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil || len(payload.Results) == 0 {
		return coordinateLookup{}, false
	}

	result := chooseGeoResult(payload.Results, location, region)
	lookup := coordinateLookup{
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Label:     coordinateLabel(result),
	}

	coordinateCache.Lock()
	coordinateCache.values[cacheKey] = lookup
	coordinateCache.Unlock()
	return lookup, true
}

func chooseGeoResult(results []geoResult, location, region string) geoResult {
	best := results[0]
	bestScore := geoResultScore(best, location, region)
	for _, result := range results[1:] {
		score := geoResultScore(result, location, region)
		if score > bestScore {
			best = result
			bestScore = score
		}
	}
	return best
}

func geoResultScore(result geoResult, location, region string) int {
	normalizedLocation := strings.ToLower(strings.TrimSpace(location))
	normalizedRegion := strings.ToLower(strings.TrimSpace(region))
	score := 0
	if strings.EqualFold(result.Name, normalizedLocation) {
		score += 8
	} else if strings.Contains(strings.ToLower(result.Name), normalizedLocation) {
		score += 3
	}
	if normalizedRegion == "" {
		return score
	}
	for _, candidate := range []string{result.Country, result.Admin1, result.Admin2} {
		normalizedCandidate := strings.ToLower(strings.TrimSpace(candidate))
		if normalizedCandidate == "" {
			continue
		}
		if normalizedCandidate == normalizedRegion {
			score += 12
		} else if strings.Contains(normalizedCandidate, normalizedRegion) || strings.Contains(normalizedRegion, normalizedCandidate) {
			score += 6
		}
	}
	return score
}

func coordinateLabel(result geoResult) string {
	parts := []string{result.Name}
	if result.Admin1 != "" && !strings.EqualFold(result.Admin1, result.Name) {
		parts = append(parts, result.Admin1)
	}
	if result.Country != "" && !strings.EqualFold(result.Country, result.Admin1) {
		parts = append(parts, result.Country)
	}
	return strings.Join(parts, ", ")
}
