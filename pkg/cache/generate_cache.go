package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/yonesbeheshti/ai-gateway/entity"
	"github.com/yonesbeheshti/ai-gateway/param"
)

const (
	cachePrefix = "chat:cache:"
)

func GenerateCacheKey(req *entity.ChatRequest) (string, error) {
	payload := param.CachePayload{
		Model:    req.Model,
		Messages: req.Messages,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cache key payload: %w", err)
	}
	hash := sha256.Sum256(data)
	return cachePrefix + hex.EncodeToString(hash[:]), nil
}
