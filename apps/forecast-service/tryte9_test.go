package main

import "testing"

func TestTryte9FullRangeRoundTrip(t *testing.T) {
	for value := tryte9Min; value <= tryte9Max; value++ {
		trits := encodeBalancedTernary9(value)
		got, err := decodeBalancedTernary(trits)
		if err != nil {
			t.Fatalf("decodeBalancedTernary(%d): %v", value, err)
		}
		if got != value {
			t.Fatalf("round trip %d = %d via %v", value, got, trits)
		}
	}
}

func TestTryte9PrefilterForcesMuteReset(t *testing.T) {
	got := evaluateTryte9(42, Tryte9BraidT1CCW, false, 1, 1, 0)
	if got.Prefilter.Passed {
		t.Fatal("prefilter passed non-zero u+v+w")
	}
	if got.EffectiveOpcode != Tryte9MuteReset || !got.Muted {
		t.Fatalf("effective opcode = %d, muted = %v", got.EffectiveOpcode, got.Muted)
	}
}

func TestTryte9BypassPreservesRequestedOpcode(t *testing.T) {
	got := evaluateTryte9(-42, Tryte9BraidT1InverseCW, true, 1, 1, 1)
	if !got.Prefilter.Passed {
		t.Fatal("bypassed prefilter did not pass")
	}
	if got.EffectiveOpcode != Tryte9BraidT1InverseCW || got.Muted {
		t.Fatalf("effective opcode = %d, muted = %v", got.EffectiveOpcode, got.Muted)
	}
}

func TestSetunUsesFullNineTritRange(t *testing.T) {
	for _, value := range []float64{-0.5, 0, 0.5} {
		engine := setunFromFloat(value)
		if engine.RawValue() > uint32(setunMax) {
			t.Fatalf("raw value %d exceeds %d", engine.RawValue(), setunMax)
		}
		if len(engine.ToTrits()) != tryte9TritCount {
			t.Fatalf("got %d trits", len(engine.ToTrits()))
		}
	}
}
