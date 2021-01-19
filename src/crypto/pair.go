package crypto

import "fmt"

// Represents default delimiter for Pair
const PairDefaultDelimiter = '-'

// Pair is an object to store information/methods of concrete pair
// primary is buy currency
// secondary is sell currency
type Pair struct {
	primary   Currency
	secondary Currency
}

// Returns primary Currency id
func (p *Pair) Primary() string {
	return p.primary.Id()
}

// Returns secondary Currency id
func (p *Pair) Secondary() string {
	return p.secondary.Id()
}

// Creates new Pair. If one of two Currency creation failed returns Currency's error.
func NewPair(primary string, secondary string) (p Pair, err error) {
	p.primary, err = NewCurrency(primary)
	if err != nil {
		return p, err
	}
	p.secondary, err = NewCurrency(secondary)
	if err != nil {
		return p, err
	}
	return p, err
}

// Represents Pair as string with default delimiter
func (p *Pair) String() string {
	return fmt.Sprintf("%s%c%s", p.primary.Id(), PairDefaultDelimiter, p.secondary.Id())
}
