package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/brandonjabr/go-web-app-bookings/internal/driver"
	"github.com/brandonjabr/go-web-app-bookings/internal/format"
	"github.com/brandonjabr/go-web-app-bookings/internal/forms"
	"github.com/brandonjabr/go-web-app-bookings/internal/helpers"
	"github.com/brandonjabr/go-web-app-bookings/internal/models"
	"github.com/brandonjabr/go-web-app-bookings/internal/render"
	"github.com/brandonjabr/go-web-app-bookings/internal/repository"
	"github.com/brandonjabr/go-web-app-bookings/internal/repository/db_repo"
	"github.com/go-chi/chi"
)

type Repository struct {
	AppConfig *config.AppConfig
	DB        repository.DatabaseRepo
}

type JSONResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
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
	checkInDate, err := format.ParseDate(strings.Replace(req.Form.Get("check_in_date"), "/", "-", -1))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	checkOutDate, err := format.ParseDate(strings.Replace(req.Form.Get("check_out_date"), "/", "-", -1))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(checkInDate, checkOutDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		repo.AppConfig.Session.Put(req.Context(), "error", "No availability for searched dates")
		http.Redirect(w, req, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := models.Reservation{
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
		RoomID:       0,
	}

	repo.AppConfig.Session.Put(req.Context(), "reservation", reservation)

	render.Template(w, req, "select-room.page.html.tmpl", &models.TemplateData{
		OtherData: data,
	})
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
	reservation, ok := repo.AppConfig.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
	}

	checkInDate, err := format.ParseDateToString(reservation.CheckInDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	checkOutDate, err := format.ParseDateToString(reservation.CheckOutDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	dateStringMap := make(map[string]string)
	dateStringMap["check_in_date"] = checkInDate
	dateStringMap["check_out_date"] = checkOutDate

	reservationData := make(map[string]interface{})
	reservationData["reservation"] = reservation

	render.Template(w, req, "reservation.page.html.tmpl", &models.TemplateData{
		Form:       forms.New(nil),
		StringData: dateStringMap,
		OtherData:  reservationData,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	checkInDate, err := format.ParseStringToDate(req.Form.Get("check_in_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	checkOutDate, err := format.ParseStringToDate(req.Form.Get("check_out_date"))
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

func (repo *Repository) SelectRoom(w http.ResponseWriter, req *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := repo.AppConfig.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	reservation.RoomID = roomID

	repo.AppConfig.Session.Put(req.Context(), "reservation", reservation)

	http.Redirect(w, req, "/reservation", http.StatusSeeOther)
}
