package exchanges

import (
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"io"
)

// Exchanger represents communication between different Exchanges
type Exchanger interface {
	Dial() error
	Serve() error
	Stop(reason interface{})
	SetLogger(writer io.Writer)
	SetPairs(...crypto.Pair) error
	SetChannels(...string) error
	Ticker() <-chan crypto.Tick
}
