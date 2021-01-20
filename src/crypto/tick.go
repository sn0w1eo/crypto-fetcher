package crypto

import "time"

// Common struct for Tick exchange
type Tick struct {
	Time    time.Time
	Pair    Pair
	BestBid float64
	BestAsk float64
}

// Wraps Ticker object to crypto.Tick object
func WrapTicker(ticker Ticker) (tick Tick, err error) {
	tick.Pair, err = ticker.Pair()
	if err != nil {
		return tick, nil
	}
	tick.Time = ticker.Timestamp()
	tick.BestBid, err = ticker.BestBid()
	if err != nil {
		return tick, nil
	}
	tick.BestAsk, err = ticker.BestAsk()
	if err != nil {
		return tick, nil
	}
	return
}
