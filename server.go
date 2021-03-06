package gobaserver

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/gomuddle/goba"
	"github.com/valyala/fasthttp"
)

// A Server provides a web interface for interacting with a goba.Goba.
type Server struct {
	goba        *goba.Goba
	router      *router.Router
	credentials []Credentials
}

// New creates a new configured Server and returns it.
func New(goba goba.Goba, validCredentials ...Credentials) *Server {
	s := Server{
		goba:        &goba,
		router:      newRouter(),
		credentials: validCredentials,
	}
	s.routes()
	return &s
}

// ListenAndServe  serves HTTP requests from the given TCP address.
func (s Server) ListenAndServe(addr string) error {
	return s.fasthttpServer().ListenAndServe(addr)
}

// ListenAndServeTLS serves HTTP requests from the given TCP address.
// certFile and keyFile are paths to TLS certificate and key files.
func (s Server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	return s.fasthttpServer().ListenAndServeTLS(addr, certFile, keyFile)
}

// fasthttpServer returns a configured fasthttp.Server.
func (s Server) fasthttpServer() *fasthttp.Server {
	server := &fasthttp.Server{
		Handler: s.router.Handler,
	}
	server.ReadBufferSize = 8192
	return server
}

// writeJSON sets the response's "Content-Type" header to
// "application/json" and writes the JSON representation of
// v to the response's body.
func (s Server) writeJSON(ctx *fasthttp.RequestCtx, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ctx.SetContentType("application/json")
	_, err = ctx.Write(data)
	return nil
}
