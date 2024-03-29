package renders

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/justinas/nosurf"
)

var appConfig *config.AppConfig
var PathToTemplate = "./templates"

func SetConfig(a *config.AppConfig) {
	appConfig = a
}

func SetupDefaultData(templateData *model.TemplateData, r *http.Request) *model.TemplateData {

	stringMap := map[string]string{}
	stringMap["csrf_token"] = nosurf.Token(r)

	errorMessage := appConfig.Session.PopString(r.Context(), "error")
	flashMessage := appConfig.Session.PopString(r.Context(), "flash")
	warningMessage := appConfig.Session.PopString(r.Context(), "warning")

	templateData.Error = errorMessage
	templateData.Flash = flashMessage
	templateData.Warning = warningMessage
	templateData.IsLogin = appConfig.Session.Exists(r.Context(), "user_id")

	templateData.StringMap = stringMap

	return templateData
}

func ServeTemplate(w http.ResponseWriter, r *http.Request, templateName string, templateData *model.TemplateData) error {

	var allCachedTemplateMap = map[string]*template.Template{}
	if appConfig.UseCache {
		allCachedTemplateMap = appConfig.TemplateCache
	} else {
		allCachedTemplateMap, _ = CreateTemplateCache()
	}

	currentTemplatePointer, inMap := allCachedTemplateMap[templateName]

	if !inMap {
		log.Println("could not get the map")
		return errors.New("could not get the map")
	}

	SetupDefaultData(templateData, r)
	err := currentTemplatePointer.Execute(w, templateData)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	cachedTemplateMap := map[string]*template.Template{}
	templatePageDir := fmt.Sprintf("%s/*.page.tmpl", PathToTemplate)
	layoutPageDir := fmt.Sprintf("%s/*.layout.tmpl", PathToTemplate)

	pages, err := filepath.Glob(templatePageDir)

	if err != nil {
		return cachedTemplateMap, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		templatePointer, err := template.New(name).ParseFiles(page)

		if err != nil {
			return cachedTemplateMap, err
		}

		anyLayouts, err := filepath.Glob(layoutPageDir)

		if err != nil {
			return cachedTemplateMap, err
		}

		if len(anyLayouts) > 0 {

			_, err := templatePointer.ParseGlob(layoutPageDir)
			if err != nil {
				log.Println(err)
			}

		}

		cachedTemplateMap[name] = templatePointer

	}

	return cachedTemplateMap, nil

}
