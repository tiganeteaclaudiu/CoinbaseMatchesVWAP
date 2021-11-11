package utils

import (
	"CoinbaseMatchesVWAP/model"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Aggregator - aggregates trade data for multiple trade pairs
type Aggregator struct {
	Utils        map[string]*VWAPUtil
	tradingPairs []string
	config       model.Config
}

// NewAggregator - initializes a new aggregator based on a config object
func NewAggregator(config model.Config) *Aggregator {
	// initialize VWAP Utils for each trade pair in configuration
	utils := map[string]*VWAPUtil{}
	for _, pair := range config.TradePairs {
		utils[pair] = NewVWAPUtil(config.Window, pair)
	}
	return &Aggregator{
		Utils:        utils,
		tradingPairs: config.TradePairs,
		config:       config,
	}
}

// ToString - prints formatted aggregated trade data to output
func (ag *Aggregator) ToString() string {
	var pairs []string

	// output VWAP for all trading pairs in aggregator
	for _, pair := range ag.tradingPairs {
		pairs = append(pairs, ag.Utils[pair].ToString())
	}

	return strings.Join(pairs, "\n")
}

// ToOutput - prints formatted aggregated trade data to output
func (ag *Aggregator) ToOutput() {
	// clears console before printing if configured so
	if ag.config.ClearConsole == true {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	fmt.Println(ag.ToString())
}
