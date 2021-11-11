package utils

import (
	"assignment/model"
	"testing"
)

// TestNewAggregator - ideal test for function
func TestNewAggregator(t *testing.T) {
	pairs := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	config := model.Config{
		TradePairs:    pairs,
		SocketAddress: "test",
		ClearConsole:  false,
		Window:        200,
	}
	result := NewAggregator(config)
	if result == nil {
		t.Error("aggregator is nil")
	}
	if len(result.Utils) == 0 {
		t.Error("failed to create utils for trading pairs")
	}
}

// TestNewAggregator_NoPairs covers case when no trading pairs are received.
// In this case, no utils should be created for aggregat
func TestNewAggregator_NoPairs(t *testing.T) {
	pairs := []string{}
	config := model.Config{
		TradePairs:    pairs,
		SocketAddress: "test",
		ClearConsole:  false,
		Window:        200,
	}
	result := NewAggregator(config)
	if result == nil {
		t.Error("aggregator is nil")
	}
	if len(result.Utils) != 0 {
		t.Errorf("created wrong number of utils: %d", len(result.Utils))
	}
}

// TestAggregator_ToString - tests ToString method of Aggregator
func TestAggregator_ToString(t *testing.T) {
	expectedValue := `Trading Pair: BTC-USD, VWAP: 1.000000
Trading Pair: ETH-USD, VWAP: 2.000000
Trading Pair: ETH-BTC, VWAP: 2.000000`
	// initialize a new aggregator with 3 trade pairs
	pairs := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	config := model.Config{
		TradePairs:    pairs,
		SocketAddress: "test",
		ClearConsole:  false,
		Window:        200,
	}
	result := NewAggregator(config)

	// Add a datapoint to the util for each trading pair
	result.Utils["BTC-USD"].Add(1,2)
	result.Utils["ETH-BTC"].Add(2,3)
	result.Utils["ETH-USD"].Add(2,3)

	str := result.ToString()
	if str != expectedValue {
		t.Errorf("expected \n%s \ngot \n%s", expectedValue, str)
	}
}