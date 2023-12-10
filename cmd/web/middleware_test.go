package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var th testHandler
	h := NoSurf(&th)

	switch v := h.(type) {
	case http.Handler:
	default:
		t.Errorf("type is not http.Handler - instead got %T", v)
	}

}

func TestSessionLoad(t *testing.T) {

}
