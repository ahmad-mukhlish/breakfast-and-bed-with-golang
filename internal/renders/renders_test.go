package renders

import (
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"net/http"
	"testing"
)

func TestSetupDefaultData(t *testing.T) {

	r, err := createRequestWithSession()

	dummyFlash := "testo 123"

	mockedSession.Put(r.Context(), "flash", dummyFlash)

	if err != nil {
		t.Error(err)
	}

	testedDefaultData := SetupDefaultData(&model.TemplateData{}, r)

	if testedDefaultData.Flash != dummyFlash {
		t.Errorf("Expected %s, but got %s", dummyFlash, testedDefaultData.Flash)
	}

}

func TestCreateTemplateCache(t *testing.T) {
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}

func TestSetConfig(t *testing.T) {
	SetConfig(appConfig)
}

func TestServeTemplate(t *testing.T) {

	r, err := createRequestWithSession()

	if err != nil {
		t.Error(err)
	}

	PathToTemplate = "./../../templates"
	mw := mockedWriter{}
	err = ServeTemplate(&mw, r, "home.page.tmpl", &model.TemplateData{})

	if err != nil {
		t.Error(err)
	}

	err = ServeTemplate(&mw, r, "non-exist.page.tmpl", &model.TemplateData{})

	if err == nil {
		t.Error("this page should not exist, therefore the error is expected")
	}
}

func createRequestWithSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/test", nil)

	if err != nil {
		return nil, err
	}

	currentContext := r.Context()

	sessionedContext, _ := mockedSession.Load(currentContext, r.Header.Get("X-Session"))

	return r.WithContext(sessionedContext), nil

}
