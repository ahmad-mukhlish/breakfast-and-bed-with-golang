package config

import (
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

// AppConfig is a global config for the whole app
type AppConfig struct {
	TemplateCache    map[string]*template.Template
	UseCache         bool
	IsProductionMode bool
	Session          *scs.SessionManager
	InfoLog          *log.Logger
	ErrorLog         *log.Logger
	MailChan         chan model.MailData
}
