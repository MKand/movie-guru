// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"context"
	"log"
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

// getMetadata retrieves metadata from the database
func (moveiDB *MovieDB) GetServerMetadata(appVersion string) (*Metadata, error) {
	query := `SELECT * FROM app_metadata WHERE "app_version" = $1;`
	metadata := &Metadata{}
	rows := moveiDB.DB.QueryRowContext(context.Background(), query, appVersion)
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
