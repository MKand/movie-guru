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

package embedder

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/genkit"
)

// Embedder wraps the Genkit embedding functionality
type Embedder struct {
	embedder ai.Embedder
	g *genkit.Genkit
}


// NewEmbedder creates a new embedder instance using Vertex AI text-embedding-005
func NewEmbedder(ctx context.Context, g *genkit.Genkit) (*Embedder, error) {
	embeddingModel := googlegenai.VertexAIEmbedder(g, "text-embedding-005")	
	return &Embedder{
		embedder: embeddingModel,
		g: g,
	}, nil
}

// GenerateEmbedding generates a vector embedding for the given text
func (e *Embedder) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}
	doc:= &ai.Document{
                Content:  []*ai.Part{ai.NewTextPart(text)},
            }
	// Generate embedding using Genkit
	 eres, err := genkit.Embed(ctx, e.g,
	 	ai.WithEmbedder(e.embedder),
	 	ai.WithDocs(doc),
	)
		if err != nil {
			return nil, err
		}
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(eres.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Extract the first embedding (we only sent one document)
	embedding := eres.Embeddings[0].Embedding
	
	// Convert to []float32 for sqlite-vec compatibility
	result := make([]float32, len(embedding))
	for i, v := range embedding {
		result[i] = float32(v)
	}

	return result, nil
}
