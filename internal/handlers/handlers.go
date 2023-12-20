package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/brandonjabr/go-web-app-bookings/internal/driver"
	"github.com/brandonjabr/go-web-app-bookings/internal/format"
	"github.com/brandonjabr/go-web-app-bookings/internal/forms"
	"github.com/brandonjabr/go-web-app-bookings/internal/helpers"
	"github.com/brandonjabr/go-web-app-bookings/internal/models"
	"github.com/brandonjabr/go-web-app-bookings/internal/render"
	"github.com/brandonjabr/go-web-app-bookings/internal/repository"
	"github.com/brandonjabr/go-web-app-bookings/internal/repository/db_repo"
)

type Repository struct {
	AppConfig *config.AppConfig
	DB        repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(appConfig *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		AppConfig: appConfig,
		DB:        db_repo.NewPostgresRepo(db.SQL, appConfig),
	}
}

func NewHandlers(repo *Repository) {
	Repo = repo
}

func (repo *Repository) Home(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "home.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "about.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) LuxurySuite(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "luxury-suite.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) StandardRoom(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "standard-room.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "contact.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) SearchAvailability(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "search-availability.page.html.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, req *http.Request) {
	checkInDate := req.Form.Get("check_in_date")
	checkOutDate := req.Form.Get("check_in_date")

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
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *Repository) Reservation(w http.ResponseWriter, req *http.Request) {
	var emptyReservation models.Reservation
	reservationData := make(map[string]interface{})
	reservationData["reservation"] = emptyReservation
	render.Template(w, req, "reservation.page.html.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		OtherData: reservationData,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	checkInDate, err := format.ParseDate(req.Form.Get("check_in_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	checkOutDate, err := format.ParseDate(req.Form.Get("check_out_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(req.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName:    req.Form.Get("first_name"),
		LastName:     req.Form.Get("last_name"),
		Email:        req.Form.Get("email"),
		PhoneNumber:  req.Form.Get("phone_number"),
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
		RoomID:       roomID,
	}

	form := forms.New(req.PostForm)

	form.Required("first_name", "last_name", "email", "check_in_date", "check_out_date", "room_id")

	form.IsValidEmail("email")

	if !form.Valid() {
		reservationData := make(map[string]interface{})
		reservationData["reservation"] = reservation

		render.Template(w, req, "reservation.page.html.tmpl", &models.TemplateData{
			Form:      form,
			OtherData: reservationData,
		})
		return
	}

	newReservationID, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomRestriction := models.RoomRestriction{
		StartDate:     checkInDate,
		EndDate:       checkOutDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = repo.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	repo.AppConfig.Session.Put(req.Context(), "reservation", reservation)

	http.Redirect(w, req, "/reservation-details", http.StatusSeeOther)
}

func (repo *Repository) ReservationDetails(w http.ResponseWriter, req *http.Request) {
	reservation, ok := repo.AppConfig.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.AppConfig.ErrorLog.Println("can't get reservation from session")
		repo.AppConfig.Session.Put(req.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	repo.AppConfig.Session.Remove(req.Context(), "reservation")

	reservationData := make(map[string]interface{})
	reservationData["reservation"] = reservation

	render.Template(w, req, "reservation-details.page.html.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		OtherData: reservationData,
	})
}
