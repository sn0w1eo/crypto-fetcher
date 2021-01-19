package crypto

// Ticker interface is used to provide  information from different Exchanges
type Ticker interface {
	Pair() (Pair, error)
	BestBid() (float64, error)
	BestAsk() (float64, error)
	Timestamp() int64
}
