package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"Supportagent/backend/internal/config"
	"Supportagent/backend/internal/database"
	"Supportagent/backend/internal/handler"
	appmcp "Supportagent/backend/internal/mcp"
	"Supportagent/backend/internal/policy"
	"Supportagent/backend/internal/service"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	// -----------------------------------------
	// 1. Load configuration
	// -----------------------------------------

	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// -----------------------------------------
	// 2. Connect to PostgreSQL
	// -----------------------------------------

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// -----------------------------------------
	// 3. Create MCP Client
	// -----------------------------------------

	ctx := context.Background()

	mcpClient, err := appmcp.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	// -----------------------------------------
	// 4. Create OpenAI client
	// -----------------------------------------

	apiKey := os.Getenv("OPENROUTER_API_KEY")

	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	// -----------------------------------------
	// 5. AI Service
	// -----------------------------------------
	policyEngine := policy.NewEngine()
	aiService := service.NewAIService(
		client,
		mcpClient,
		policyEngine,
	)

	// -----------------------------------------
	// 6. Chat Handler
	// -----------------------------------------

	chatHandler := handler.NewChatHandler(
		aiService,
	)

	// -----------------------------------------
	// 7. HTTP routes
	// -----------------------------------------

	http.HandleFunc("/health", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Support Agent Backend is healthy"))
	})

	http.HandleFunc(
		"/chat",
		chatHandler.Chat,
	)

	// -----------------------------------------
	// 8. Start server
	// -----------------------------------------

	log.Println("Support Agent Backend running on :8080")

	if err := http.ListenAndServe(
		":8080",
		enableCORS(http.DefaultServeMux),
	); err != nil {
		log.Fatal(err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set(
			"Access-Control-Allow-Origin",
			"http://localhost:3000",
		)

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, OPTIONS",
		)

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
