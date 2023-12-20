package render

import (
	"net/http"
	"testing"

	"github.com/brandonjabr/go-web-app-bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var tmplData models.TemplateData

	req, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(req.Context(), "flash", "111")
	result := AddDefaultData(&tmplData, req)

	if result.Flash != "111" {
		t.Error("Test request failed")
	}
}

func TestRenderTemplate(t *testing.T) {
	TEMPLATE_PATH = "../../templates"

	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	appConfig.TemplateCache = templateCache

	req, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var tw testHTTPWriter

	err = Template(&tw, req, "home.page.html.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error(err)
	}

	err = Template(&tw, req, "nonexistant.page.html.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that does not exist.")
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(appConfig)
}

func TestCreateTemplateCache(t *testing.T) {
	TEMPLATE_PATH = "../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	req, err := http.NewRequest("GET", "/test-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := req.Context()
	ctx, _ = session.Load(ctx, req.Header.Get("X-Session"))
	req = req.WithContext(ctx)

	return req, nil
}
