package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	req := httptest.NewRequest("POST", "/test-url", nil)
	form := New(req.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid form, should be valid")
	}
}

func TestForm_Required(t *testing.T) {
	req := httptest.NewRequest("POST", "/test-url", nil)
	form := New(req.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form was valid when required fields are missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	req = httptest.NewRequest("POST", "/test-url", nil)

	req.PostForm = postData
	form = New(req.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("got missing fields error when required fields are present")
	}
}

func TestForm_MinLength(t *testing.T) {
	req := httptest.NewRequest("POST", "/test-url", nil)

	postData := url.Values{}
	postData.Add("email", "test")

	req.PostForm = postData
	form := New(postData)

	form.MinLength("email", 5)

	if form.Valid() {
		t.Error("form had no errors when email field was too short")
	}

	postData = url.Values{}
	postData.Add("email", "test@gmail.com")

	req.PostForm = postData
	form = New(postData)

	form.MinLength("email", 5)

	if !form.Valid() {
		t.Error("form had min length error when email field was long enough")
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	testGet := form.Has("notfound")
	if testGet {
		t.Error("found field that is not in form")
	}

	postData = url.Values{}
	postData.Add("email", "test")

	form = New(postData)

	testGet = form.Has("email")
	if !testGet {
		t.Error("could not get email field present in form")
	}
}

func TestForm_GetErrors(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.Has("notfound")

	if form.Errors.Get("notfound") == "" {
		t.Error("should have an error for field, but did not")
	}

	postData = url.Values{}
	postData.Add("email", "test")

	form = New(postData)

	form.Has("email")

	if form.Errors.Get("email") != "" {
		t.Error("had error for field that should not have an error")
	}
}

func TestForm_IsValidEmail(t *testing.T) {
	postData := url.Values{}
	postData.Add("email", "test")

	form := New(postData)

	form.IsValidEmail("email")

	if form.Valid() {
		t.Error("got valid email when email is not valid")
	}

	postData = url.Values{}
	postData.Add("email", "test@gmail.com")

	form = New(postData)

	form.IsValidEmail("email")

	if !form.Valid() {
		t.Error("got invalid email when email is valid")
	}
}
