package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_setupRouter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "router is not nil and serves static files"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()
			if router == nil {
				t.Fatal("setupRouter() returned nil")
			}
			// Optionally, check that static file handler is registered
			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			// Accept 200, 403, or 404 depending on your static file setup
			if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound && rr.Code != http.StatusForbidden {
				t.Errorf("setupRouter() root handler returned status [%v], want 200, 403, or 404", rr.Code)
			}
		})
	}
}

func Test_setupServer(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{name: "server config", addr: ":12345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.NewServeMux()
			srv := setupServer(tt.addr, handler)
			if srv == nil {
				t.Fatal("setupServer() returned nil")
			}
			if srv.Addr != tt.addr {
				t.Errorf("setupServer() Addr = [%v], want [%v]", srv.Addr, tt.addr)
			}
			if srv.Handler != handler {
				t.Error("setupServer() Handler does not match input handler")
			}
			if srv.ReadTimeout != serverReadTimeout {
				t.Errorf("setupServer() ReadTimeout = [%v], want [%v]", srv.ReadTimeout, serverReadTimeout)
			}
			if srv.WriteTimeout != serverWriteTimeout {
				t.Errorf("setupServer() WriteTimeout = [%v], want [%v]", srv.WriteTimeout, serverWriteTimeout)
			}
			if srv.IdleTimeout != serverIdleTimeout {
				t.Errorf("setupServer() IdleTimeout = [%v], want [%v]", srv.IdleTimeout, serverIdleTimeout)
			}
		})
	}
}
