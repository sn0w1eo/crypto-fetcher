package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPair_DefaultDelimiter(t *testing.T) {
	var expected rune
	assert.IsType(t, expected, PairDefaultDelimiter)
}

func TestPair_Primary(t *testing.T) {
	id := "BTC"
	pair := Pair{
		primary: Currency{id: id},
	}
	assert.Equal(t, id, pair.Primary())
}

func TestPair_Secondary(t *testing.T) {
	id := "USD"
	pair := Pair{
		secondary: Currency{id: id},
	}
	assert.Equal(t, id, pair.Secondary())
}

func TestNewPair(t *testing.T) {
	cases := []struct {
		primary   string
		secondary string
		expected  Pair
		hasError  bool
	}{
		{"BTC", "USD", Pair{primary: Currency{id: "BTC"}, secondary: Currency{id: "USD"}}, false},
		{"ETH", "USD", Pair{primary: Currency{id: "ETH"}, secondary: Currency{id: "USD"}}, false},
		{"ETH", "", Pair{}, true},
		{"", "USD", Pair{}, true},
	}
	for _, testCase := range cases {
		pair, err := NewPair(testCase.primary, testCase.secondary)
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, pair)
	}
}

func TestPair_String(t *testing.T) {
	cases := []struct {
		primary    string
		secondary  string
		delimiters []rune
		expected   string
	}{
		{"BTC", "USD", []rune{'/'}, "BTC/USD"},
		{"BTC", "USD", []rune{'|', '-'}, "BTC|USD"},
		{"BTC", "USD", nil, fmt.Sprintf("BTC%cUSD", PairDefaultDelimiter)},
	}
	for _, testCase := range cases {
		pair, _ := NewPair("BTC", "USD")
		str := pair.String(testCase.delimiters...)
		assert.Equal(t, testCase.expected, str)
	}
}
