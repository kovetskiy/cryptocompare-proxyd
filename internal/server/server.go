package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/cache"
	"github.com/kovetskiy/cryptocompare-proxyd/internal/cryptocompare"
	"github.com/reconquest/pkg/log"
)

const (
	apiPath = "/api/v1/price"
)

// Server listens for new http/websocket connections, serves the
// requests/connections and
// processes the cryptocompare data from cache or from upstream.
type Server struct {
	listenAddress string
	http          *http.Server
	websocket     *websocket.Upgrader

	cache  cache.Cache
	client cryptocompare.Client
	ttl    int
}

// New instance of Server.
func New(
	listenAddress string,
	cache cache.Cache,
	client cryptocompare.Client,
	ttl int,
) (*Server, error) {
	return &Server{
		listenAddress: listenAddress,
		cache:         cache,
		client:        client,
		ttl:           ttl,
	}, nil
}

// ListenAndServe listens and serves received http connections.
func (server *Server) ListenAndServe() error {
	server.http = &http.Server{
		Handler: server,
		Addr:    server.listenAddress,
	}

	server.websocket = &websocket.Upgrader{
		ReadBufferSize:  1,
		WriteBufferSize: 1,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	log.Infof(nil, "the http server starting at %s", server.listenAddress)

	return server.http.ListenAndServe()
}

// Close immediately closes all active http connections.
func (server *Server) Close() error {
	return server.http.Close()
}

// ServeHTTP is invoked by net/http package when gets a connection from
// net/http.Server.
func (server *Server) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request,
) {
	ip := request.RemoteAddr

	// this might happen in case of a docker container under a nginx reverse
	// proxy where we will have X-Forwarded-For with a real ip and the IP
	// address used to send data is an internal ip
	if request.Header.Get("X-Forwarded-For") != "" {
		ip = request.Header.Get("X-Forwarded-For")
	}

	hasQuery := len(request.URL.Query()) > 0

	tag := "REST"
	if !hasQuery {
		tag = "WEBSOCKET"
	}

	log.Debugf(
		nil,
		"%10s\t%15s\t%s\t%4s\t%s",
		tag,
		ip,
		request.Header.Get("User-Agent"),
		request.Method,
		request.URL.String(),
	)

	switch {
	case request.URL.Path == apiPath && hasQuery:
		server.handleREST(response, request)

	case request.URL.Path == apiPath && !hasQuery:
		server.handleWebsocket(response, request)

	default:
		response.WriteHeader(http.StatusNotFound)
	}
}
