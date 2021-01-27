package coinbase

import (
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestCoinbase_SetPairs(t *testing.T) {
	btc_usd, _ := crypto.NewPair("btc", "usd")
	btc_eur, _ := crypto.NewPair("btc", "eur")
	eth_usd, _ := crypto.NewPair("eth", "usd")
	eth_eur, _ := crypto.NewPair("eth", "eur")
	btc_eth, _ := crypto.NewPair("btc", "eth")

	cases := []struct {
		pairs    []crypto.Pair
		len      int
		hasError bool
	}{
		{[]crypto.Pair{btc_usd, btc_eur, eth_eur, btc_eth}, 4, false},
		{[]crypto.Pair{eth_usd}, 1, false},
		{[]crypto.Pair{}, 0, true},
	}
	for _, testCase := range cases {
		c := Coinbase{}
		err := c.SetPairs(testCase.pairs...)
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		assert.Equal(t, testCase.len, len(c.pairs))
	}
}

func TestCoinbase_SetLogger(t *testing.T) {
	c := Coinbase{}
	c.SetLogger(os.Stdout)
	assert.Equal(t, os.Stdout, c.logger.Writer())
}

func TestCoinbase_Ticker(t *testing.T) {
	t.Run("Ticker is singleton", func(t *testing.T) {
		expected := make(chan crypto.Tick, 1)
		c := Coinbase{
			tick:     expected,
			channels: []string{},
		}
		assert.EqualValues(t, expected, c.Ticker())
		assert.EqualValues(t, []string{}, c.channels)
	})

	t.Run("Ticker appends channel for subscribe and initializes chan for reading tick messages", func(t *testing.T) {
		c := Coinbase{}
		gotChan := c.Ticker()
		assert.Equal(t, c.channels[0], tickerChannelName)

		passValue := crypto.Tick{T: time.Now()}
		c.tick <- passValue
		assert.Equal(t, passValue, <-gotChan)
		assert.Equal(t, 1, cap(c.tick))
	})
}

func Test_isValidSetup(t *testing.T) {
	cases := []struct {
		pairs    []crypto.Pair
		channels []string
		hasError bool
	}{
		{[]crypto.Pair{{}, {}}, []string{tickerChannelName}, false},
		{nil, []string{tickerChannelName}, true},
		{[]crypto.Pair{}, []string{tickerChannelName}, true},
		{[]crypto.Pair{{}, {}}, nil, true},
		{[]crypto.Pair{{}, {}}, []string{}, true},
	}
	for _, testCase := range cases {
		c := Coinbase{}
		c.pairs = testCase.pairs
		c.channels = testCase.channels
		err := c.isValidSetup()
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
	}
}

func Test_parseMessageType(t *testing.T) {
	cases := []struct {
		message         []byte
		expectedMessage coinbaseMessage
		hasError        bool
	}{
		{[]byte("{\"type\":\"anything\"}"), coinbaseMessage{Type: "anything"}, false},
		{[]byte("{\"field\":\"anything\"}"), coinbaseMessage{}, true},
		{[]byte("{}"), coinbaseMessage{}, true},
		{[]byte("some message"), coinbaseMessage{}, true},
	}
	for _, testCase := range cases {
		actualMessageType, err := parseMessageType(testCase.message)
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedMessage.Type, actualMessageType)
	}
}
