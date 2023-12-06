package renders

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/justinas/nosurf"
)

var appConfig *config.AppConfig

func SetConfig(a *config.AppConfig) {
	appConfig = a
}

func setUpDefaultData(templateData *model.TemplateData, r *http.Request) *model.TemplateData {

	stringMap := map[string]string{}
	stringMap["csrf_token"] = nosurf.Token(r)
	templateData.StringMap = stringMap

	return templateData
}

func ServeTemplate(w http.ResponseWriter, r *http.Request, templateName string, templateData *model.TemplateData) {

	var allCachedTemplateMap = map[string]*template.Template{}
	if appConfig.UseCache {
		allCachedTemplateMap = appConfig.TemplateCache
	} else {
		allCachedTemplateMap, _ = CreateTemplateCache()
	}

	currentTemplatePointer, inMap := allCachedTemplateMap[templateName]

	if !inMap {
		log.Println("could not get the map")
	}

	setUpDefaultData(templateData, r)
	err := currentTemplatePointer.Execute(w, templateData)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	cachedTemplateMap := map[string]*template.Template{}
	templatePageDir := "./templates/*.page.tmpl"
	layoutPageDir := "./templates/*.layout.tmpl"

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
