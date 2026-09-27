package main

import (
	"context"

	"github.com/yonesbeheshti/ai-gateway/delivery/httpserver"
)

func main() {

	server := httpserver.New("", 8080)
	ctx := context.Background()
	if err := server.Serve(ctx); err != nil {
		panic(err)
	}

}
