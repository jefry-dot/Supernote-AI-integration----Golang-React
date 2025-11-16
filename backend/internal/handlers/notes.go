package handlers

import (
	"net/http"
	"strconv"

	"supernote-ai/backend/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotesHandler handles note-related endpoints
type NotesHandler struct {
	db *database.DB
}

// NewNotesHandler creates a new notes handler
func NewNotesHandler(db *database.DB) *NotesHandler {
	return &NotesHandler{db: db}
}

// CreateNoteRequest represents the request body for creating a note
type CreateNoteRequest struct {
	Title    string                 `json:"title" binding:"required"`
	Content  string                 `json:"content" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// GetAllNotes retrieves all notes with pagination
func (h *NotesHandler) GetAllNotes(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get notes from database
	notes, err := h.db.GetAllNotes(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve notes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes":  notes,
		"limit":  limit,
		"offset": offset,
		"count":  len(notes),
	})
}

// GetNoteByID retrieves a single note by ID
func (h *NotesHandler) GetNoteByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid note ID",
		})
		return
	}

	note, err := h.db.GetNoteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

// DeleteNote deletes a note by ID
func (h *NotesHandler) DeleteNote(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid note ID",
		})
		return
	}

	if err := h.db.DeleteNote(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Note deleted successfully",
	})
}

// SearchNotes performs semantic search on notes
// We'll implement this later when we add Gemini embeddings
func (h *NotesHandler) SearchNotes(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Search not implemented yet - needs Gemini embeddings",
	})
}
