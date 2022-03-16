package updater

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kovetskiy/cryptocompare-proxyd/internal/cache"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/cryptocompare"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

// Updater has only one function — to update the entries in the database by the
// given list of fsyms/tsyms.
type Updater struct {
	client cryptocompare.Client
	cache  cache.Cache

	fsyms []string
	tsyms []string

	updateInterval int

	done chan struct{}
}

// New instance of Updater.
func New(
	client cryptocompare.Client,
	cache cache.Cache,
	fsyms []string,
	tsyms []string,
	updateInterval int,
) (*Updater, error) {
	return &Updater{
		client:         client,
		cache:          cache,
		fsyms:          fsyms,
		tsyms:          tsyms,
		updateInterval: updateInterval,
		done:           make(chan struct{}),
	}, nil
}

// Update is a core function of Updater and is invoked by Serve(). It updates
// the prices in the cache storage.
func (updater *Updater) Update() error {
	startedAt := time.Now()

	log.Debugf(
		karma.
			Describe("fsyms", strings.Join(updater.fsyms, ",")).
			Describe("tsyms", strings.Join(updater.tsyms, ",")),
		"updater: updating the price list",
	)

	list, err := updater.client.GetPriceList(updater.fsyms, updater.tsyms)
	if err != nil {
		return karma.Format(err, "get price list")
	}

	// first we need to check if the price list has everything is according to
	// what we have asked for.
	for _, fsym := range updater.fsyms {
		if _, ok := list.Raw[fsym]; !ok {
			return fmt.Errorf(
				"the received price list (raw) doesn't have %q",
				fsym,
			)
		}

		if _, ok := list.Display[fsym]; !ok {
			return fmt.Errorf(
				"the received price list (display) doesn't have %q",
				fsym,
			)
		}

		for _, tsym := range updater.tsyms {
			if _, ok := list.Raw[fsym][tsym]; !ok {
				return fmt.Errorf(
					"the received price list of %s (raw) doesn't have %q",
					fsym,
					tsym,
				)
			}

			if _, ok := list.Display[fsym][tsym]; !ok {
				return fmt.Errorf(
					"the received price list of %s (display) doesn't have %q",
					fsym,
					tsym,
				)
			}
		}
	}

	for _, fsym := range updater.fsyms {
		for _, tsym := range updater.tsyms {
			err := updater.cache.Write(
				context.Background(),
				startedAt,
				fsym,
				tsym,
				list.Raw[fsym][tsym],
				list.Display[fsym][tsym],
			)
			if err != nil {
				return karma.Format(err, "cache write of %s to %s", fsym, tsym)
			}
		}
	}

	return nil
}

// Serve is expected to be running in a goroutine. It waits for the specified
// time and invokes the Update() method.
func (updater *Updater) Serve() error {
	log.Infof(nil, "the updater has started")

	for {
		select {
		case <-time.After(time.Duration(updater.updateInterval) * time.Second):
			//
		case <-updater.done:
			return nil
		}

		err := updater.Update()
		if err != nil {
			return err
		}
	}
}

// Close immediately stops the Updater instance.
func (updater *Updater) Close() {
	close(updater.done)
}
