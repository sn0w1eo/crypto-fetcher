package crypto

import "time"

// Common struct for Tick exchange
type Tick struct {
	T   time.Time
	P   Pair
	Bid float64
	Ask float64
}

// Returns pair in crypto.Pair
func (t Tick) Pair() (Pair, error) {
	return t.P, nil
}

// Returns Timestamp
func (t Tick) Timestamp() time.Time {
	return t.T
}

// Returns BestBid in float64
func (t Tick) BestBid() (float64, error) {
	return t.Bid, nil
}

// Returns BestAsk in float64
func (t Tick) BestAsk() (float64, error) {
	return t.Ask, nil
}

// Wraps Ticker object to crypto.Tick object
func WrapTicker(ticker Ticker) (tick Tick, err error) {
	tick.T = ticker.Timestamp()
	tick.P, err = ticker.Pair()
	if err != nil {
		return tick, err
	}
	tick.Bid, err = ticker.BestBid()
	if err != nil {
		return tick, err
	}
	tick.Ask, err = ticker.BestAsk()
	if err != nil {
		return tick, err
	}
	return tick, nil
}
