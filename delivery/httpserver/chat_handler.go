package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yonesbeheshti/ai-gateway/param"
)

func (s *Server) chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJson(w, map[string]any{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB limit
	var req param.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJson(w, map[string]any{"error": "invalid request body"}, http.StatusBadRequest)
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
