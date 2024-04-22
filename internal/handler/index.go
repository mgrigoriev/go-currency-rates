package handler

import "net/http"

func (s *AppServer) indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := s.tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
