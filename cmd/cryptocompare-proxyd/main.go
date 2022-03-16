package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kovetskiy/cryptocompare-proxyd/internal/cache"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/config"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/cryptocompare"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/server"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/updater"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"

	"github.com/docopt/docopt-go"
)

var (
	version = "[manual build]"
	usage   = "cryptocompare-proxyd " + version + `

This is a reverse proxy service with caching for the cryptocompare API.

Usage:
  cryptocompare-proxyd [options] [-R]
  cryptocompare-proxyd -h | --help
  cryptocompare-proxyd --version

Options:
  -R --read-only       Run in read only mode.
  -c --config <value>  Use the specified configuration file. [default: /etc/cryptocompare-proxyd.conf]
  --debug              Print debug messages.
  -h --help            Show this screen.
  --version            Show version.
`
)

// Opts describes command-line options.
type Opts struct {
	ValueConfig  string `docopt:"--config"`
	FlagDebug    bool   `docopt:"--debug"`
	FlagReadOnly bool   `docopt:"--read-only"`
}

func main() {
	args, err := docopt.ParseArgs(usage, nil, "cryptocompare-proxyd "+version)
	if err != nil {
		// the only case is that the developer made a mistake in the usage
		// variable
		panic(err)
	}

	var opts Opts
	err = args.Bind(&opts)
	if err != nil {
		log.Fatalf(err, "unable to bind the specified command-line arguments")
	}

	if opts.FlagDebug {
		log.SetLevel(log.LevelDebug)
	}

	config, err := config.Load(opts.ValueConfig)
	if err != nil {
		log.Fatalf(err, "unable to load the configuration")
	}

	cache, err := cache.New(
		config.DatabaseAddress,
		config.DatabaseName,
		config.DatabaseUsername,
		config.DatabasePassword,
	)
	if err != nil {
		log.Fatalf(err, "unable to initialize cache instance")
	}

	err = cache.Boot()
	if err != nil {
		log.Fatalf(err, "unable to boot cache instance")
	}

	defer cache.Close()

	client, err := cryptocompare.New(version)
	if err != nil {
		log.Fatalf(err, "unable to initialize cryptocompare http client")
	}

	var refresher *updater.Updater
	if !opts.FlagReadOnly {
		refresher, err = updater.New(
			client,
			cache,
			config.Fsyms,
			config.Tsyms,
			config.UpdateInterval,
		)
		if err != nil {
			log.Fatalf(err, "unable to initialize updater")
		}

		err = refresher.Update()
		if err != nil {
			log.Fatalf(err, "unable to update the symbols data")
		}
	}

	server, err := server.New(
		config.ListenAddress,
		cache,
		client,
		config.CacheTTL,
	)
	if err != nil {
		log.Fatalf(err, "unable to initialize http server instance")
	}

	serve(server, refresher)
}

func serve(
	server *server.Server,
	refresher *updater.Updater,
) {
	var serveError error

	done := make(chan struct{})
	workers := &sync.WaitGroup{}

	workers.Add(1)
	go func() {
		defer workers.Done()

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			// this is a corner case when received a SIGINT signal and have
			// to shutdown the HTTP server.
			serveError = karma.Format(
				err,
				"http server: unable to listen and serve",
			)

			close(done)
		}
	}()

	if refresher != nil {
		workers.Add(1)
		go func() {
			defer workers.Done()

			err := refresher.Serve()
			if err != nil {
				serveError = karma.Format(
					err,
					"updater: unable to serve",
				)

				close(done)
			}
		}()
	}

	// The server can shut down in two cases:
	// it's either user/container-orchestrator interaction: by sending os signal
	// or an error
	//
	// in both cases we want to try to gracefully shut down everything we can

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case signal := <-signals:
		log.Infof(
			nil,
			"the server is shutting down due to received signal '%v'",
			signal,
		)

	case <-done:
		if serveError != nil {
			log.Errorf(
				serveError,
				"the server is shutting down due to an error",
			)
		}
	}

	err := server.Close()
	if err != nil {
		log.Errorf(err, "unable to gracefully shutdown the http server")
	} else {
		log.Infof(nil, "the server has gracefully shut down")
	}

	if refresher != nil {
		refresher.Close()

		log.Infof(nil, "the updater has gracefully shut down")
	}

	workers.Wait()
}
