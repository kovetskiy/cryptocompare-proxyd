package cache

import (
	"context"
	"time"

	"github.com/kovetskiy/cryptocompare-proxyd/internal/cryptocompare"
)

// Cache is a persistent storage for price lists.
type Cache interface {
	// Boot is required to ensure that the instance is ready for read/write
	// operations.
	// It may connect to a database, it depends on the implementation.
	Boot() error

	// Close tries to close all current connections.
	Close()

	// Read returns a list of entities by the given parameters.
	Read(
		ctx context.Context,
		FromSymbols []string,
		ToSymbols []string,
		ttl int,
	) ([]Entity, error)

	// Write saves the specified data into the internal storage.
	Write(
		ctx context.Context,
		at time.Time,
		fromSymbol string,
		toSymbol string,
		raw cryptocompare.RawPrice,
		display cryptocompare.DisplayPrice,
	) error
}

// New instance of cache, currently postgres supported only.
func New(
	address string,
	database string,
	username string,
	password string,
) (Cache, error) {
	return &postgres{
		address:  address,
		database: database,
		username: username,
		password: password,
	}, nil
}
