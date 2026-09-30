package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/yonesbeheshti/ai-gateway/delivery/httpserver"
	"github.com/yonesbeheshti/ai-gateway/provider/ollama"
	"github.com/yonesbeheshti/ai-gateway/provider/openrouter"
	"github.com/yonesbeheshti/ai-gateway/service/aiservice"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Note: .env file not found, reading from system env")
	}

	cfg := httpserver.Config{
		Addr: "",
		Port: 8080,
	}

	// TODO: save them with the model name -> gemma2:2b: ollamaP
	ollamaP := ollama.New()
	openrouterP := openrouter.New()
	providers := map[string]aiservice.Provider{
		"ollama":     ollamaP,
		"openrouter": openrouterP,
	}

	aiSvc := aiservice.New(ollamaP, providers)
	server := httpserver.New(cfg, aiSvc)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := server.Serve(ctx); err != nil {
		panic(err)
	}



}
