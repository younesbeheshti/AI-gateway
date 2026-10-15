package httpserver

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/yonesbeheshti/ai-gateway/entity"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteJson(w, map[string]string{"error": "missing Authorization header"}, http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer sk-gw-..."
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")
		if apiKey == authHeader || apiKey == "" {
			WriteJson(w, map[string]string{"error": "invalid Authorization format, use Bearer <token>"}, http.StatusUnauthorized)
			return
		}

		user, err := s.userSvc.GetUserByApiKey(r.Context(), apiKey)
		if err != nil || user == nil {
			WriteJson(w, map[string]string{"error": "unauthorized: invalid api key"}, http.StatusUnauthorized)
			return
		}

		// Optional: check simulated token balance
		//if user.TokenBalance <= 0 {
		//	WriteJson(w, map[string]string{"error": "insufficient token balance"}, http.StatusPaymentRequired)
		//	return
		//}

		log.Println("authenticated user:", user)

		// Pass user into context for handlers
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(userContextKey).(*entity.User)

		if !s.rateLimiter.Allow(user.ID) {
			WriteJson(w, map[string]any{"error": "rate limit exceeded"}, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
