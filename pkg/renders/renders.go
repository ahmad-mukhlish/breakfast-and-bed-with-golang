package renders

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ahmad-mukhlish/annahdloh-landing-page/pkg/config"
	"github.com/ahmad-mukhlish/annahdloh-landing-page/pkg/model"
)

var app *config.AppConfig

func SetConfig(a *config.AppConfig) {
	app = a
}

func setUpDefaultData(templateData *model.TemplateData) *model.TemplateData {
	return templateData
}

func ServeTemplate(w http.ResponseWriter, templateName string, templateData *model.TemplateData) {

	var allCachedTemplateMap = map[string]*template.Template{}
	if app.UseCache {
		allCachedTemplateMap = app.TemplateCache
	} else {
		allCachedTemplateMap, _ = CreateTemplateCache()
	}

	currentTemplatePointer, inMap := allCachedTemplateMap[templateName]

	if !inMap {
		log.Println("could not get the map")
	}

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

			templatePointer.ParseGlob(layoutPageDir)

		}

		cachedTemplateMap[name] = templatePointer

	}

	return cachedTemplateMap, nil

}
