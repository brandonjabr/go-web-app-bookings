package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/brandonjabr/go-web-app-bookings/internal/driver"
	"github.com/brandonjabr/go-web-app-bookings/internal/handlers"
	"github.com/brandonjabr/go-web-app-bookings/internal/helpers"
	"github.com/brandonjabr/go-web-app-bookings/internal/models"
	"github.com/brandonjabr/go-web-app-bookings/internal/render"
)

const PORT = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	serve := &http.Server{
		Addr:    PORT,
		Handler: routes(&appConfig),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})

	appConfig.Production = false

	appConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.Production

	appConfig.Session = session

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=brandonjabr password=")
	if err != nil {
		log.Fatal("cannot connect to database")
		return nil, err
	}

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("could not create template cache")
		return nil, err
	}

	appConfig.TemplateCache = templateCache
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&appConfig)

	helpers.NewHelpers(&appConfig)

	return db, nil
}
