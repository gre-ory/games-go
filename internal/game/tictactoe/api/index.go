package api

import (
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.tpl", nil)
}
