package param

import "github.com/yonesbeheshti/ai-gateway/entity"

type CachePayload struct {
	Model    string           `json:"model"`
	Messages []entity.Message `json:"messages"`
}
