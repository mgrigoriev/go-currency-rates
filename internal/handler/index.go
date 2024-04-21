package handler

import "net/http"

func (s *HTTPServer) indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := s.tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
