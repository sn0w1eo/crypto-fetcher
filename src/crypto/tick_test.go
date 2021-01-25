package crypto

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

type fTick struct {
	primary   string
	secondary string
	bestBid   string
	bestAsk   string
	timestamp int64
}

func (tick fTick) Pair() (Pair, error) {
	return NewPair(tick.primary, tick.secondary)
}

func (tick fTick) BestBid() (float64, error) {
	return strconv.ParseFloat(tick.bestBid, 64)
}

func (tick fTick) BestAsk() (float64, error) {
	return strconv.ParseFloat(tick.bestAsk, 64)
}
func (tick fTick) Timestamp() time.Time {
	return time.Unix(tick.timestamp, 0)
}

func TestWrapTicker(t *testing.T) {
	cases := []struct {
		fakeTick fTick
		want     Tick
		hasError bool
	}{
		{
			fakeTick: fTick{"btc", "usd", "1", "10", 1},
			want:     Tick{Pair: Pair{primary: Currency{id: "BTC"}, secondary: Currency{id: "USD"}}, BestBid: 1.0, BestAsk: 10.0, Time: time.Unix(1, 0)},
			hasError: false,
		},
		{
			fakeTick: fTick{"", "usd", "1", "10", 0},
			want:     Tick{},
			hasError: true,
		},
		{
			fakeTick: fTick{"btc", "", "1", "10", 1},
			want:     Tick{},
			hasError: true,
		},
		{
			fakeTick: fTick{"btc", "usd", "WRONG", "10", 1},
			want:     Tick{},
			hasError: true,
		},
		{
			fakeTick: fTick{"btc", "usd", "1", "WRONG", 1},
			want:     Tick{},
			hasError: true,
		},
	}
	for _, testCase := range cases {
		tick, err := WrapTicker(testCase.fakeTick)
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, testCase.want, tick)
	}

}
