package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"supernote-ai/backend/internal/config"
	"supernote-ai/backend/internal/database"
	"supernote-ai/backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware (untuk development)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	notesHandler := handlers.NewNotesHandler(db)

	// Setup routes
	setupRoutes(router, healthHandler, notesHandler)

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		db.Close()
		os.Exit(0)
	}()

	// Start server
	port := ":" + cfg.Port
	log.Printf("🚀 Server starting on http://localhost%s", port)
	log.Printf("📝 Environment: %s", cfg.Environment)
	log.Printf("✅ Database: Connected")
	log.Printf("\n📚 Available endpoints:")
	log.Printf("  GET    /health          - Health check")
	log.Printf("  GET    /api/notes       - Get all notes")
	log.Printf("  GET    /api/notes/:id   - Get note by ID")
	log.Printf("  DELETE /api/notes/:id   - Delete note")
	log.Printf("  POST   /api/search      - Search notes (TODO)")
	log.Printf("  POST   /api/chat        - Chat with AI (TODO)\n")

	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, notesHandler *handlers.NotesHandler) {
	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// Notes endpoints
		notes := api.Group("/notes")
		{
			notes.GET("", notesHandler.GetAllNotes)
			notes.GET("/:id", notesHandler.GetNoteByID)
			notes.DELETE("/:id", notesHandler.DeleteNote)
		}

		// Search endpoint (placeholder)
		api.POST("/search", notesHandler.SearchNotes)

		// Chat endpoint (placeholder - akan kita implement nanti)
		api.POST("/chat", func(c *gin.Context) {
			c.JSON(501, gin.H{
				"error": "Chat endpoint not implemented yet",
			})
		})
	}
}
