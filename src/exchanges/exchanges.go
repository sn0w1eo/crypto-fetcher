package exchanges

import (
	"fmt"
	"github.com/Sn0w1eo/crypto-fetcher/src/exchanges/coinbase"
)

// Protocol represents which transport protocol will be used for data fetching
type Protocol int

const (
	WebSocket Protocol = iota + 1
)

type Exchange int

const (
	Coinbase Exchange = iota + 1
)

func New(exchange Exchange, protocol Protocol) (Exchanger, error) {
	switch exchange {
	case Coinbase:
		switch protocol {
		case WebSocket:
			return coinbase.NewWS(), nil
		default:
			return nil, fmt.Errorf("coinbase protocol not implemented yet: %d", protocol)
		}
	default:
		return nil, fmt.Errorf("exchange not found")
	}
}
