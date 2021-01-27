package storage

import (
	"fmt"
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	"github.com/Sn0w1eo/crypto-fetcher/src/storage/mysql"
)

type Type int

const (
	MySQL Type = iota + 1
)

// Storage used like abstraction interface to manipulate with different storage providers
type Storage interface {
	// Open used for pass DSN information to open storage source correctly
	Open(dsn string) error
	// WriteTick writes Ticker to storage
	WriteTick(ticker crypto.Ticker) error
	// Close closes current storage
	Close() error
}

// Creates new storage based on Type. Returns error if Type not found
func New(storageType Type) (Storage, error) {
	switch storageType {
	case MySQL:
		return mysql.New()
	default:
		return nil, fmt.Errorf("storage not implemented yet")
	}
}
