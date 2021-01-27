package coinbase

import (
	"fmt"
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTick_Timestamp(t *testing.T) {
	tick := Tick{
		Time: time.Unix(1, 0),
	}
	assert.Equal(t, time.Unix(1, 0), tick.Timestamp())
}

func TestTick_BestBid(t *testing.T) {
	cases := []struct {
		bid      string
		expected float64
		hasError bool
	}{
		{"123.2", 123.2, false},
		{"ABC", 0, true},
		{"", 0, true},
	}
	for _, testCase := range cases {
		tick := Tick{Bid: testCase.bid}
		actual, err := tick.BestBid()
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestTick_BestAsk(t *testing.T) {
	cases := []struct {
		ask      string
		expected float64
		hasError bool
	}{
		{"123.2", 123.2, false},
		{"ABC", 0, true},
		{"", 0, true},
	}
	for _, testCase := range cases {
		tick := Tick{Ask: testCase.ask}
		actual, err := tick.BestAsk()
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestTick_Pair(t *testing.T) {
	t.Run("split by coinbase delimiter", func(t *testing.T) {
		cases := []struct {
			productId string
			hasError  bool
		}{
			{fmt.Sprintf("BTC%cUSD", PairDelimiter), false},
			{fmt.Sprintf("ETH%cXRP", PairDelimiter), false},
			{"BTC/USD", true},
			{"ETH|USD", true},
			{"BTCADA", true},
		}
		for _, testCase := range cases {
			tick := Tick{
				ProductId: testCase.productId,
			}
			_, err := tick.Pair()
			if testCase.hasError {
				assert.Error(t, err)
				continue
			}
			assert.NoError(t, err)
		}
	})
	t.Run("wraptick returns crypto.Tick", func(t *testing.T) {
		btc_usd, _ := crypto.NewPair("BTC", "USD")
		cases := []struct {
			productId    string
			expectedPair crypto.Pair
			hasError     bool
		}{
			{fmt.Sprintf("BTC%cUSD", PairDelimiter), btc_usd, false},
			{"A-B", crypto.Pair{}, true},
		}
		for _, testCase := range cases {
			tick := Tick{
				ProductId: testCase.productId,
			}
			pair, err := tick.Pair()
			if testCase.hasError {
				assert.Error(t, err)
				continue
			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedPair, pair)
		}
	})
}

func Test_parseTick(t *testing.T) {
	btc_usd, _ := crypto.NewPair("BTC", "USD")

	cases := []struct {
		msg          []byte
		expectedTick crypto.Tick
		hasError     bool
	}{

		{
			msg: []byte("{\"time\":\"1970-01-01T00:00:00Z\", \"best_bid\":\"2\", \"best_ask\":\"1\", \"product_id\":\"BTC-USD\"}"),
			expectedTick: crypto.Tick{
				T:   time.Unix(0, 0),
				P:   btc_usd,
				Bid: 2,
				Ask: 1,
			},
			hasError: false,
		},
		{
			msg:          []byte("{\"time\":\"1970-01-01T00:00:00Z\", \"best_bid\":2, \"best_ask\":1, \"product_id\":\"BTC-USD\"}"),
			expectedTick: crypto.Tick{},
			hasError:     true,
		},
	}

	for _, testCase := range cases {
		tick, err := parseTick(testCase.msg)
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)

		assert.Equal(t, time.Duration(0), testCase.expectedTick.T.Sub(tick.T))
		assert.Equal(t, testCase.expectedTick.P, tick.P)
		assert.Equal(t, testCase.expectedTick.Bid, tick.Bid)
		assert.Equal(t, testCase.expectedTick.Ask, tick.Ask)

	}
}
