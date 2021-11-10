package main

import (
	"assignment/helpers"
	"assignment/model"
	"assignment/utils"
	"assignment/websocketClient"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// loads configuration from JSON file
func loadConfig() (model.Config, error) {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := model.Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		return model.Config{}, err
	}

	return configuration, nil
}

func timeExecution(start time.Time) {
	fmt.Println("Execution time: ", time.Since(start))
}

func main() {
	defer timeExecution(time.Now())
	// load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to load configuration. err: %v", err))
		return
	}
	socketAddr := flag.String("addr", config.SocketAddress, "http service address")
	if socketAddr == nil {
		fmt.Println(fmt.Sprintf("Websocket address not configured. err: %v", err))
		return
	}

	// channel that handles shutdown
	done := make(chan struct{})
	// channel that handles read from websocket client
	read := make(chan []byte)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// initialize an VWAP Utility for each trading pair
	aggregator := utils.NewAggregator(config)

	// initialize websocket client
	client, err := websocketClient.NewSocketClient(socketAddr, done)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to open socket client: %v", err))
		return
	}
	// subscribe to matches channel for all trading pairs
	err = client.SubscribeToMatches(helpers.GetSubscribeToMatchesMessage(config.TradePairs))
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to subscribe client to matches channel: %v", err))
		return
	}
	// start reading from websocket
	client.Read(read)

	// start reading from channel
	go startRead(read, aggregator)

	// listen for shutdown events (keyboard interrupt or done)
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			client.Close()

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// startRead - starts reading from input channel, outputs data to aggregator
func startRead(read chan []byte, aggregator *utils.Aggregator) {
	var err error
	for {
		select {
		case message := <-read:
			var dataPoint model.DataPoint
			err = json.Unmarshal(message, &dataPoint)
			if err != nil {
				fmt.Println(err)
				return
			}

			// ignore messages other than matches
			if dataPoint.Type != "last_match" && dataPoint.Type != "match" {
				continue
			}

			// parse price and volume of tranzaction (size)
			price, err := strconv.ParseFloat(dataPoint.Price, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			size, err := strconv.ParseFloat(dataPoint.Size, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			// add data point to VWAP util based on Trading Pair in message (Product ID)
			aggregator.Utils[dataPoint.ProductID].Add(price, size)
			// output all Trading Pairs data using aggregator
			fmt.Println(aggregator.ToString())
		}
	}
}
