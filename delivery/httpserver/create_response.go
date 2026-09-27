package httpserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yonesbeheshti/ai-gateway/param"
)

func (s *Server) createResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJson(w, map[string]any{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	data, err :=io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
	}
	var req *param.ChatRequest
	if err := json.Unmarshal(data,req); err != nil {
		WriteJson(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
	}





}
