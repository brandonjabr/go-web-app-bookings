package main

import (
	"testing"

	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("type is not *chi.Mux - instead got %T", v)
	}
}
