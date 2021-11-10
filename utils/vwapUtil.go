package utils

import (
	"assignment/helpers"
	"fmt"
	"math"
)

// VWAPUtil represents a set of utilities designed to calculate the VWAP (Volume weighted average price) of a specific trading pair
type VWAPUtil struct {
	// Pair represents the trading pair the utility is initialized for
	// Defined as identifier
	Pair string
	// prices - slice of trading prices in slide window
	prices []float64
	// volumes - slice of trading volumes in slide window
	volumes []float64
	// maxPrice - maximum trading price in slide window
	maxPrice float64
	// minPrice - minimum trading price in slide window
	minPrice float64
	// cumulatedVolume - total volume of trades in slide window
	cumulatedVolume float64
	// cumulatedTPV - total TPV of trades in slide window
	cumulatedTPV float64
	// window - size of trading window. Limits calculations of VWAP to last N trades.
	window int
}

// NewVWAPUtil initializes a new VWAPUtil for a traiding pair
func NewVWAPUtil(window int, pair string) *VWAPUtil {
	// window represents maximum number of data points to slide
	return &VWAPUtil{
		window:   window,
		minPrice: math.MaxFloat64,
		Pair:     pair,
		prices:   []float64{},
	}
}

func (ag *VWAPUtil) removeLast() {
	// subtract volume of oldest data point
	ag.cumulatedVolume -= ag.volumes[0]
	ag.volumes = ag.volumes[1:]

	// check if maximum price is price we are about to delete
	if ag.maxPrice == ag.prices[0] {
		ag.maxPrice = helpers.GetMaxFloat(ag.prices[1:])
	}

	// check if minimum price is price we are about to delete
	if ag.minPrice == ag.prices[0] {
		ag.minPrice = helpers.GetMinFloat(ag.prices[1:])
	}

	ag.prices = ag.prices[1:]
}

// GetTypicalPrice - calculates TPV of current slide window
func (ag *VWAPUtil) GetTypicalPrice(lastPrice float64) float64 {
	return (ag.maxPrice + ag.minPrice + lastPrice) / 3
}

// Add - adds a new data point to slide window
func (ag *VWAPUtil) Add(newPrice, newVolume float64) {
	// if max number of trades in window has been reached, discard tail
	if len(ag.prices) == ag.window {
		ag.removeLast()
	}

	// check if new data point's price is either the largest or smallest in slide window
	if newPrice > ag.maxPrice {
		ag.maxPrice = newPrice
	}

	if newPrice < ag.minPrice {
		ag.minPrice = newPrice
	}

	// keep track of new data point's price and volume
	ag.prices = append(ag.prices, newPrice)
	ag.volumes = append(ag.volumes, newVolume)

	// add volume to cumulated volume in slide window
	ag.cumulatedVolume += newVolume
	// add TPV to cumulated TPV in slide window
	ag.cumulatedTPV = ag.GetTypicalPrice(newPrice) * ag.cumulatedVolume
}

// GetVWAP - calculate VWAP of current slide window
func (ag *VWAPUtil) GetVWAP() float64 {
	return ag.cumulatedTPV / ag.cumulatedVolume
}

// ToString - output as string
func (ag *VWAPUtil) ToString() string {
	return fmt.Sprintf("Trading Pair: %s, VWAP: %f", ag.Pair, ag.GetVWAP())
}
