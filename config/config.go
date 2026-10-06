package config

import (
	redisadapter "github.com/yonesbeheshti/ai-gateway/adapter/redis"
	"github.com/yonesbeheshti/ai-gateway/delivery/httpserver"
	"github.com/yonesbeheshti/ai-gateway/pkg/ratelimiter"
	"github.com/yonesbeheshti/ai-gateway/provider/openrouter"
	"github.com/yonesbeheshti/ai-gateway/repository/postgres"
	"github.com/yonesbeheshti/ai-gateway/service/aiservice"
)

type Config struct {
	Redis       redisadapter.Config `koanf:"redis"`
	Httpserver  httpserver.Config   `koanf:"httpserver"`
	OpenRouter  openrouter.Config   `koanf:"openrouter"`
	AiSvc       aiservice.Config    `koanf:"aiservice"`
	DBConfig    postgres.Config     `koanf:"dbconfig"`
	RateLimiter ratelimiter.Config  `koanf:"ratelimiter"`
}
