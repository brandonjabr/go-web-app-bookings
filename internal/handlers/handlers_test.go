package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brandonjabr/go-web-app-bookings/internal/format"
	"github.com/brandonjabr/go-web-app-bookings/internal/models"
)

var handlersTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"GET home", "/", "GET", http.StatusOK},
	{"GET contact", "/contact", "GET", http.StatusOK},
	{"GET about", "/about", "GET", http.StatusOK},
	{"GET luxury suite", "/rooms/luxury-suite", "GET", http.StatusOK},
	{"GET standard room", "/rooms/standard-room", "GET", http.StatusOK},
	{"GET search availability", "/search-availability", "GET", http.StatusOK},
	{"GET reservation", "/reservation", "GET", http.StatusOK},
	{"GET reservation details", "/reservation-details", "GET", http.StatusOK},
	// {"POST search availability", "/reservation", "POST", []postData{
	// 	{key: "check_in_date", value: "01-01-2024"},
	// 	{key: "check_out_date", value: "01-05-2024"},
	// }, http.StatusOK},
	// {"POST search availability JSON", "/reservation", "POST", []postData{
	// 	{key: "check_in_date", value: "01-01-2024"},
	// 	{key: "check_out_date", value: "01-05-2024"},
	// }, http.StatusOK},
	// {"POST reservation", "/reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Johnson"},
	// 	{key: "email", value: "test@gmail.com"},
	// 	{key: "phone_number", value: "111-111-1111"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, handlersTest := range handlersTests {
		resp, err := testServer.Client().Get(testServer.URL + handlersTest.url)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != handlersTest.expectedStatusCode {
			t.Errorf("for test %s, expected status code %d but got %d", handlersTest.name, handlersTest.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:   1,
			Name: "Standard Room",
		},
	}

	// Test case where reservation is valid
	req, _ := http.NewRequest("GET", "/reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	responseRecorder := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("reservation handler returned unexpected response code: got %d, expected %d", responseRecorder.Code, http.StatusOK)
	}

	// Test case where reservation not in session
	req, _ = http.NewRequest("GET", "/reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned unexpected response code: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Test case where room does not exist
	req, _ = http.NewRequest("GET", "/reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	reservation.RoomID = 100

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned unexpected response code: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	// Test case where post reservation is valid
	body := "check_in_date=05-10-2050"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=05-11-2050")
	body = fmt.Sprintf("%s&%s", body, "first_name=John")
	body = fmt.Sprintf("%s&%s", body, "last_name=Johnson")
	body = fmt.Sprintf("%s&%s", body, "email=john@smith.com")
	body = fmt.Sprintf("%s&%s", body, "phone_number=1112223333")
	body = fmt.Sprintf("%s&%s", body, "room_id=2")

	req, _ := http.NewRequest("POST", "/reservation", strings.NewReader(body))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned unexpected response code: got %d, expected %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Test case where post reservation body is missing
	req, _ = http.NewRequest("POST", "/reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned unexpected response code: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Test case where check in date is invalid
	body = "check_in_date=invalid"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=05-11-2050")
	body = fmt.Sprintf("%s&%s", body, "first_name=John")
	body = fmt.Sprintf("%s&%s", body, "last_name=Johnson")
	body = fmt.Sprintf("%s&%s", body, "email=john@smith.com")
	body = fmt.Sprintf("%s&%s", body, "phone_number=1112223333")
	body = fmt.Sprintf("%s&%s", body, "room_id=2")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned unexpected response code for invalid check in date: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Test case where check out date is invalid
	body = "check_in_date=05-10-2050"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=invalid")
	body = fmt.Sprintf("%s&%s", body, "first_name=John")
	body = fmt.Sprintf("%s&%s", body, "last_name=Johnson")
	body = fmt.Sprintf("%s&%s", body, "email=john@smith.com")
	body = fmt.Sprintf("%s&%s", body, "phone_number=1112223333")
	body = fmt.Sprintf("%s&%s", body, "room_id=2")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned unexpected response code for invalid check out date: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Test case where room ID is invalid
	body = "check_in_date=05-10-2050"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=05-11-2050")
	body = fmt.Sprintf("%s&%s", body, "first_name=John")
	body = fmt.Sprintf("%s&%s", body, "last_name=Johnson")
	body = fmt.Sprintf("%s&%s", body, "email=john@smith.com")
	body = fmt.Sprintf("%s&%s", body, "phone_number=1112223333")
	body = fmt.Sprintf("%s&%s", body, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned unexpected response code for invalid room ID: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Test case where required field is missing
	body = "check_in_date=05-10-2050"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=05-11-2050")
	body = fmt.Sprintf("%s&%s", body, "last_name=Johnson")
	body = fmt.Sprintf("%s&%s", body, "email=john@smith.com")
	body = fmt.Sprintf("%s&%s", body, "phone_number=1112223333")
	body = fmt.Sprintf("%s&%s", body, "room_id=1")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned unexpected response code for required field missing: got %d, expected %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Test case where inserting reservation fails due to invalid room ID
	body = "check_in_date=05-10-2050"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=05-11-2050")
	body = fmt.Sprintf("%s&%s", body, "first_name=John")
	body = fmt.Sprintf("%s&%s", body, "last_name=Johnson")
	body = fmt.Sprintf("%s&%s", body, "email=john@smith.com")
	body = fmt.Sprintf("%s&%s", body, "phone_number=1112223333")
	body = fmt.Sprintf("%s&%s", body, "room_id=100")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned unexpected response code for room does not exist: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ReservationDetails(t *testing.T) {
	checkInDate, _ := format.ParseStringToDate("05-10-2050")
	checkOutDate, _ := format.ParseStringToDate("05-11-2050")

	reservation := models.Reservation{
		FirstName:    "John",
		LastName:     "Johnson",
		Email:        "john@smith.com",
		PhoneNumber:  "1112223333",
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
		RoomID:       1,
		Room: models.Room{
			ID:   1,
			Name: "Standard Room",
		},
	}

	// Test case where reservation is valid
	req, _ := http.NewRequest("GET", "/reservation-details", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	responseRecorder := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationDetails)

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("reservation details handler returned unexpected response code: got %d, expected %d", responseRecorder.Code, http.StatusOK)
	}

	// Test case where reservation not in session
	req, _ = http.NewRequest("GET", "/reservation-details", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation details handler returned unexpected response code: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {

	// Test case where room is available
	body := "check_in_date=12-20-2023"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=12-21-2023")
	body = fmt.Sprintf("%s&%s", body, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(body))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(responseRecorder, req)

	var response JSONResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Error("failed to parse json response from POST to search availability")
	}

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("post search availablity handler returned unexpected response code for rooms are not available: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Test case where room is not available
	body = "check_in_date=12-24-2023"
	body = fmt.Sprintf("%s&%s", body, "check_out_date=12-24-2023")
	body = fmt.Sprintf("%s&%s", body, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(responseRecorder, req)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Error("failed to parse json response from POST to search availability")
	}

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("post search availablity handler returned unexpected response code for rooms are not available: got %d, expected %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
