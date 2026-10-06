package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yonesbeheshti/ai-gateway/param"
)

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJson(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	var req param.LoginRequest
	if err := json.Unmarshal(data, &req); err != nil {
		WriteJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	fmt.Println("req:", req)

	// send to user service and generate user token with that has amount of token to use
	resp, err := s.userSvc.Login(ctx, &req)
	if err != nil {
		WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, resp, http.StatusOK)
}
