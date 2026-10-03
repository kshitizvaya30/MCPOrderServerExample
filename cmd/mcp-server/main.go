package main

import (
	"context"
	"log"

	"Supportagent/backend/internal/config"
	"Supportagent/backend/internal/database"
	appmcp "Supportagent/backend/internal/mcp"
	"Supportagent/backend/internal/repository"
	"Supportagent/backend/internal/service"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository)

	server := appmcp.NewServer(orderService)

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		log.Fatal(err)
	}
}
