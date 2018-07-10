package main

import "testing"

func TestService(t *testing.T) {
	service, err := NewService("", "")
	if err != nil {
		t.Error(err)
		return
	}

	resp := service.handleTask(`{"url":"https://google.com", "body":""}`)
	if resp == nil {
		t.Error("Handler return nil")
		return
	}

	if resp.Success == nil {
		t.Error("Service did not return success response")
		return
	}
	t.Log(resp.Success.StatusCode)

	resp = service.handleTask(`{"body":""}`)
	if resp == nil {
		t.Error("Handler return nil")
		return
	}

	if resp.Error == nil {
		t.Error("Service did not return error response")
		return
	}
	t.Log(resp.Error.Message)

}
