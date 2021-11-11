package main

import (
	"assignment/model"
	"assignment/utils"
	"testing"
	"time"
)

var config = model.Config{
	TradePairs:   []string{"BTC-USD", "ETH-USD", "ETH-BTC"},
	ClearConsole: true,
	Window:       200,
}

// TestValidateConfig_Valid - attempts to validate a valid config
func TestValidateConfig_Valid(t *testing.T) {
	config := model.Config{
		TradePairs:    []string{"BTC-USD", "ETH-USD", "ETH-BTC"},
		SocketAddress: "test",
		ClearConsole:  true,
		Window:        200,
	}
	err := validateConfig(config)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

// TestValidateConfig_Invalid - attempts to validate an invalid config
func TestValidateConfig_Invalid(t *testing.T) {
	// socket address is missing
	config := model.Config{
		TradePairs:   []string{"BTC-USD", "ETH-USD", "ETH-BTC"},
		ClearConsole: true,
		Window:       200,
	}
	err := validateConfig(config)
	if err == nil {
		t.Error(`expected "No SOCKET_ADDRESS in configuration" error, got nil`)
	}
}

// TestStartRead - tests startRead function
func TestStartRead(t *testing.T) {
	expected := `Trading Pair: BTC-USD, VWAP: 64632.950000
Trading Pair: ETH-USD, VWAP: NaN
Trading Pair: ETH-BTC, VWAP: NaN`
	read := make(chan []byte)
	aggregator := utils.NewAggregator(config)
	done := make(chan struct{})
	go startRead(read, aggregator, done)

	testMsg := `{"type":"match","trade_id":234704065,"maker_order_id":"8c5d05a4-41c8-41f8-abc0-a49f09072bfd","taker_order_id":"d9827159-d335-4a68-931d-5a66ee0f1de3","side":"sell","size":"0.00002416","price":"64632.95","product_id":"BTC-USD","sequence":30995303205,"time":"2021-11-11T08:35:56.588997Z"}`
	// send message to read channel. Message should be processed in startRead for loop
	read <- []byte(testMsg)

	time.Sleep(time.Second * 2)

	// check if message was read and successfully processed by aggregator
	if aggregator.ToString() != expected {
		t.Errorf("expected %s got %s", expected, aggregator.ToString())
	}
}

// TestStartRead_InvalidDataPoint - tests startRead function
// Invalid data point message is sent to read channel
func TestStartRead_InvalidDataPoint(t *testing.T) {
	expected := `Trading Pair: BTC-USD, VWAP: NaN
Trading Pair: ETH-USD, VWAP: NaN
Trading Pair: ETH-BTC, VWAP: NaN`
	read := make(chan []byte)
	aggregator := utils.NewAggregator(config)
	done := make(chan struct{})
	go startRead(read, aggregator, done)

	testMsg := `invalid message`
	// send message to read channel. Message should be processed in startRead for loop
	read <- []byte(testMsg)

	time.Sleep(time.Second * 5)

	for {
		select {
		// check if done message was receive
		case <-done:
			// aggregator should have no data (Nan VWAP should be present for all trading pairs)
			if aggregator.ToString() != expected {
				t.Errorf("expected %s got %s", expected, aggregator.ToString())
			}
			return
		case <-time.After(2 * time.Second):
			// check if 2 seconds passed without receiving a done message
			t.Errorf("no done message received")
			return
		}
	}
}
