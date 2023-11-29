package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

type AppConfig struct {
	TemplateCache    map[string]*template.Template
	UseCache         bool
	InfoLog          *log.Logger
	IsProductionMode bool
	Session          *scs.SessionManager
	ResRoutePath     string
}
