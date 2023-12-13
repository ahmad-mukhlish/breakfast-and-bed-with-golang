package renders

import (
	"encoding/gob"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"os"
	"testing"
	"time"
)

var mockedSession *scs.SessionManager
var mockedAppConfig config.AppConfig

func TestMain(m *testing.M) {

	setupSession()
	os.Exit(m.Run())

}

func setupSession() {
	mockedSession = scs.New()

	mockedSession.Lifetime = 24 * time.Hour
	mockedSession.Cookie.Secure = false
	mockedSession.Cookie.Persist = true
	mockedSession.Cookie.SameSite = http.SameSiteLaxMode

	//register custom types here with gob.register
	gob.Register(model.Reservation{})

	mockedAppConfig.Session = mockedSession
	appConfig = &mockedAppConfig

}

type mockedWriter struct {
}

func (mw *mockedWriter) Header() http.Header {
	return http.Header{}
}

func (mw *mockedWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (mw *mockedWriter) WriteHeader(int) {
}
