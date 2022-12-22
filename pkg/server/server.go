package server

import (
	"context"
	"log"
	"net"
	"net/http"
)

type HTTPServer struct {
	Address string
	Context *context.Context
	Cancel  context.CancelFunc
	Server  *http.Server
	Mux     *http.ServeMux
}

func NewHTTPServer(address string, ctx context.Context, cancel context.CancelFunc) *HTTPServer {
	log.Println("creating server")

	httpsv := &HTTPServer{
		Address: address,
		Context: &ctx,
		Cancel:  cancel,
		Mux:     http.NewServeMux(),
	}

	server := http.Server{
		Addr:    httpsv.Address,
		Handler: httpsv.Mux,
		BaseContext: func(l net.Listener) context.Context {
			return *httpsv.Context
		},
	}

	httpsv.Server = &server

	return httpsv
}

func (sv *HTTPServer) Start() {
	log.Println("starting server")

	err := sv.Server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	sv.Cancel()
}
