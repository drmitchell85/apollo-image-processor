# Image Processing Service APIs

## Batch Upload API

### POST /upload
Status: In development

Uploads multiple images for processing in a single batch.

**Request:**
- Content-Type: multipart/form-data
- Body: 
  - files[]: Multiple image files (JPG, PNG only)

**Response:**
```json
{
    "batch_id": "uuid",
    "status": "created",
    "total_images": 5,
    "created_at": "2024-01-26T10:00:00Z"
}
```

**Processing:**
1. Validates all uploaded files are valid images (JPG/PNG)
2. Generates a new batch UUID
3. Creates batch record in PostgreSQL with status "created"
4. For each image:
   - Generates image UUID
   - Stores original file to filesystem
   - Creates image record in PostgreSQL
5. Updates batch record with total image count
6. Queues batch for processing
7. Returns batch ID and initial status

**Error Responses:**
- 400: Invalid file type
- 413: Batch too large
- 500: Server error

## Batch Status API

### GET /batch/{id}
Status: Planned

Retrieves status and progress information for a specific batch.

**Parameters:**
- id: Batch UUID (path parameter)

**Response:**
```json
{
    "batch_id": "uuid",
    "status": "processing",
    "created_at": "2024-01-26T10:00:00Z",
    "completed_at": null,
    "total_images": 5,
    "processed_images": 2,
    "progress": 40
}
```

**Processing:**
1. Looks up batch by UUID in PostgreSQL
2. Calculates current progress
3. Returns batch status and progress information

**Error Responses:**
- 404: Batch not found
- 500: Server error

### GET /batch/{id}/images
Status: Planned

Retrieves details of all images in a specific batch.

**Parameters:**
- id: Batch UUID (path parameter)

**Response:**
```json
{
    "batch_id": "uuid",
    "images": [
        {
            "id": "uuid",
            "original_name": "photo1.jpg",
            "status": "processed",
            "error": null,
            "created_at": "2024-01-26T10:00:00Z",
            "processed_at": "2024-01-26T10:01:00Z"
        },
        // ... additional images
    ]
}
```

**Processing:**
1. Looks up batch by UUID in PostgreSQL
2. Retrieves all associated image records
3. Returns image details

**Error Responses:**
- 404: Batch not found
- 500: Server error