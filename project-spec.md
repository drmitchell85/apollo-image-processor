# Batch Image Processing System

## Architecture

### Components
- REST API (Go/Gin)
- CLI Tool (Go/Cobra)
- Message Queue (RabbitMQ)
- Worker Service (Go)
- PostgreSQL (primary datastore)
- Redis (queue management/caching)
- Local filesystem storage

### Data Flow
1. Images uploaded via API or CLI
2. Metadata stored in PostgreSQL
3. Job queued in Redis/RabbitMQ
4. Worker processes batch
5. Results saved to filesystem
6. Status updated in PostgreSQL

## Image Processing Specifications
- Input: JPG, PNG
- Output: WebP
- Transformations:
  - Standard: 800x800px (preserved ratio)
  - Thumbnail: 200x200px
  - Grayscale version
  - Metadata stripped
  - Quality: 85%

## Database Schema

### PostgreSQL Tables
```sql
-- Batches
CREATE TABLE batches (
    id UUID PRIMARY KEY,
    status VARCHAR(20),
    created_at TIMESTAMP,
    completed_at TIMESTAMP,
    total_images INT,
    processed_images INT
);

-- Images
CREATE TABLE images (
    id UUID PRIMARY KEY,
    batch_id UUID REFERENCES batches(id),
    original_name VARCHAR(255),
    original_path VARCHAR(255),
    processed_path VARCHAR(255),
    status VARCHAR(20),
    error TEXT,
    created_at TIMESTAMP,
    processed_at TIMESTAMP
);
```

### Redis Keys
- batch:{id}:status
- batch:{id}:progress
- processing:queue

## API Endpoints

### POST /upload
- Accepts multiple images
- Returns batch ID

### GET /batch/{id}
- Returns batch status and progress

### GET /batch/{id}/images
- Returns all images in batch

## CLI Commands
```bash
imgtool watch /path/to/dir  # Watch directory
imgtool status <batch-id>   # Check batch status
imgtool list                # List recent batches
```

## Development Phases

### Phase 1: Core Infrastructure
- API setup
- Database schema
- Basic file handling

### Phase 2: Processing Pipeline
- Queue integration
- Worker service
- Image processing

### Phase 3: CLI Tool
- Directory watching
- Status checking
- Batch management

### Phase 4: Monitoring
- Progress tracking
- Error handling
- Logging
