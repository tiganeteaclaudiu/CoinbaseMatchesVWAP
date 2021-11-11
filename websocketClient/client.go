package websocketClient

import (
	"errors"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

// SocketClient - defines the methods of a socket client
type SocketClient interface {
	SubscribeToMatches(subscriptionMessage string) error
	Read(output chan []byte)
	Close() error
}

type socketClient struct {
	conn *websocket.Conn
	done chan struct{}
}

// NewSocketClient - initializes a new socket client
func NewSocketClient(address *string, done chan struct{}) (SocketClient, error) {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "wss", Host: *address}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	return &socketClient{
		conn: c,
		done: done,
	}, nil
}

// SubscribeToMatches - sends subscribe message to matches channel
func (cl *socketClient) SubscribeToMatches(subscriptionMessage string) error {
	// starts at 0 retries
	return cl.subscribeToMatches(0, subscriptionMessage)
}

// subscribeToMatches - sends subscribe message to matches channel
func (cl *socketClient) subscribeToMatches(retries int, subscriptionMessage string) error {
	// limit to 5 retries
	if retries == 5 {
		return errors.New("too many retries")
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	err := cl.conn.WriteMessage(websocket.TextMessage, []byte(subscriptionMessage))
	if err != nil {
		log.Println("failed to write:", err)
		time.Sleep(time.Second * 2)
		// retry
		cl.subscribeToMatches(retries+1, subscriptionMessage)
		return err
	}

	return nil
}

// Read - starts reading from websocket channel
func (cl *socketClient) Read(output chan []byte) {
	go func() {
		defer close(cl.done)
		for {
			_, message, err := cl.conn.ReadMessage()
			if err != nil {
				log.Println("failed to read:", err)
				return
			}

			output <- message
		}
	}()
}

// Close - closes websocket connection
func (cl *socketClient) Close() error {
	return cl.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
