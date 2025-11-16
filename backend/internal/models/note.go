package models

import (
	"time"

	"github.com/google/uuid"
)

// Note represents a note document with embeddings
type Note struct {
	ID           uuid.UUID              `json:"id"`
	UserID       *uuid.UUID             `json:"user_id,omitempty"`
	Title        string                 `json:"title"`
	Content      string                 `json:"content"`
	ChunkContent string                 `json:"chunk_content"`
	Embedding    []float32              `json:"-"` // Hidden from JSON, stored as vector
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Note            *Note   `json:"note"`
	SimilarityScore float32 `json:"similarity_score"`
}
