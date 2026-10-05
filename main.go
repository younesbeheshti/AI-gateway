package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	redisadapter "github.com/yonesbeheshti/ai-gateway/adapter/redis"
	"github.com/yonesbeheshti/ai-gateway/config"
	"github.com/yonesbeheshti/ai-gateway/delivery/httpserver"
	"github.com/yonesbeheshti/ai-gateway/provider/ollama"
	"github.com/yonesbeheshti/ai-gateway/provider/openrouter"
	"github.com/yonesbeheshti/ai-gateway/repository/redis"
	"github.com/yonesbeheshti/ai-gateway/service/aiservice"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Note: .env file not found, reading from system env")
	}

	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("config: %+v\n", cfg)

	redisAdapter := redisadapter.New(cfg.Redis)
	cacheWithRedis := redis.New(redisAdapter)

	// TODO: save them with the model name -> gemma2:2b: ollamaP
	ollamaP := ollama.New()
	openrouterP := openrouter.New(cfg.OpenRouter)
	providers := map[string]aiservice.Provider{
		"ollama":     ollamaP,
		"openrouter": openrouterP,
	}

	aiSvc := aiservice.New(cfg.AiSvc, ollamaP, providers, cacheWithRedis)
	server := httpserver.New(cfg.Httpserver, aiSvc)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := server.Serve(ctx); err != nil {
		panic(err)
	}

}
