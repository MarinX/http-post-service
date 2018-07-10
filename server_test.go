package main

import (
	"net/http"
	"testing"
	"time"
)

const (
	httpPort = ":8181"
	url      = "http://localhost" + httpPort
)

func getServer() *Server {
	return NewServer(httpPort, nil)
}
func TestServer(t *testing.T) {
	go getServer().ListenAndServe()
	time.Sleep(time.Second)

	// should fail
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Server return wrong HTTP code, want 400 got", resp.StatusCode)
		return

	}

	resp, err = http.Post(url, "application/json", nil)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("Server return wrong HTTP code, want 200 got", resp.StatusCode)
	}
}
