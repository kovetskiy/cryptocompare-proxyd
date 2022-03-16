package server

import (
	"encoding/json"
	"io"

	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

func writeJSON(writer io.Writer, msg interface{}) {
	err := json.NewEncoder(writer).Encode(msg)
	if err != nil {
		log.Errorf(err, "server: write json")
	}
}

func writeErrorJSON(writer io.Writer, err error) {
	log.Error(err)

	response := struct {
		Error string `json:"error"`
	}{}

	if karmic, ok := err.(karma.Karma); ok {
		response.Error = karmic.GetMessage()
	} else {
		response.Error = err.Error()
	}

	writeJSON(writer, response)
}
