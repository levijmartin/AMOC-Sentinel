package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchOpenMeteoDailyHistoricalRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("start_date"); got != "2024-04-05" {
			t.Fatalf("start_date = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"latitude":13.1,"longitude":-59.6,"timezone":"America/Barbados",
			"daily":{"time":["2024-04-05"],"weather_code":[61],"temperature_2m_max":[29.4],"temperature_2m_min":[24.2],"precipitation_sum":[7.5],"wind_speed_10m_max":[31.2]}
		}`))
	}))
	defer server.Close()

	originalEndpoint := archiveWeatherEndpoint
	originalClient := weatherClient
	archiveWeatherEndpoint = server.URL
	weatherClient = server.Client()
	defer func() {
		archiveWeatherEndpoint = originalEndpoint
		weatherClient = originalClient
	}()

	date, _ := time.Parse("2006-01-02", "2024-04-05")
	got, kind, err := fetchOpenMeteoDaily(context.Background(), 13.1, -59.6, date)
	if err != nil {
		t.Fatal(err)
	}
	if kind != "historical" || got.WeatherCode != 61 || got.PrecipitationMM != 7.5 || got.Timezone != "America/Barbados" {
		t.Fatalf("unexpected result: kind=%s weather=%+v", kind, got)
	}
}

func TestWeatherCodeLabel(t *testing.T) {
	cases := map[int]string{0: "Clear sky", 3: "Overcast", 61: "Rain", 80: "Rain showers", 95: "Thunderstorm"}
	for code, want := range cases {
		if got := weatherCodeLabel(code); got != want {
			t.Errorf("weatherCodeLabel(%d) = %q, want %q", code, got, want)
		}
	}
}
