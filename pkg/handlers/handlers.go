package handlers

import (
	"context"
	"net/http"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/model"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/renders"
)

const IPAddressKey = "ip_address"

type Repository struct {
	AppConfig *config.AppConfig
}

var Repo *Repository

func CreateRepository(appConfig *config.AppConfig) *Repository {

	return &Repository{
		AppConfig: appConfig,
	}

}

func CreateHandlers(repository *Repository) {
	Repo = repository
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, "about.page.tmpl", initializedTempalte)
}

func (repo *Repository) Learn(w http.ResponseWriter, r *http.Request) {

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, "learn.page.tmpl", initializedTempalte)
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {

	IPAddrress := r.RemoteAddr
	repo.AppConfig.Session.Put(r.Context(), IPAddressKey, IPAddrress)

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, "home.page.tmpl", initializedTempalte)

}

func initiateTemplate(
	appConfig *config.AppConfig,
	context context.Context) *model.TemplateData {

	stringMap := map[string]string{}

	stringMap["test"] = "this is some string"
	stringMap[IPAddressKey] = appConfig.Session.GetString(context, IPAddressKey)

	templateData := model.TemplateData{
		StringMap: stringMap,
	}

	return &templateData

}
