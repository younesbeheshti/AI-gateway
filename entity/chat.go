package entity

type ChatRequest struct {
	Model    string
	Messages []Message
	Stream   bool
}

type Message struct {
	Role    string
	Content string
}

type ChatResponse struct {
	Content string
}
