package coinbase

import (
	"encoding/json"
	"fmt"
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"io"
	"log"
)

// Delimiter according to Coinbase API
const PairDelimiter = '-'

// Ticker channel name according to Coinbase API
const tickerChannelName = "ticker"

// Coinbase is a base object for all other Protocol
type Coinbase struct {
	tick chan crypto.Tick

	pairs    []crypto.Pair
	channels []string

	logger *log.Logger
}

// Base message exchange format provided by Coinbase API
type coinbaseMessage struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

// Struct for subscription option according to rules of Coinbase API
type coinbaseSubscribe struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

// Returns error if some of required fields hasn't been initialized
func (cb *Coinbase) isValidSetup() error {
	if len(cb.pairs) == 0 {
		return fmt.Errorf("pairs aren't set")
	}
	if len(cb.channels) == 0 {
		return fmt.Errorf("channels aren't set")
	}
	return nil
}

// Simple log function
func (cb *Coinbase) log(v interface{}) {
	if cb.logger != nil {
		cb.logger.Println(v)
	}
}

// Sets logger as io.Writer interface
func (cb *Coinbase) SetLogger(w io.Writer) {
	cb.logger = log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
}

// Sets slice of crypto.Pair, will be used for subscribe later on
func (cb *Coinbase) SetPairs(pairs ...crypto.Pair) error {
	if len(pairs) == 0 {
		return fmt.Errorf("at least one pair should be set")
	}
	cb.pairs = make([]crypto.Pair, len(pairs))
	copy(cb.pairs, pairs)
	return nil
}

// Returns chan of crypto.Tick
// Chan will be closed on connection lost or after Stop() method
func (cb *Coinbase) Ticker() <-chan crypto.Tick {
	if cb.tick != nil {
		return cb.tick
	}
	cb.channels = append(cb.channels, tickerChannelName)
	cb.tick = make(chan crypto.Tick, 1)
	return cb.tick
}

// Returns type of message from Coinbase server
func parseMessageType(msg []byte) (string, error) {
	cbMsg := coinbaseMessage{}
	err := json.Unmarshal(msg, &cbMsg)
	if cbMsg.Type == "" {
		return "", fmt.Errorf("message type is empty")
	}
	return cbMsg.Type, err
}
