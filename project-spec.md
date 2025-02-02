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
CREATE TABLE public.batches (
	batch_id uuid NOT NULL,
	status public."batch_status" NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	completed_at timestamp NULL,
	total_images int4 DEFAULT 0 NOT NULL,
	processed_images int4 DEFAULT 0 NOT NULL,
	CONSTRAINT batches_pk PRIMARY KEY (batch_id)
);
CREATE INDEX idx_batches_created_at ON public.batches USING btree (created_at);
CREATE INDEX idx_batches_status ON public.batches USING btree (status);

-- Images
CREATE TABLE public.images (
	image_id uuid NOT NULL,
	batch_id uuid NOT NULL,
	status public.image_status DEFAULT 'pending'::image_status NOT NULL,
	error text NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	processed_at timestamp NULL,
	image_name text,
	image bytea,
	CONSTRAINT error_only_on_failure CHECK ((((status = 'failed'::image_status) AND (error IS NOT NULL)) OR ((status <> 'failed'::image_status) AND (error IS NULL)))),
	CONSTRAINT images_pkey PRIMARY KEY (image_id),
	CONSTRAINT processed_at_required CHECK ((((status = 'completed'::image_status) AND (processed_at IS NOT NULL)) OR (status <> 'completed'::image_status)))
);
CREATE INDEX idx_images_created_at ON public.images USING btree (created_at);
CREATE INDEX idx_images_status ON public.images USING btree (status);


-- public.images foreign keys

ALTER TABLE public.images ADD CONSTRAINT images_batch_id_fkey FOREIGN KEY (batch_id) REFERENCES public.batches(batch_id);
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
