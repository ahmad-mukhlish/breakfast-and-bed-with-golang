package helper

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
)

var appConfig *config.AppConfig

func SetHelperAppConfig(a *config.AppConfig) {
	appConfig = a
}

func CatchClientError(w http.ResponseWriter, statusCode int) {

	appConfig.ErrorLog.Println("Client Error With Status Code", statusCode)
	http.Error(w, http.StatusText(statusCode), statusCode)

}

func CatchServerError(w http.ResponseWriter, err error) {
	stackTrace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	appConfig.ErrorLog.Println(stackTrace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ConvertStringSQLToTime(timeString, dateFormat string) (time.Time, error) {

	timeResult, err := time.Parse(dateFormat, timeString)
	if err != nil {
		return time.Now(), err
	}

	return timeResult, nil

}

func IsAuthenticated(r *http.Request) bool {
	return appConfig.Session.Exists(r.Context(), "user_id")
}
