package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	tryte9Min       = -9841
	tryte9Max       = 9841
	tryte9Bias      = 9841
	tryte9TritCount = 9
)

type Tryte9Opcode uint8

const (
	Tryte9Hold Tryte9Opcode = iota
	Tryte9BraidT1CCW
	Tryte9BraidT1InverseCW
	Tryte9MuteReset
)

type Tryte9Prefilter struct {
	Bypassed bool `json:"bypassed"`
	U        int  `json:"u"`
	V        int  `json:"v"`
	W        int  `json:"w"`
	Sum      int  `json:"sum"`
	Passed   bool `json:"passed"`
}

type Tryte9Response struct {
	Value           int               `json:"value"`
	Trits           []int             `json:"trits"`
	TritLabel       string            `json:"tritLabel"`
	Rank            uint16            `json:"rank"`
	RequestedOpcode Tryte9Opcode      `json:"requestedOpcode"`
	EffectiveOpcode Tryte9Opcode      `json:"effectiveOpcode"`
	OpcodeBits      string            `json:"opcodeBits"`
	Action          string            `json:"action"`
	Muted           bool              `json:"muted"`
	Packed          uint32            `json:"packed"`
	PackedHex       string            `json:"packedHex"`
	Prefilter       Tryte9Prefilter   `json:"prefilter"`
	Boundary        map[string]string `json:"boundary"`
}

func handleTryte9EvaluateAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	value, err := parseBoundedInt(r, "value", tryte9Min, tryte9Max)
	if err != nil {
		http.Error(w, `{"error":"value must be an integer from -9841 through 9841"}`, http.StatusBadRequest)
		return
	}
	opcodeValue, err := parseBoundedInt(r, "opcode", 0, 3)
	if err != nil {
		http.Error(w, `{"error":"opcode must be 0, 1, 2, or 3"}`, http.StatusBadRequest)
		return
	}

	bypass, err := strconv.ParseBool(defaultString(r.URL.Query().Get("bypass"), "false"))
	if err != nil {
		http.Error(w, `{"error":"bypass must be true or false"}`, http.StatusBadRequest)
		return
	}
	u, err := parseOptionalTrit(r, "u")
	if err != nil {
		http.Error(w, `{"error":"u must be -1, 0, or 1"}`, http.StatusBadRequest)
		return
	}
	v, err := parseOptionalTrit(r, "v")
	if err != nil {
		http.Error(w, `{"error":"v must be -1, 0, or 1"}`, http.StatusBadRequest)
		return
	}
	wv, err := parseOptionalTrit(r, "w")
	if err != nil {
		http.Error(w, `{"error":"w must be -1, 0, or 1"}`, http.StatusBadRequest)
		return
	}

	writeJSON(w, evaluateTryte9(value, Tryte9Opcode(opcodeValue), bypass, u, v, wv))
}

func evaluateTryte9(value int, requested Tryte9Opcode, bypass bool, u, v, w int) Tryte9Response {
	trits := encodeBalancedTernary9(value)
	rank := uint16(value + tryte9Bias)
	prefilter := Tryte9Prefilter{Bypassed: bypass, U: u, V: v, W: w, Sum: u + v + w}
	prefilter.Passed = bypass || prefilter.Sum == 0

	effective := requested
	if !prefilter.Passed {
		effective = Tryte9MuteReset
	}
	packed := uint32(rank) | uint32(effective)<<15

	return Tryte9Response{
		Value:           value,
		Trits:           trits,
		TritLabel:       joinTrits(trits),
		Rank:            rank,
		RequestedOpcode: requested,
		EffectiveOpcode: effective,
		OpcodeBits:      fmt.Sprintf("%02b", effective),
		Action:          effective.action(),
		Muted:           effective == Tryte9MuteReset,
		Packed:          packed,
		PackedHex:       fmt.Sprintf("0x%05X", packed),
		Prefilter:       prefilter,
		Boundary: map[string]string{
			"module": "Sentinel Decode optional prefilter",
			"status": "software reference demo; not FPGA/ASIC or physical YBCO qualification",
		},
	}
}

func encodeBalancedTernary9(value int) []int {
	trits := make([]int, tryte9TritCount)
	n := value
	for i := tryte9TritCount - 1; i >= 0; i-- {
		remainder := n % 3
		n /= 3
		switch remainder {
		case 2:
			remainder = -1
			n++
		case -2:
			remainder = 1
			n--
		}
		trits[i] = remainder
	}
	return trits
}

func decodeBalancedTernary(trits []int) (int, error) {
	value := 0
	for _, trit := range trits {
		if trit < -1 || trit > 1 {
			return 0, fmt.Errorf("invalid trit %d", trit)
		}
		value = value*3 + trit
	}
	return value, nil
}

func (o Tryte9Opcode) action() string {
	switch o {
	case Tryte9Hold:
		return "identity_hold"
	case Tryte9BraidT1CCW:
		return "braid_t1_counterclockwise"
	case Tryte9BraidT1InverseCW:
		return "braid_t1_inverse_clockwise"
	default:
		return "illegal_mute_reset"
	}
}

func parseBoundedInt(r *http.Request, key string, min, max int) (int, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return value, nil
}

func parseOptionalTrit(r *http.Request, key string) (int, error) {
	raw := defaultString(r.URL.Query().Get(key), "0")
	value, err := strconv.Atoi(raw)
	if err != nil || value < -1 || value > 1 {
		return 0, fmt.Errorf("invalid trit")
	}
	return value, nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func joinTrits(trits []int) string {
	parts := make([]string, len(trits))
	for i, trit := range trits {
		parts[i] = strconv.Itoa(trit)
	}
	return strings.Join(parts, " ")
}
