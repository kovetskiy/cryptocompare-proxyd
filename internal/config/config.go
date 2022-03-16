package config

import (
	"github.com/kovetskiy/ko"
)

// Config is a configuration variables stored in environment variables on in a
// file.
type Config struct {
	// ListenAddress to listen for HTTP connections on.
	ListenAddress string `yaml:"listen_address" required:"true" env:"LISTEN_ADDRESS" default:":8080"`

	// UpdateInterval is a duration of time (seconds) which should be used for refreshing
	// cached data from the cryptocompare service.
	//
	// Note that the cryptocompare service caches data itself, the ttl is 10
	// seconds, so there is not much sense to assign this value to less than 10
	// seconds.
	UpdateInterval int `yaml:"update_interval" required:"true" env:"UPDATE_INTERVAL" default:"30"`

	// CacheTTL is a duration of time (seconds) to treat cache entries as
	// expired.
	CacheTTL int `yaml:"cache_ttl" required:"true" env:"CACHE_TTL" default:"120"`

	// Fsyms is a cryptocurrency symbols of interest.
	Fsyms []string `yaml:"fsyms,inline" required:"true" env:"FSYMS" default:"[BTC]"`

	// Tsyms is a cryptocurrency symbols list to convert into.
	Tsyms []string `yaml:"tsyms" required:"true" env:"TSYMS" default:"[USD]"`

	// DatabaseAddress is an address of a database to connect to.
	DatabaseAddress string `yaml:"database_address" required:"true" env:"DATABASE_ADDRESS" default:"localhost:5432"`

	// DatabaseName is a name of database the program should talk with.
	DatabaseName string `yaml:"database_name" required:"true" env:"DATABASE_NAME" default:"cryptocompare-proxyd-dev"`

	// DatabaseUsername is a username that should be used to auth in the
	// database.
	DatabaseUsername string `yaml:"database_username" required:"true" env:"DATABASE_USERNAME" default:"cryptocompare-proxyd-dev"`

	// DatabasePassword is a password that should be used to auth in the
	// database.
	DatabasePassword string `yaml:"database_password" required:"false" env:"DATABASE_PASSWORD" default:"cryptocompare-proxyd-dev"`
}

// Load read the given file or reads environment variables, returns instance of
// Config with values from a file or environment variables.
func Load(path string) (*Config, error) {
	config := &Config{}
	err := ko.Load(path, config, ko.RequireFile(false))
	if err != nil {
		return nil, err
	}

	return config, nil
}
