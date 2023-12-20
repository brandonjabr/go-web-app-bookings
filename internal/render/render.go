package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/brandonjabr/go-web-app-bookings/internal/models"
	"github.com/justinas/nosurf"
)

var TEMPLATE_PATH = "./templates"

var appConfig *config.AppConfig
var functions = template.FuncMap{}

func NewRenderer(c *config.AppConfig) {
	appConfig = c
}

func AddDefaultData(tmplData *models.TemplateData, req *http.Request) *models.TemplateData {
	tmplData.Flash = appConfig.Session.PopString(req.Context(), "flash")
	tmplData.Error = appConfig.Session.PopString(req.Context(), "error")
	tmplData.Warning = appConfig.Session.PopString(req.Context(), "warning")
	tmplData.CSRFToken = nosurf.Token(req)
	return tmplData
}

func Template(w http.ResponseWriter, req *http.Request, tmpl string, tmplData *models.TemplateData) error {
	var tmplCache map[string]*template.Template

	if appConfig.UseCache {
		tmplCache = appConfig.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	template, inCache := tmplCache[tmpl]
	if !inCache {
		log.Println("template not found in cache")
		return errors.New("template not found in cache")
	}

	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, req)

	_ = template.Execute(buf, tmplData)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html.tmpl", TEMPLATE_PATH))
	if err != nil {
		return tmplCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		matchingFiles, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html.tmpl", TEMPLATE_PATH))
		if err != nil {
			return tmplCache, err
		}

		if len(matchingFiles) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.layout.html.tmpl", TEMPLATE_PATH))
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = tmplSet
	}

	return tmplCache, nil
}
