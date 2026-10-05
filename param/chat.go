package param

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	WithCache bool      `json:"with_cache"`
}
type ChatResponse struct {
	Content string `json:"message"`
}
