package cryptocompare

// PriceList is a list of prices from the cryptocompare service without any
// modifications.
type PriceList struct {
	Raw     map[string]map[string]RawPrice     `json:"RAW"`
	Display map[string]map[string]DisplayPrice `json:"DISPLAY"`
}

// RawPrice is expected to use by machines and not expected to be shown to
// humans.
type RawPrice struct {
	Price           float64 `json:"PRICE"`
	Volume24Hour    float64 `json:"VOLUME24HOUR"`
	Volume24HourTo  float64 `json:"VOLUME24HOURTO"`
	Open24Hour      float64 `json:"OPEN24HOUR"`
	High24Hour      float64 `json:"HIGH24HOUR"`
	Low24Hour       float64 `json:"LOW24HOUR"`
	Change24Hour    float64 `json:"CHANGE24HOUR"`
	ChangePct24Hour float64 `json:"CHANGEPCT24HOUR"`
	Supply          float64 `json:"SUPPLY"`
	Mktcap          float64 `json:"MKTCAP"`
}

// DisplayPrice is a kind of price that can be shown to a human. Prices contain
// dollar sign, thousands are written as K, etc.
type DisplayPrice struct {
	Price           string `json:"PRICE"`
	Volume24Hour    string `json:"VOLUME24HOUR"`
	Volume24HourTo  string `json:"VOLUME24HOURTO"`
	Open24Hour      string `json:"OPEN24HOUR"`
	High24Hour      string `json:"HIGH24HOUR"`
	Low24Hour       string `json:"LOW24HOUR"`
	Change24Hour    string `json:"CHANGE24HOUR"`
	ChangePct24Hour string `json:"CHANGEPCT24HOUR"`
	Supply          string `json:"SUPPLY"`
	Mktcap          string `json:"MKTCAP"`
}
