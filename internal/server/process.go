package server

import (
	"context"
	"fmt"
	"io"

	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

func (server *Server) process(
	response io.Writer,
	fsyms []string,
	tsyms []string,
) error {
	if len(fsyms) == 0 {
		return errFsymsEmpty
	}

	if len(tsyms) == 0 {
		return errTsymsEmpty
	}

	entities, err := server.cache.Read(
		context.Background(),
		fsyms,
		tsyms,
		server.ttl,
	)
	if err != nil {
		return karma.Format(err, "cache: read data failed")
	}

	list := newPriceList(entities)

	missing := []pair{}
	for _, fsym := range fsyms {
		for _, tsym := range tsyms {
			if !hasRawPrice(list, fsym, tsym) ||
				!hasDisplayPrice(list, fsym, tsym) {
				// found a pair that is not in our cache
				missing = append(
					missing,
					pair{fsym: fsym, tsym: tsym},
				)
			}
		}
	}

	if len(missing) == 0 {
		writeJSON(response, list)
		return nil
	}

	// some useful list of pairs for analytics
	log.Warning(
		karma.
			Describe("pairs", fmt.Sprintf("%v", missing)).
			Format(nil, "the user requested pairs missing in the cache storage"),
	)

	upstreamList, err := server.client.GetPriceList(fsyms, tsyms)
	if err != nil {
		return karma.Format(err, "upstream: request price list failed")
	}

	writeJSON(response, upstreamList)

	return nil
}
