package main

import "net/http"

type httpServer struct {
	server *http.Server
}

const httpPortEnvKey = "HTTP_PORT"

func newHTTPServer(handler http.Handler, getenv func(string) string) *httpServer {
	return &httpServer{
		server: &http.Server{
			Addr:    ":" + getenv(httpPortEnvKey),
			Handler: handler,
		},
	}
}

func (s *httpServer) run() error {
	return s.server.ListenAndServe()
}
