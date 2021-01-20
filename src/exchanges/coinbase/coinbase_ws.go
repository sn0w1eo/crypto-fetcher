package coinbase

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// CoinbaseWS is used for WebSocket Protocol
type CoinbaseWS struct {
	Coinbase
	conn *websocket.Conn

	// finish chan
	done chan bool
}

// Dials to predefined URL in Protocol
// Returns error if setup isn't valid
// Returns error is Dial to server failed
func (cbw *CoinbaseWS) Dial() error {
	err := cbw.isValidSetup()
	if err != nil {
		cbw.log(err)
		return err
	}

	cbw.conn, _, err = websocket.DefaultDialer.Dial(protocols[cbw.protocol], nil)
	if err != nil {
		cbw.log(err)
		return err
	}
	return nil
}

// Reader is invoked by Serve method
// Reader starts read message out of connection and sends it to dedicated chan (e.g. tick)
// Returns on done
// On connection closed invoke method Stop()
// Closes dedicated chan on return
func (cbw *CoinbaseWS) reader() {
	defer close(cbw.tick)
	for {
		select {
		case <-cbw.done:
			_ = cbw.conn.Close()
			return
		default:
			_, msg, err := cbw.conn.ReadMessage()
			if err != nil {
				cbw.Stop("connection closed")
				return
			}
			msgType, err := parseMessageType(msg)
			switch msgType {
			case tickerChannelName:
				tick, err := parseTick(msg)
				if err != nil {
					cbw.log(err)
					continue
				}
				if cbw.tickUsage {
					cbw.tick <- tick
				}
			}
		}
	}
}

// Builds subscription message and sends it over websocket connection
// Important: connection should be established
// Returns error on send fail
func (cbw *CoinbaseWS) subscribe() error {
	var pairs []string
	for _, pair := range cbw.pairs {
		pairs = append(pairs, pair.String(PairDelimiter))
	}
	s := coinbaseSubscribe{
		Type:       "subscribe",
		ProductIds: pairs,
		Channels:   cbw.channels,
	}
	err := cbw.conn.WriteJSON(s)
	if err != nil {
		cbw.log(fmt.Errorf("WriteJSON failed: %s. message: %v", err.Error(), s))
	}
	return err
}

// Sends value to done chan. Logs reason of stop
// Disables tickUsage flag
func (cbw *CoinbaseWS) Stop(reason interface{}) {
	cbw.done <- true
	cbw.tickUsage = false
	if reason != nil {
		cbw.log(reason)
	}
}

// Serve invokes subscribe method and starts reader. And waits for done chan to finish
func (cbw *CoinbaseWS) Serve() error {
	err := cbw.subscribe()
	if err != nil {
		return err
	}

	go cbw.reader()

	for {
		select {
		case _, ok := <-cbw.done:
			if ok {
				close(cbw.done)
			}
			return nil
		}
	}
}
