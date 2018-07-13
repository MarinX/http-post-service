package main

import (
	"sync"
	"testing"
)

func TestService(t *testing.T) {
	service, err := NewService("", "")
	if err != nil {
		t.Error(err)
		return
	}

	req := &Request{URL: "https://google.com"}
	resp := service.handleTask(req, nil)
	if resp == nil {
		t.Error("Handler return nil")
		return
	}

	if resp.Success == nil {
		t.Error("Service did not return success response")
		return
	}
	t.Log(resp.Success.StatusCode)

	req.URL = ""
	resp = service.handleTask(req, nil)
	if resp == nil {
		t.Error("Handler return nil")
		return
	}

	if resp.Error == nil {
		t.Error("Service did not return error response")
		return
	}
	t.Log(resp.Error.Message)

	var responses []*Response
	wg := new(sync.WaitGroup)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			responses = append(responses, service.handleTask(req, wg))
		}(wg)
	}
	// wait for finish
	wg.Wait()
	if len(responses) < 5 {
		t.Error("Service did not return exatch batch, want 5 got", len(responses))
	}

	for _, val := range responses {
		if val.Error == nil {
			t.Error("Service should return error, got nil insted")
		}
	}

}
