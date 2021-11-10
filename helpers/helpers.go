package helpers

import (
	"fmt"
	"math"
	"strings"
)

// GetSubscribeToMatchesMessage builds a matches channel subscription message based on desired trading pairs
func GetSubscribeToMatchesMessage(pairs []string) string {
	var wrappedPairs []string
	// wrap each pair in double quotes
	for _, pair := range pairs {
		wrappedPairs = append(wrappedPairs, `"`+pair+`"`)
	}
	return fmt.Sprintf(`
	{
	   "type":"subscribe",
	   "channels":[
		  {
			 "name":"matches",
			 "product_ids":[%s]
		  }
	   ]
	}
`, strings.Join(wrappedPairs, ","))
}

// GetMaxFloat - returns maximum in float64 slice
func GetMaxFloat(vals []float64) float64 {
	var max float64
	for _, val := range vals {
		if val > max {
			max = val
		}
	}

	return max
}

// GetMinFloat - returns minimum in float64 slice
func GetMinFloat(vals []float64) float64 {
	min := math.MaxFloat64
	for _, val := range vals {
		if val < min {
			min = val
		}
	}

	return min
}
