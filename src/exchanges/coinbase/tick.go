package coinbase

import (
	"fmt"
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"strconv"
	"strings"
	"time"
)

// Coinbase Tick message format
type Tick struct {
	Time      time.Time `json:"time"`
	ProductId string    `json:"product_id"`
	Bid       string    `json:"best_bid"`
	Ask       string    `json:"best_ask"`
}

// Returns BestBid in float64 from Coinbase's Tick format
func (t Tick) BestBid() (bb float64, err error) {
	bb, err = strconv.ParseFloat(t.Bid, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert tick's best_bid to float64: %s", t.Ask)
	}
	return
}

// Returns BestAsk in float64 from Coinbase's Tick format
func (t Tick) BestAsk() (ba float64, err error) {
	ba, err = strconv.ParseFloat(t.Ask, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert tick's best_ask to float64: %s", t.Ask)
	}
	return
}

// Returns Timestamp from Coinbase's Tick format
func (t Tick) Timestamp() time.Time {
	return t.Time
}

// Returns pairs in crypto.Pair from Coinbase's Tick format
func (t Tick) Pair() (pair crypto.Pair, err error) {
	currencies := strings.Split(t.ProductId, string(PairDelimiter))
	if len(currencies) != 2 {
		return pair, fmt.Errorf("failed to split currencies by delimiter: %c", PairDelimiter)
	}
	pair, err = crypto.NewPair(currencies[0], currencies[1])
	if err != nil {
		return
	}
	return
}
