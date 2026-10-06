package httpserver

import "net/http"

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: get userId from token/context
		userId := int64(5)
		if !s.rateLimiter.Allow(userId) {
			WriteJson(w, map[string]any{"error": "rate limit exceeded"}, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
