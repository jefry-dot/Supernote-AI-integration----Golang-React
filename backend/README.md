# Supernote AI Backend

Backend API untuk Supernote AI Integration dengan Golang, PostgreSQL (Supabase), dan Google Gemini.

## Features

- ✅ RESTful API dengan Gin framework
- ✅ PostgreSQL database dengan pgvector extension
- ✅ Vector embeddings untuk semantic search
- ✅ Connection pooling dan health checks
- 🚧 Gemini API integration (coming soon)
- 🚧 Semantic search (coming soon)
- 🚧 RAG-based chatbot (coming soon)

## Setup

### 1. Prerequisites

- Go 1.21 atau lebih baru
- PostgreSQL dengan pgvector extension (Supabase)
- Gemini API key

### 2. Environment Variables

Copy `.env.example` ke `.env` dan isi dengan credentials Anda:

```bash
cp .env.example .env
```

Edit `.env` dan update:
- `DATABASE_URL` - Connection string Supabase
- `GEMINI_API_KEY` - API key dari Google AI Studio

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Server

```bash
go run cmd/server/main.go
```

Server akan running di `http://localhost:8080`

## API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "database": "connected",
  "service": "supernote-ai-backend",
  "version": "1.0.0"
}
```

### Get All Notes
```
GET /api/notes?limit=20&offset=0
```

### Get Note by ID
```
GET /api/notes/:id
```

### Delete Note
```
DELETE /api/notes/:id
```

### Search Notes (TODO)
```
POST /api/search
```

### Chat with AI (TODO)
```
POST /api/chat
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── config/               # Configuration
│   ├── database/             # Database connection & queries
│   ├── handlers/             # HTTP handlers
│   ├── models/               # Data models
│   ├── services/             # Business logic (AI, embeddings)
│   └── utils/                # Helper functions
├── .env                      # Environment variables
├── go.mod                    # Go dependencies
└── README.md
```

## Development

### Run in development mode
```bash
ENV=development go run cmd/server/main.go
```

### Build for production
```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## Next Steps

- [ ] Implement Gemini embeddings service
- [ ] Implement semantic search
- [ ] Implement RAG-based chatbot
- [ ] Add file upload for notes
- [ ] Add authentication
- [ ] Add rate limiting
