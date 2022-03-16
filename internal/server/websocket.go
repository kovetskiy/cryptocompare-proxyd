package server

import (
	"encoding/json"
	"net/http"

	"github.com/reconquest/karma-go"
)

type websocketQuery struct {
	Fsyms []string `json:"fsyms"`
	Tsyms []string `json:"tsyms"`
}

func (server *Server) handleWebsocket(
	response http.ResponseWriter,
	request *http.Request,
) {
	connection, err := server.websocket.Upgrade(response, request, nil)
	if err != nil {
		writeErrorJSON(response, err)
		return
	}

	defer connection.Close()

	wsWriter := websocketWriter{connection: connection}

	for {
		_, reader, err := connection.NextReader()
		if err != nil {
			break
		}

		var query websocketQuery
		err = json.NewDecoder(reader).Decode(&query)
		if err != nil {
			writeErrorJSON(
				wsWriter,
				karma.Format(err, "json decoding failed"),
			)

			return
		}

		err = server.process(wsWriter, query.Fsyms, query.Tsyms)
		if err != nil {
			writeErrorJSON(wsWriter, err)
		}
	}
}
