package handlers

import (
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
	"net/http"
)

func (m *HandlerRepository) AdminDashboard(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	_ = renders.ServeTemplate(w, r, "admin.dashboard.page.tmpl", initializedTemplate)

}
