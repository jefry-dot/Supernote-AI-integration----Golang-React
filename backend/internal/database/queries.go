package database

import (
	"context"
	"encoding/json"
	"fmt"

	"supernote-ai/backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// InsertNote inserts a new note with embedding into the database
func (db *DB) InsertNote(ctx context.Context, note *models.Note) error {
	query := `
		INSERT INTO notes (id, user_id, title, content, chunk_content, embedding, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(note.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Generate UUID if not provided
	if note.ID == uuid.Nil {
		note.ID = uuid.New()
	}

	err = db.Pool.QueryRow(ctx, query,
		note.ID,
		note.UserID,
		note.Title,
		note.Content,
		note.ChunkContent,
		note.Embedding, // pgx handles []float32 to vector conversion
		metadataJSON,
		note.CreatedAt,
		note.UpdatedAt,
	).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert note: %w", err)
	}

	return nil
}

// SearchNotesBySimilarity performs vector similarity search
func (db *DB) SearchNotesBySimilarity(ctx context.Context, queryEmbedding []float32, limit int) ([]*models.SearchResult, error) {
	query := `
		SELECT
			id, user_id, title, content, chunk_content, metadata, created_at, updated_at,
			embedding <=> $1 AS similarity_score
		FROM notes
		ORDER BY similarity_score ASC
		LIMIT $2
	`

	rows, err := db.Pool.Query(ctx, query, queryEmbedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}
	defer rows.Close()

	var results []*models.SearchResult
	for rows.Next() {
		note := &models.Note{}
		var metadataJSON []byte
		var similarityScore float32

		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.ChunkContent,
			&metadataJSON,
			&note.CreatedAt,
			&note.UpdatedAt,
			&similarityScore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Unmarshal metadata
		if err := json.Unmarshal(metadataJSON, &note.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		results = append(results, &models.SearchResult{
			Note:            note,
			SimilarityScore: similarityScore,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// GetNoteByID retrieves a note by its ID
func (db *DB) GetNoteByID(ctx context.Context, id uuid.UUID) (*models.Note, error) {
	query := `
		SELECT id, user_id, title, content, chunk_content, metadata, created_at, updated_at
		FROM notes
		WHERE id = $1
	`

	note := &models.Note{}
	var metadataJSON []byte

	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.ChunkContent,
		&metadataJSON,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("note not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	// Unmarshal metadata
	if err := json.Unmarshal(metadataJSON, &note.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return note, nil
}

// GetAllNotes retrieves all notes (with pagination)
func (db *DB) GetAllNotes(ctx context.Context, limit, offset int) ([]*models.Note, error) {
	query := `
		SELECT id, user_id, title, content, chunk_content, metadata, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		var metadataJSON []byte

		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.ChunkContent,
			&metadataJSON,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Unmarshal metadata
		if err := json.Unmarshal(metadataJSON, &note.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return notes, nil
}

// DeleteNote deletes a note by ID
func (db *DB) DeleteNote(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}
