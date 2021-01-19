package crypto

import (
	"fmt"
	"strings"
)

// Currency is an object to store information/methods of concrete currency
// id represents symbol
type Currency struct {
	id string
}

// Sets id of Currency in uppercase.
// Restrictions: id's length should be a least 2 symbols
func (c *Currency) SetId(symbol string) error {
	if len(symbol) < 2 {
		return fmt.Errorf("bad symbol: %s", symbol)
	}
	c.id = strings.ToUpper(symbol)
	return nil
}

// Returns id of Currency
func (c *Currency) Id() string {
	return c.id
}

// Creates new Currency and sets id
func NewCurrency(id string) (c Currency, err error) {
	err = c.SetId(id)
	return c, err
}
