package param

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message Message `json:"message"`
}
type ChatResponse struct {
	Message Message `json:"message"`
}
