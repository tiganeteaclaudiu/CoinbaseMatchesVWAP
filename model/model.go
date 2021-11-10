package model

// DataPoint models a single data point received from coinbase websocket matches channel
type DataPoint struct {
	Type      string `json:"type"`
	Side      string `json:"side"`
	Size      string `json:"size"`
	Price     string `json:"price"`
	ProductID string `json:"product_id"`
}