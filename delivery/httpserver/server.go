package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	Addr string
	Port int
}

func New(addr string, port int) *Server {
	return &Server{
		Addr: addr,
		Port: port,
	}
}

func (s *Server) Serve(ctx context.Context) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/health_check", s.healthCheck)

	fmt.Println("server is running on", s.Addr, s.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.Addr, s.Port), handler)

}

func WriteJson(w http.ResponseWriter, data map[string]any, status int) {
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
