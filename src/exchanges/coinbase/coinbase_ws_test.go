package coinbase

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestNewWS(t *testing.T) {
	cbw := NewWS()
	assert.NotEqual(t, nil, cbw)
}

func TestCoinbaseWS_Dial(t *testing.T) {
	cbw := NewWS()
	cbw.Dial()
}

func TestCoinbaseWS_Stop(t *testing.T) {
	t.Run("Stop() sends true to done chan", func(t *testing.T) {
		cbw := CoinbaseWS{
			done: make(chan bool, 1),
		}
		cbw.Stop(nil)
		actual := <-cbw.done
		close(cbw.done)
		assert.Equal(t, true, actual)
	})

	t.Run("Stop() param will be logged", func(t *testing.T) {
		b := bytes.Buffer{}
		w := bufio.NewWriter(&b)
		cbw := CoinbaseWS{
			Coinbase: Coinbase{
				logger: log.New(w, "", 0),
			},
			done: make(chan bool, 1),
		}
		expected := "user stopped"
		cbw.Stop(expected)
		w.Flush()
		close(cbw.done)
		assert.Equal(t, expected+"\n", b.String())
	})
}
