package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var handlersTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"GET home", "/", "GET", []postData{}, http.StatusOK},
	{"GET contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"GET about", "/about", "GET", []postData{}, http.StatusOK},
	{"GET luxury suite", "/rooms/luxury-suite", "GET", []postData{}, http.StatusOK},
	{"GET standard room", "/rooms/standard-room", "GET", []postData{}, http.StatusOK},
	{"GET search availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"GET reservation", "/reservation", "GET", []postData{}, http.StatusOK},
	{"GET reservation details", "/reservation-details", "GET", []postData{}, http.StatusOK},
	{"POST search availability", "/reservation", "POST", []postData{
		{key: "check_in_date", value: "01-01-2024"},
		{key: "check_out_date", value: "01-05-2024"},
	}, http.StatusOK},
	{"POST search availability JSON", "/reservation", "POST", []postData{
		{key: "check_in_date", value: "01-01-2024"},
		{key: "check_out_date", value: "01-05-2024"},
	}, http.StatusOK},
	{"POST reservation", "/reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Johnson"},
		{key: "email", value: "test@gmail.com"},
		{key: "phone_number", value: "111-111-1111"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, handlersTest := range handlersTests {
		if handlersTest.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + handlersTest.url)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != handlersTest.expectedStatusCode {
				t.Errorf("for test %s, expected status code %d but got %d", handlersTest.name, handlersTest.expectedStatusCode, resp.StatusCode)
			}
		} else {
			vals := url.Values{}
			for _, x := range handlersTest.params {
				vals.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+handlersTest.url, vals)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != handlersTest.expectedStatusCode {
				t.Errorf("for test %s, expected status code %d but got %d", handlersTest.name, handlersTest.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
