package main

import "testing"

func TestChooseGeoResultUsesRegionAndExactName(t *testing.T) {
	results := []geoResult{
		{Name: "Portland", Latitude: 45.52345, Longitude: -122.67621, Country: "United States", Admin1: "Oregon"},
		{Name: "Bridgetown", Latitude: 13.10732, Longitude: -59.62021, Country: "Barbados", Admin1: "Saint Michael"},
	}

	got := chooseGeoResult(results, "Bridgetown", "Barbados")
	if got.Name != "Bridgetown" || got.Country != "Barbados" {
		t.Fatalf("chooseGeoResult() = %#v, want Bridgetown, Barbados", got)
	}
}

func TestCoordinateLabelAvoidsDuplicateRegion(t *testing.T) {
	result := geoResult{Name: "Los Angeles", Admin1: "California", Country: "United States"}
	got := coordinateLabel(result)
	want := "Los Angeles, California, United States"
	if got != want {
		t.Fatalf("coordinateLabel() = %q, want %q", got, want)
	}
}
