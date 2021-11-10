package model

// Config models configuration file for app
type Config struct {
	TradePairs    []string `json:"TRADE_PAIRS"`
	SocketAddress string   `json:"SOCKET_ADDRESS"`
	ClearConsole  bool     `json:"CLEAR_CONSOLE"`
	Window        int     `json:"WINDOW"`
}
