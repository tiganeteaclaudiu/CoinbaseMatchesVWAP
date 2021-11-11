package utils

import (
	"fmt"
	"testing"
)

// CreateVWAPUtil - creates a complete VWAPUtil for testing purposes
func CreateVWAPUtil(window int) VWAPUtil {
	return VWAPUtil{
		Pair:            "BTC-USD",
		cumulatedVolume: 10,
		volumes:         []float64{1, 2, 3, 2, 3},
		maxPrice:        5.1,
		minPrice:        1.1,
		prices:          []float64{3.2, 1.1, 2.22, 5.1, 2.13},
		window:          window,
		cumulatedTPV:    33.733333,
	}
}

func roundToTwoDecimal(val float64) string {
	return fmt.Sprintf("%.2f", val)
}

// TestNewVWAPUtil - ideal test for function
func TestNewVWAPUtil(t *testing.T) {
	result := NewVWAPUtil(200, "BTC-USD")
	if result == nil {
		t.Error("VWAP Util is nil")
	}
}

// TestVWAPUtil_removeLast - tests removeLast method of VWAPUtil
func TestVWAPUtil_removeLast(t *testing.T) {
	util := CreateVWAPUtil(100)

	util.removeLast()

	if util.cumulatedVolume != 9.000000 {
		t.Errorf("expected %f got %f", 9.000000, util.cumulatedVolume)
	}

	if util.maxPrice != 5.100000 {
		t.Errorf("expected %f got %f", 5.100000, util.maxPrice)
	}

	if util.minPrice != 1.100000 {
		t.Errorf("expected %f got %f", 1.100000, util.minPrice)
	}

	if len(util.prices) != 4 {
		t.Errorf("expected %d got %d", 4, len(util.prices))
	}

	if len(util.volumes) != 4 {
		t.Errorf("expected %d got %d", 4, len(util.volumes))
	}
}

// TestVWAPUtil_GetTypicalPrice - tests GetTypicalPrice method of VWAPUtil
func TestVWAPUtil_GetTypicalPrice(t *testing.T) {
	expected := 2.733333
	util := CreateVWAPUtil(100)

	TPV := util.GetTypicalPrice(2)

	if roundToTwoDecimal(TPV) != roundToTwoDecimal(expected) {
		t.Errorf("expected %f got %f", expected, TPV)
	}
}

// TestVWAPUtil_Add - tests Add method of VWAPUtil
func TestVWAPUtil_Add(t *testing.T) {
	util := CreateVWAPUtil(100)

	util.Add(0.5, 2)

	if util.cumulatedVolume != 12.000000 {
		t.Errorf("expected %f got %f", 12.000000, util.cumulatedVolume)
	}

	if util.maxPrice != 5.100000 {
		t.Errorf("expected %f got %f", 5.100000, util.maxPrice)
	}

	if util.cumulatedTPV != 24.400000 {
		t.Errorf("expected %f got %f", 24.400000, util.cumulatedTPV)
	}

	// new min price is set
	if util.minPrice != 0.500000 {
		t.Errorf("expected %f got %f", 0.500000, util.minPrice)
	}

	if len(util.prices) != 6 {
		t.Errorf("expected %d got %d", 6, len(util.prices))
	}

	if len(util.volumes) != 6 {
		t.Errorf("expected %d got %d", 6, len(util.volumes))
	}
}

// TestVWAPUtil_Add_FullWindow - tests Add method of VWAPUtil
// length of trades is equal to window so last trade must be dropped before adding new one
func TestVWAPUtil_Add_FullWindow(t *testing.T) {
	// there are 5 trades in util by default
	util := CreateVWAPUtil(5)

	util.Add(0.5, 2)

	// one trade is added and the oldest is deleted, so volume should not increase by new volume (2)
	if util.cumulatedVolume != 11.000000 {
		t.Errorf("expected %f got %f", 11.000000, util.cumulatedVolume)
	}

	if util.maxPrice != 5.100000 {
		t.Errorf("expected %f got %f", 5.100000, util.maxPrice)
	}

	// new min price is set
	if util.minPrice != 0.500000 {
		t.Errorf("expected %f got %f", 0.500000, util.minPrice)
	}

	// length should stay equal to window (one element is added and last one is deleted)
	if len(util.prices) != 5 {
		t.Errorf("expected %d got %d", 5, len(util.prices))
	}

	// length should stay equal to window (one element is added and last one is deleted)
	if len(util.volumes) != 5 {
		t.Errorf("expected %d got %d", 5, len(util.volumes))
	}
}

// TestVWAPUtil_GetVWAP - tests GetVWAP method of VWAPUtil
func TestVWAPUtil_GetVWAP(t *testing.T) {
	// there are 5 trades in util by default
	util := CreateVWAPUtil(5)
	expected := 3.373333

	result := util.GetVWAP()

	if roundToTwoDecimal(result) != roundToTwoDecimal(expected) {
		t.Errorf("expected %f got %f", 3.373333, result)
	}
}

// TestVWAPUtil_ToString - tests ToString method of VWAPUtil
func TestVWAPUtil_ToString(t *testing.T) {
	expected := `Trading Pair: BTC-USD, VWAP: 3.373333`
	// there are 5 trades in util by default
	util := CreateVWAPUtil(5)

	result := util.ToString()

	if result != expected {
		t.Errorf("expected %s got %s", expected, result)
	}
}
