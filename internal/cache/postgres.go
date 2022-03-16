package cache

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kovetskiy/cryptocompare-proxyd/internal/cryptocompare"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// We make sure postgres implements the Cache interface.
var _ Cache = (*postgres)(nil)

type postgres struct {
	address  string
	database string
	username string
	password string

	db *bun.DB
}

func (postgres *postgres) Boot() error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		postgres.username,
		postgres.password,
		postgres.address,
		postgres.database,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	// the official documentation of bun recommends adding this setting, but
	// it's too noisy. It does contain some useful information such as
	// queries and how much time it takes to perform the query.
	//
	// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	db.RegisterModel((*entity)(nil))

	log.Debugf(nil, "postgres: ensure table schema")

	_, err := db.NewCreateTable().
		Model((*entity)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return karma.Format(err, "postgres: create/ensure table")
	}

	postgres.db = db

	return nil
}

func (postgres *postgres) Close() {
	err := postgres.db.Close()
	if err != nil {
		log.Errorf(err, "postgres: close")
	}
}

func (postgres *postgres) Write(
	ctx context.Context,
	at time.Time,
	fromSymbol string,
	toSymbols string,
	raw cryptocompare.RawPrice,
	display cryptocompare.DisplayPrice,
) error {
	_, err := postgres.db.NewInsert().Model(&entity{
		At:      at,
		Fsym:    fromSymbol,
		Tsym:    toSymbols,
		Raw:     raw,
		Display: display,
	}).On("CONFLICT ON CONSTRAINT fsym_tsym DO UPDATE").Exec(ctx)
	if err != nil {
		return karma.Format(err, "postgres: insert")
	}

	return nil
}

func (postgres *postgres) Read(
	ctx context.Context,
	fromSymbols []string,
	toSymbolss []string,
	ttl int,
) ([]Entity, error) {
	entities := []entity{}

	err := postgres.db.NewSelect().
		Model((*entity)(nil)).
		Where(
			"fsym IN (?) AND tsym IN (?) AND at > NOW() - INTERVAL '?'",
			bun.In(fromSymbols),
			bun.In(toSymbolss),
			ttl,
		).
		Scan(ctx, &entities)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, karma.Format(err, "postgres: select")
	}

	// converting implementations to interfaces
	result := make([]Entity, len(entities))
	for i, entity := range entities {
		result[i] = Entity(entity)
	}

	return result, nil
}
