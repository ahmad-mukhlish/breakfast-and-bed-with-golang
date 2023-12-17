package helper

import (
	"fmt"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"net/http"
	"runtime/debug"
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
