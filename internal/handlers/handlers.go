package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/brandonjabr/go-web-app-bookings/internal/models"
	"github.com/brandonjabr/go-web-app-bookings/internal/render"
)

type Repository struct {
	AppConfig *config.AppConfig
}

var Repo *Repository

func NewRepo(appConfig *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: appConfig,
	}
}

func NewHandlers(repo *Repository) {
	Repo = repo
}

func (repo *Repository) Home(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "home.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "about.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) LuxurySuite(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "luxury-suite.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) StandardRoom(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "standard-room.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "contact.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) SearchAvailability(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "search-availability.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, req *http.Request) {
	checkInDate := req.Form.Get("checkindate")
	checkOutDate := req.Form.Get("checkoutdate")

	w.Write([]byte(fmt.Sprintf("Posted to search availability - check in date is %s and check out date is %s", checkInDate, checkOutDate)))
}

type JSONResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) AvailabilityJSON(w http.ResponseWriter, req *http.Request) {
	response := JSONResponse{
		OK:      true,
		Message: "Room is available!",
	}

	out, err := json.MarshalIndent(response, "", "		")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *Repository) Reservation(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "reservation.page.html.tmpl", &models.TemplateData{})
}
