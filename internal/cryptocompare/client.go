package cryptocompare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

const (
	apiURI = "https://min-api.cryptocompare.com/data/pricemultifull"
)

type remoteResponse struct {
	Response string `json:"Response"`
	Message  string `json:"Message"`
}

// Client can talk to cryptocompare and return current prices.
type Client interface {
	GetPriceList(fsyms []string, tsyms []string) (*PriceList, error)
}

type client struct {
	version string
	http    *http.Client
}

// New creates a new client to talk to cryptocompare.
func New(version string) (Client, error) {
	return &client{
		http: &http.Client{},
	}, nil
}

// GetPriceList makes a HTTP request and returns the current prices.
func (client *client) GetPriceList(
	fsyms []string,
	tsyms []string,
) (*PriceList, error) {
	query := url.Values{}
	query.Add("fsyms", strings.Join(fsyms, ","))
	query.Add("tsyms", strings.Join(tsyms, ","))

	uri := apiURI + "?" + query.Encode()

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, karma.Format(err, "new request")
	}

	request.Header.Set("User-Agent", "cryptocompare-proxyd/"+client.version)

	log.Debugf(nil, "client: GET request to %s", uri)

	response, err := client.http.Do(request)
	if err != nil {
		return nil, karma.Format(err, "http GET request")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code, expected: %v, but got %v",
			http.StatusOK,
			response.Status,
		)
	}

	// we read the contents completely instead of streaming into json.Decoder
	// because we are going to use the output in debug messages in case of error
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, karma.Format(err, "read response body")
	}

	var list PriceList
	err = json.Unmarshal(contents, &list)
	if err != nil {
		log.Errorf(
			karma.
				Describe("response", string(contents)).
				Reason(err),
			"unable to decode json response",
		)

		return nil, karma.Format(err, "decode json response")
	}

	if len(list.Raw) == 0 && len(list.Display) == 0 {
		// seems like a corner case so we should double check that this is not
		// an error response
		var remoteError remoteResponse

		err := json.Unmarshal(contents, &remoteError)
		if err == nil && remoteError.Response == "Error" {
			// our case is when we don't have problems decoding the JSON, otherwise
			// it would fail even in previous json.Unmarshal cases.
			return nil, karma.
				Describe("contents", string(contents)).
				Format(nil, "the remote server returned an error: %s", remoteError.Message)
		}
	}

	return &list, nil
}
