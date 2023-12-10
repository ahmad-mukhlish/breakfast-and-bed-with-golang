package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
)

// AppConfig is a global config for the whole app
type AppConfig struct {
	TemplateCache    map[string]*template.Template
	UseCache         bool
	IsProductionMode bool
	Session          *scs.SessionManager
}
