package httpserver

import (
	"net/http"
)

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJson(w,
			map[string]any{"error": "method not allowed"},
			http.StatusMethodNotAllowed,
		)
		return
	}



	WriteJson(w,
		map[string]any{
			"status": "ok",
		}, http.StatusOK,
	)
}
