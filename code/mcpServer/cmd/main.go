package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"log/slog"
	"github.com/firebase/genkit/go/plugins/mcp"
	types "github.com/mkand/movieSearch/pkg/types"
	"github.com/mkand/movieSearch/pkg/db"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/ai"
)

func main(){
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	g := genkit.Init(context.Background(), genkit.WithPlugins(&googlegenai.VertexAI{}))

	movieSearchDB, err := db.GetDB(ctx, g)
	if err != nil {
		slog.Error("Failed to initialize database:", err)
	}
	
	// Tool 1: Search for movies

	genkit.DefineTool(g, "search_movies", "Find relevant movies based on a user's statement",
		func(ctx *ai.ToolContext, userQuery types.UserQuery) ([]types.MovieContextSchema, error) {
			results, err := movieSearchDB.PerformSearch(userQuery)
			if err != nil{
				slog.Error("failed to get movies", "error", err)
				return nil, err 
			}
			return results, nil
		})

	// Start MCP server
	server := mcp.NewMCPServer(g, mcp.MCPServerOptions{
		Name: "movie-search",
	})
	slog.Info("Starting MCP server", "name", "movie-search", "tools", server.ListRegisteredTools())

	if err := server.ServeStdio(); err != nil && err != context.Canceled {
		slog.Error("MCP server error", "error", err)
		os.Exit(1)
	}
}
