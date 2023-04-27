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
	URL               string
	basicAuth         bool
	basicAuthUser     string
	basicAuthPassword string
}

type H map[string]any

type HandlerFunc func(*Context)

func New() *Server {
	return &Server{mux: http.NewServeMux()}
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

func (s *Server) Start(port int) error {
	s.server = httptest.NewUnstartedServer(s.mux)
	// port == 0 : random port
	if port > 1 {
		if port > 65535 {
			return fmt.Errorf("invalid port number: %d", port)
		}
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			return fmt.Errorf("error on Listen: %w", err)
		}
		s.server.Listener = listener
	}
	s.server.Start()
	s.URL = s.server.URL
	return nil
}

func (s *Server) Close() {
	s.server.Close()
}
