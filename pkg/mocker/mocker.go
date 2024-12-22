package mocker

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
)

type Server struct {
	mux               *http.ServeMux
	server            *httptest.Server
	port              int
	basicAuth         bool
	basicAuthUser     string
	basicAuthPassword string
}

type H map[string]any

type HandlerFunc func(*Context)

func New() *Server {
	return NewWithPort(0)
}

func NewWithPort(port int) *Server {
	return &Server{mux: http.NewServeMux(), port: port}
}

func (s *Server) SetBasicAuth(username string, password string) {
	s.basicAuthUser = username
	s.basicAuthPassword = password
	s.basicAuth = true
}

func (s *Server) GET(pattern string, handler HandlerFunc) {
	s.addRoute(pattern, handler)
}

func (s *Server) addRoute(pattern string, handler HandlerFunc) {
	f := func(w http.ResponseWriter, r *http.Request) {
		if s.basicAuth && !s.verifyBasicAuth(w, r) {
			return
		}
		handler(&Context{Writer: w, Request: r})
	}
	s.mux.HandleFunc(pattern, f)
}

func (s *Server) verifyBasicAuth(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok || username != s.basicAuthUser || password != s.basicAuthPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="mocker"`)
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func (s *Server) Start() error {
	s.server = httptest.NewUnstartedServer(s.mux)
	if s.port != 0 {
		s.server.Listener, _ = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	}
	s.server.Start()
	return nil
}

func (s *Server) Close() {
	s.server.Close()
}

func (s *Server) URL() string {
	return s.server.URL
}
