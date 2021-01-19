package crypto

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_Currency_SetId(t *testing.T) {
	cases := []struct {
		actual   string
		expected string
		hasError bool
	}{
		{"usd", "USD", false},
		{"eTH", "ETH", false},
		{"BtC", "BTC", false},
		{"ar", "AR", false},
		{"x", "", true},
		{"", "", true},
	}
	for _, testCase := range cases {
		c := Currency{}
		err := c.SetId(testCase.actual)
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		privateId := reflect.ValueOf(c).FieldByName("id")
		assert.Equal(t, testCase.expected, privateId.String())
	}
}

func Test_NewCurrency(t *testing.T) {
	cases := []struct {
		actual   string
		expected string
		hasError bool
	}{
		{"USD", "USD", false},
		{"AR", "AR", false},
		{"X", "", true},
		{"", "", true},
	}
	for _, testCase := range cases {
		c, err := NewCurrency(testCase.actual)
		if testCase.hasError {
			assert.Error(t, err)
			continue
		}
		privateId := reflect.ValueOf(c).FieldByName("id")
		assert.Equal(t, testCase.expected, privateId.String())
	}
}

func Test_Currency_Id(t *testing.T) {
	cases := []struct {
		actual   string
		expected string
	}{
		{"USD", "USD"},
		{"AR", "AR"},
	}
	for _, testCase := range cases {
		c, _ := NewCurrency(testCase.actual)
		assert.Equal(t, testCase.expected, c.Id())
	}
}
