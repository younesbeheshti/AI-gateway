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
	Usage Usage
}


type Usage struct {
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
}