package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/brandonjabr/go-web-app-bookings/internal/models"
)

var session *scs.SessionManager
var testAppConfig config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	testAppConfig.Production = false

	testAppConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testAppConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testAppConfig.Session = session

	appConfig = &testAppConfig

	os.Exit(m.Run())
}

type testHTTPWriter struct{}

func (tw *testHTTPWriter) Header() http.Header {
	var header http.Header
	return header
}

func (tw *testHTTPWriter) WriteHeader(i int) {

}

func (tw *testHTTPWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
