package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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

func (s *Server) Serve(ctx context.Context) error {
	handler := http.NewServeMux()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second, // Allow enough time for LLM response
		IdleTimeout:  120 * time.Second,
	}

	// TODO: add get handler to list all the model for the user
	handler.HandleFunc("/health_check", s.healthCheck)
	handler.HandleFunc("/api/chat/completions", s.chatHandler)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server listening on %s:%d\n", s.cfg.Addr, s.cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("http server failed: %w", err)
	case <-ctx.Done():
		log.Println("Shutdown signal received, shutting down gracefully...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}
	log.Println("Server exited cleanly")
	return nil
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
