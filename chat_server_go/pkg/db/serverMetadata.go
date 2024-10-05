package db

import (
	"context"
	"log"
	"os"
	"time"
)

// Metadata stores application metadata
type Metadata struct {
	AppVersion               string `json:"app_version"`
	TokenAudience            string `json:"token_audience"`
	HistoryLength            int    `json:"history_length"`
	MaxUserMessageLen        int    `json:"max_user_message_len"`
	CorsOrigin               string `json:"cors_origin"`
	RetrieverLength          int    `json:"retriever_length"`
	GoogleChatModelName      string `json:"google_chat_model_name"`
	GoogleEmbeddingModelName string `json:"google_embedding_model_name"`
	FrontEndDomain           string `json:"front_end_domain"`
}

func (d *MovieDB) GetMetadata(ctx context.Context, appVersion string) (*Metadata, error) {
	if os.Getenv("LOCAL") == "true" {
		return getMetadataLocal()
	}
	return d.getServerMetadata(ctx, appVersion)
}

// getMetadata retrieves metadata from the database
func (d *MovieDB) getServerMetadata(ctx context.Context, appVersion string) (*Metadata, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	query := `SELECT * FROM app_metadata WHERE "app_version" = $1;`
	metadata := &Metadata{}
	rows := d.DB.QueryRowContext(dbCtx, query, appVersion)
	err := rows.Scan(
		&metadata.AppVersion,
		&metadata.TokenAudience,
		&metadata.HistoryLength,
		&metadata.MaxUserMessageLen,
		&metadata.CorsOrigin,
		&metadata.RetrieverLength,
		&metadata.GoogleChatModelName,
		&metadata.GoogleEmbeddingModelName,
		&metadata.FrontEndDomain,
	)
	if err != nil {
		return metadata, err
	}
	log.Println(metadata)
	return metadata, nil
}

func getMetadataLocal() (*Metadata, error) {
	metadata := &Metadata{
		AppVersion:               "local",
		HistoryLength:            10,
		MaxUserMessageLen:        1000,
		CorsOrigin:               "http://localhost:5173",
		RetrieverLength:          10,
		GoogleChatModelName:      "gemini-1.5-flash",
		GoogleEmbeddingModelName: "text-embedding-004",
	}
	log.Println(metadata)
	return metadata, nil
}
