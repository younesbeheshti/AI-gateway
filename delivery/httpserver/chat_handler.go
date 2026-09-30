package httpserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/yonesbeheshti/ai-gateway/param"
)

func (s *Server) chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJson(w, map[string]any{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	var req param.ChatRequest
	if err := json.Unmarshal(data, &req); err != nil {
		WriteJson(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	log.Println("req:", req)

	resp, err := s.aiSvc.GenerateResponse(ctx, &req)
	if err != nil {
		WriteJson(w, map[string]any{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, resp, http.StatusOK)

}
