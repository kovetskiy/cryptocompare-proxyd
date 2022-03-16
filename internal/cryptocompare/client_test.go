package cryptocompare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetPriceList_ReturnsValidData(t *testing.T) {
	test := assert.New(t)

	client, err := New("testing")
	if !test.NoError(err) {
		test.FailNow("should be able to initialize an instance of client")
	}

	test.NotNil(client)

	price, err := client.GetPriceList([]string{"BTC"}, []string{"USD", "EUR"})
	test.NoError(err)
	test.NotNil(price)

	test.Len(price.Raw, 1)
	test.Len(price.Display, 1)

	test.Contains(price.Raw, "BTC")

	test.Len(price.Raw["BTC"], 2)
	test.Contains(price.Raw["BTC"], "USD")
	test.Contains(price.Raw["BTC"], "EUR")

	test.Contains(price.Display, "BTC")

	test.Len(price.Display["BTC"], 2)
	test.Contains(price.Display["BTC"], "USD")
	test.Contains(price.Display["BTC"], "EUR")
}

func TestClient_GetPriceList_ReturnsError(t *testing.T) {
	test := assert.New(t)

	client, err := New("testing")
	if !test.NoError(err) {
		test.FailNow("should be able to initialize an instance of client")
	}

	test.NotNil(client)

	price, err := client.GetPriceList([]string{"blah"}, []string{"blah"})
	test.Nil(price)
	test.Error(err)

	test.Contains(err.Error(), "market does not exist ")
}
