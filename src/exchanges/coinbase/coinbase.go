package coinbase

import (
	"encoding/json"
	"fmt"
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"github.com/Sn0w1eo/crypto-fetcher/src/exchanges"
	"io"
	"log"
)

// Delimiter according to Coinbase API
const PairDelimiter = '-'

// Ticker channel name according to Coinbase API
const tickerChannelName = "ticker"

// Protocol represents which transport protocol will be used for data fetching
type Protocol int

const (
	WebSocket Protocol = iota + 1
)

// Represents URL according to Protocol
var protocols = map[Protocol]string{
	WebSocket: "wss://ws-feed.pro.coinbase.com",
}

// Coinbase is a base object for all other Protocol
type Coinbase struct {
	protocol Protocol

	tick      chan crypto.Tick
	tickUsage bool

	pairs    []crypto.Pair
	channels []string

	logger *log.Logger
}

// Base message exchange format provided by Coinbase API
type coinbaseMessage struct {
	Type  string `json:"type"`
	Error error  `json:"error"`
}

// Struct for subscription option according to rules of Coinbase API
type coinbaseSubscribe struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

// Creates new Coinbase Exchanger. Depends on the protocol
// WebSocket protocol defined as const
// error occurs if protocol has no implementation yet
func New(p Protocol) (exchanges.Exchanger, error) {
	switch p {
	case WebSocket:
		c := new(CoinbaseWS)
		c.protocol = p
		c.tick = make(chan crypto.Tick, 1)
		c.done = make(chan bool, 1)
		return c, nil
	default:
		return nil, fmt.Errorf("protocol not found: %d", p)
	}
}

// Returns error if some of required fields hasn't been initialized
func (cb *Coinbase) isValidSetup() error {
	if cb.logger == nil {
		return fmt.Errorf("logger isn't set")
	}
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
	cb.logger.Println(v)
}

// Sets logger as io.Writer interface
func (cb *Coinbase) SetLogger(w io.Writer) {
	cb.logger = log.New(w, "", log.Ldate|log.Ltime)
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

// Sets channels for subscription
// Note: Coinbase's channel names are predefined:
// ticker, heartbeat, level2
func (cb *Coinbase) SetChannels(channels ...string) error {
	if len(channels) == 0 {
		return fmt.Errorf("at least one channel should be set")
	}
	cb.channels = make([]string, len(channels))
	copy(cb.channels, channels)
	return nil
}

// Returns chan of crypto.Tick and sets tickUsage flag to true.
// Chan will be closed on connection lost or after Stop() method
func (cb *Coinbase) Ticker() <-chan crypto.Tick {
	cb.tickUsage = true
	return cb.tick
}

// Returns type of message from Coinbase server
func parseMessageType(msg []byte) (string, error) {
	cbMsg := coinbaseMessage{}
	err := json.Unmarshal(msg, &cbMsg)
	return cbMsg.Type, err
}

// Returns crypto.Tick out of msg. Error occurs on unmarshall or wrap failure
func parseTick(msg []byte) (tick crypto.Tick, err error) {
	t := Tick{}
	err = json.Unmarshal(msg, &t)
	if err != nil {
		return tick, fmt.Errorf("wrong message format, unable to unmarshall: %s", string(msg))
	}
	tick, err = crypto.WrapTicker(t)
	if err != nil {
		return tick, fmt.Errorf("failed to wrap tick: %v", t)
	}
	return
}
