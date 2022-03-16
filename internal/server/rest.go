package server

import (
	"errors"
	"net/http"
	"strings"
)

var (
	errFsymsEmpty = errors.New("fsyms param is empty")
	errTsymsEmpty = errors.New("tsyms param is empty")
)

func (server *Server) handleREST(
	response http.ResponseWriter,
	request *http.Request,
) {
	fsyms := strings.Split(request.URL.Query().Get("fsyms"), ",")
	tsyms := strings.Split(request.URL.Query().Get("tsyms"), ",")

	err := server.process(response, fsyms, tsyms)
	if err != nil {
		writeErrorJSON(response, err)
		return
	}
}
