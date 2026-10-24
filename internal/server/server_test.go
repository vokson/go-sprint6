package server

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInitApp(t *testing.T) {
	logger := log.New(io.Discard, "", 0)

	app := InitApp(logger)

	if app == nil {
		t.Fatal("InitApp returned nil")
	}

	if app.Logger != logger {
		t.Error("logger is not set on the server struct")
	}

	if app.Server == nil {
		t.Fatal("http.Server is nil")
	}

	if app.Server.Addr != ":8080" {
		t.Errorf("Addr = %q, want %q", app.Server.Addr, ":8080")
	}

	if app.Server.ErrorLog != logger {
		t.Error("ErrorLog is not set to the provided logger")
	}

	if app.Server.Handler == nil {
		t.Error("Handler is nil")
	}

	if app.Server.ReadTimeout != 5*time.Second {
		t.Errorf("ReadTimeout = %v, want %v", app.Server.ReadTimeout, 5*time.Second)
	}

	if app.Server.WriteTimeout != 10*time.Second {
		t.Errorf("WriteTimeout = %v, want %v", app.Server.WriteTimeout, 10*time.Second)
	}

	if app.Server.IdleTimeout != 15*time.Second {
		t.Errorf("IdleTimeout = %v, want %v", app.Server.IdleTimeout, 15*time.Second)
	}
}

func TestServerServesIndex(t *testing.T) {
	t.Chdir("../..")

	logger := log.New(io.Discard, "", 0)
	app := InitApp(logger)

	ts := httptest.NewServer(app.Server.Handler)
	t.Cleanup(ts.Close)

	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("response body is empty")
	}
}
