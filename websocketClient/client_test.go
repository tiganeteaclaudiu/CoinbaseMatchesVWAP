package websocketClient

import (
	"assignment/helpers"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/require"
)

// WSTEST package usage based on source at https://github.com/posener/wstest/blob/master/dialer_test.go
// dialer for test purposes, can't handle multiple websocket connections concurrently
type handler struct {
	*websocket.Conn
	upgrader websocket.Upgrader
	Upgraded chan struct{}
}

func (s *handler) connect(w http.ResponseWriter, r *http.Request) {
	defer close(s.Upgraded)
	var err error
	s.Conn, err = s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
}

func (s *handler) Close() error {
	if s.Conn == nil {
		return nil
	}
	return s.Conn.Close()
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve")
	fmt.Println(r)
	switch r.URL.Path {
	case "/ws":
		s.connect(w, r)

	case "/ws/delay":
		<-time.After(500 * time.Millisecond)
		s.connect(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

// TestSocketClient_SubscribeToMatches tests the SubscribeToMatches method
func TestSocketClient_SubscribeToMatches(t *testing.T) {
	expectedSentMessage := `
	{
	   "type":"subscribe",
	   "channels":[
		  {
			 "name":"matches",
			 "product_ids":["BTC-USD"]
		  }
	   ]
	}
`
	t.Parallel()
	var (
		s    = &handler{Upgraded: make(chan struct{})}
		d    = wstest.NewDialer(s)
		done = make(chan struct{})
	)

	c, _, err := d.Dial("ws://example.org/ws", nil)
	require.Nil(t, err)

	<-s.Upgraded

	client := &socketClient{
		conn: c,
		done: done,
	}

	go client.SubscribeToMatches(helpers.GetSubscribeToMatchesMessage([]string{"BTC-USD"}))

	// read message that method sent to websocket
	_, m, err := s.ReadMessage()
	if string(m) != expectedSentMessage {
		t.Errorf("expected %s got %s", expectedSentMessage, string(m))
	}

	err = c.Close()

	err = s.Close()
}

// TestSocketClient_Read tests the Read method
func TestSocketClient_Read(t *testing.T) {
	dataPointMessage := `{"type":"match","trade_id":234704065,"maker_order_id":"8c5d05a4-41c8-41f8-abc0-a49f09072bfd","taker_order_id":"d9827159-d335-4a68-931d-5a66ee0f1de3","side":"sell","size":"0.00002416","price":"64632.95","product_id":"BTC-USD","sequence":30995303205,"time":"2021-11-11T08:35:56.588997Z"}`
	t.Parallel()
	var (
		s    = &handler{Upgraded: make(chan struct{})}
		d    = wstest.NewDialer(s)
		done = make(chan struct{})
	)

	c, _, err := d.Dial("ws://example.org/ws", nil)
	require.Nil(t, err)

	<-s.Upgraded

	client := &socketClient{
		conn: c,
		done: done,
	}

	output := make(chan []byte)

	// start reading from websocket channel
	client.Read(output)

	// write message to channel
	go func() {
		err := s.WriteMessage(websocket.TextMessage, []byte(dataPointMessage))
		require.Nil(t, err)
		done <- struct{}{}
	}()

	// Read method should load message written to output chan
	result := <-output
	if string(result) != dataPointMessage {
		t.Errorf("expected %s got %s", dataPointMessage, string(result))
	}

}
