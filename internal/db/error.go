package forum

import (
	"net/http"
	"text/template"
)

// Set error message for given status code
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	var message string
	template, err := template.ParseFiles(
		"web/templates/error.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch status {
	case http.StatusBadRequest:
		message = http.StatusText(status)
	case http.StatusInternalServerError:
		message = http.StatusText(status)
	case http.StatusNotFound:
		message = http.StatusText(status)
	}
	if err := template.Execute(w, map[string]string{"ErrorMessage": message}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
