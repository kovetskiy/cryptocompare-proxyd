package cache

import (
	"time"

	"github.com/kovetskiy/cryptocompare-proxyd/internal/cryptocompare"
	"github.com/uptrace/bun"
)

// Entity describes a cache record.
type Entity interface {
	StoredAt() time.Time
	FromSymbol() string
	ToSymbol() string
	RawPrice() cryptocompare.RawPrice
	DisplayPrice() cryptocompare.DisplayPrice
}

type entity struct {
	bun.BaseModel `bun:"table:pricelist,alias:u"`

	ID int64 `bun:",pk,autoincrement"`

	At time.Time `bun:"at,type:timestamp"`

	Fsym string `bun:"fsym,unique:fsym_tsym"`

	Tsym string `bun:"tsym,unique:fsym_tsym"`

	Raw     cryptocompare.RawPrice     `bun:"raw,type:jsonb"`
	Display cryptocompare.DisplayPrice `bun:"display,type:jsonb"`
}

func (entity entity) StoredAt() time.Time {
	return entity.At
}

func (entity entity) RawPrice() cryptocompare.RawPrice {
	return entity.Raw
}

func (entity entity) DisplayPrice() cryptocompare.DisplayPrice {
	return entity.Display
}

func (entity entity) FromSymbol() string {
	return entity.Fsym
}

func (entity entity) ToSymbol() string {
	return entity.Tsym
}
