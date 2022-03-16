package server

import (
	"github.com/kovetskiy/cryptocompare-proxyd/internal/cache"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/cryptocompare"
)

type pair struct {
	fsym string
	tsym string
}

func hasRawPrice(
	list *cryptocompare.PriceList,
	fsym string,
	tsym string,
) bool {
	if _, ok := list.Raw[fsym]; !ok {
		return false
	}

	if _, ok := list.Raw[fsym][tsym]; !ok {
		return false
	}

	return true
}

func hasDisplayPrice(
	list *cryptocompare.PriceList,
	fsym string,
	tsym string,
) bool {
	if _, ok := list.Display[fsym]; !ok {
		return false
	}

	if _, ok := list.Display[fsym][tsym]; !ok {
		return false
	}

	return true
}

func newPriceList(entities []cache.Entity) *cryptocompare.PriceList {
	list := &cryptocompare.PriceList{
		Raw:     map[string]map[string]cryptocompare.RawPrice{},
		Display: map[string]map[string]cryptocompare.DisplayPrice{},
	}

	for _, entity := range entities {
		if _, ok := list.Raw[entity.FromSymbol()]; !ok {
			list.Raw[entity.FromSymbol()] = map[string]cryptocompare.RawPrice{}
		}

		if _, ok := list.Raw[entity.FromSymbol()][entity.ToSymbol()]; !ok {
			list.Raw[entity.FromSymbol()][entity.ToSymbol()] = entity.RawPrice()
		}

		if _, ok := list.Display[entity.FromSymbol()]; !ok {
			list.Display[entity.FromSymbol()] = map[string]cryptocompare.DisplayPrice{}
		}

		if _, ok := list.Display[entity.FromSymbol()][entity.ToSymbol()]; !ok {
			list.Display[entity.FromSymbol()][entity.ToSymbol()] = entity.DisplayPrice()
		}
	}

	return list
}
