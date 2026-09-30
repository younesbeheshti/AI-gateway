package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yonesbeheshti/ai-gateway/service/aiservice"
)

type Config struct {
	Addr string
	Port int
}

type Server struct {
	cfg   Config
	aiSvc aiservice.Service
}

func New(cfg Config, aiSvc aiservice.Service) *Server {
	return &Server{
		cfg:   cfg,
		aiSvc: aiSvc,
	}
}

func (s *Server) Serve(_ context.Context) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/health_check", s.healthCheck)
	handler.HandleFunc("/api/chat", s.chatHandler)

	fmt.Println("server is running on", s.cfg.Addr, s.cfg.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.Port), handler)

}

func WriteJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal server error"}`))
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
