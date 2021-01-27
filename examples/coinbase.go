package main

import (
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"github.com/Sn0w1eo/crypto-fetcher/src/exchanges"
	"github.com/Sn0w1eo/crypto-fetcher/src/storage"
	"log"
	"os"
)

const MySQL_DSN = "user:password@tcp(127.0.0.1:3306)/dbname"
const DefDelimiter = '-'

func TickWriteFunc(tickCh <-chan crypto.Tick, st storage.Storage, name string) {
	for {
		select {
		case tick, ok := <-tickCh:
			if !ok {
				return
			}
			err := st.WriteTick(tick)
			if err != nil {
				log.Fatal(name, err)
			}
		}
	}
}

func TickerWriter(tickCh <-chan crypto.Tick, st storage.Storage, pairs []crypto.Pair) {
	// Named channels map
	channels := map[string]chan crypto.Tick{}
	for _, pair := range pairs {
		ch := make(chan crypto.Tick, 1)
		name := pair.String(DefDelimiter)
		channels[name] = ch

		go TickWriteFunc(ch, st, name)
	}

	for {
		select {
		case t, ok := <-tickCh:
			if !ok {
				for key, _ := range channels {
					close(channels[key])
				}
				return
			}
			pair, _ := t.Pair()
			name := pair.String(DefDelimiter)
			channels[name] <- t
		}
	}
}

func main() {
	eth_btc, _ := crypto.NewPair("eth", "btc")
	btc_usd, _ := crypto.NewPair("btc", "usd")
	btc_eur, _ := crypto.NewPair("btc", "eur")
	pairs := []crypto.Pair{eth_btc, btc_usd, btc_eur}

	// Create new Exchanger
	ex, err := exchanges.New(exchanges.Coinbase, exchanges.WebSocket)
	if err != nil {
		panic(err)
	}

	// Set pairs
	err = ex.SetPairs(pairs...)
	if err != nil {
		panic(err)
	}

	// Set logger as io.Writer
	ex.SetLogger(os.Stdout)

	// Tick channel
	tickCh := ex.Ticker()

	// New Storage
	st, err := storage.New(storage.MySQL)
	if err != nil {
		panic(err)
	}

	// Open connection to Storage by DSN
	err = st.Open(MySQL_DSN)
	if err != nil {
		panic(err)
	}

	go TickerWriter(tickCh, st, pairs)

	// Dial to Exchanger
	err = ex.Dial()
	if err != nil {
		panic(err)
	}

	// Serves connection
	err = ex.Serve()
	if err != nil {
		panic(err)
	}

}
