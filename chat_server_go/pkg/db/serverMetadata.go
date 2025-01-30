package db

import (
	"context"
	"log"
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
	ServerDomain             string `json:"server_domain"`
}

func (d *MovieDB) GetMetadata(ctx context.Context, appVersion string) (*Metadata, error) {
	return d.getServerMetadata(ctx, appVersion)
}

// getMetadata retrieves metadata from the database
func (d *MovieDB) getServerMetadata(ctx context.Context, appVersion string) (*Metadata, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	query := `SELECT * FROM app_metadata WHERE "appversion" = $1;`
	metadata := &Metadata{}
	rows := d.DB.QueryRowContext(dbCtx, query, appVersion)
	err := rows.Scan(
		&metadata.AppVersion,
		&metadata.TokenAudience,
		&metadata.HistoryLength,
		&metadata.MaxUserMessageLen,
		&metadata.CorsOrigin,
		&metadata.RetrieverLength,
		&metadata.ServerDomain,
		&metadata.GoogleChatModelName,
		&metadata.GoogleEmbeddingModelName,
	)
	if err != nil {
		return metadata, err
	}
	log.Println(metadata)
	return metadata, nil
}
