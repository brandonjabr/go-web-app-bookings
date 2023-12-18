package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/brandonjabr/go-web-app-bookings/internal/config"
)

var appConfig *config.AppConfig

func NewHelpers(ac *config.AppConfig) {
	appConfig = ac
}

func ClientError(w http.ResponseWriter, status int) {
	appConfig.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	appConfig.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
