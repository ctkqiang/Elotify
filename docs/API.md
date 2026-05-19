# Push Notification Service API

## Overview

This is a comprehensive Push Notification Service built with ElysiaJS, TypeScript, and Prisma ORM, following the MVVM (Model-View-ViewModel) architectural pattern.

Every API operation is logged with full context for debugging:
- Request/response timestamps
- Operation duration metrics
- User and request IDs for tracing
- Complete error information with stack traces
- Database query details
- Authentication events

See the README.md Logging & Debugging section for details on viewing and controlling logs.

## Architecture

### MVVM Pattern

- **Model** (`src/models/`): Database layer - handles direct database operations
  - `CampaignModel`: Campaign lifecycle management
  - `SubscriptionModel`: User subscription management
  - `NotificationLogModel`: Notification delivery tracking

- **ViewModel** (`src/viewmodels/`): Business logic layer
  - `PushNotificationViewModel`: Orchestrates models, handles business rules, enforces permissions

- **View** (`src/views/`): API response layer
  - DTOs (Data Transfer Objects)
  - Response mappers
  - Input validation contracts

### Security Layers

1. **Authentication** (`src/middleware/auth.middleware.ts`)
   - Bearer token validation
   - Admin role verification
   - Authorization checks

2. **Validation** (`src/middleware/validation.middleware.ts`)
   - Input sanitization
   - Field length validation
   - Type checking
   - Date validation

3. **Error Handling** (`src/middleware/error.middleware.ts`)
   - Centralized error handling
   - HTTP status code mapping
   - Structured error responses
   - Request logging

## Authentication

All endpoints (except `/health` and `/`) require Bearer token authentication.

```
Authorization: Bearer <token>
```

### Available Tokens (for development)
- `admin-token-super` - SUPER_ADMIN role
- `admin-token-regular` - ADMIN role

**For production**: Replace with JWT or OAuth2 implementation.

## API Endpoints

### POST - Push Notifications

#### 1. Push to All Users
```
POST /push/all
Authorization: Bearer <token>
Content-Type: application/json

{
  "title_identifier": "maintenance_alert",
  "locale_contents": {
    "en": {
      "title": "System Maintenance",
      "body": "System will be down for 2 hours",
      "icon_url": "https://example.com/icon.png",
      "action_url": "https://example.com/status"
    },
    "zh": {
      "title": "系统维护",
      "body": "系统将停机2小时"
    }
  },
  "scheduled_at": "2024-05-20T10:00:00Z"
}
```

**Response (201)**:
```json
{
  "campaign_id": "uuid-here",
  "total_subscribers": 15000,
  "timestamp": "2024-05-18T10:00:00Z"
}
```

#### 2. Push to Specific User
```
POST /push/user/{user_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title_identifier": "order_update",
  "locale_contents": {
    "en": {
      "title": "Order Shipped",
      "body": "Your order has been shipped"
    }
  }
}
```

#### 3. Push to Segment/Group
```
POST /push/group/{segment_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title_identifier": "vip_offer",
  "locale_contents": {
    "en": {
      "title": "Exclusive VIP Offer",
      "body": "Get 50% off premium features"
    }
  }
}
```

#### 4. Push to Topic
```
POST /push/topic/{topic_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title_identifier": "product_update",
  "locale_contents": {
    "en": {
      "title": "New Product Available",
      "body": "Check out our latest product"
    }
  }
}
```

### GET - Retrieve Notifications

#### 1. Health Check
```
GET /health
```

**Response (200)**:
```json
{
  "status": "ok",
  "timestamp": "2024-05-18T10:00:00Z",
  "version": "1.0.0"
}
```

#### 2. Latest Notifications
```
GET /notification/latest
```

**Response (200)**:
```json
[
  {
    "id": "uuid",
    "title_identifier": "maintenance_alert",
    "status": "QUEUED",
    "created_by": 1,
    "created_at": "2024-05-18T10:00:00Z",
    "scheduled_at": "2024-05-20T10:00:00Z"
  }
]
```

#### 3. User Notifications
```
GET /notification/user/{user_id}
Authorization: Bearer <token>
```

**Response (200)**:
```json
[
  {
    "id": "uuid",
    "campaign_id": "uuid",
    "subscription_id": "uuid",
    "user_id": "user-123",
    "sent_by": 1,
    "status": "SENT",
    "http_status_code": 201,
    "created_at": "2024-05-18T10:00:00Z"
  }
]
```

#### 4. Segment Notifications
```
GET /notification/group/{segment_id}
Authorization: Bearer <token>
```

#### 5. Topic Notifications
```
GET /notification/topic/{topic_id}
Authorization: Bearer <token>
```

### PUT - Update Notifications

```
PUT /notification/{notification_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "SENDING"
}
```

**Valid statuses**: `DRAFT`, `QUEUED`, `SENDING`, `COMPLETED`, `FAILED`

**Response (200)**:
```json
{
  "id": "uuid",
  "title_identifier": "maintenance_alert",
  "status": "SENDING",
  "created_by": 1,
  "created_at": "2024-05-18T10:00:00Z"
}
```

### DELETE - Delete Notifications

```
DELETE /notification/{notification_id}
Authorization: Bearer <token>
```

**Response (204)**: No content

## Error Responses

### Unauthorized (401)
```json
{
  "error": "UnauthorizedException",
  "message": "Missing or invalid Authorization header",
  "status_code": 401,
  "timestamp": "2024-05-18T10:00:00Z"
}
```

### Forbidden (403)
```json
{
  "error": "ForbiddenException",
  "message": "Insufficient permissions",
  "status_code": 403,
  "timestamp": "2024-05-18T10:00:00Z"
}
```

### Validation Error (400)
```json
{
  "error": "ValidationException",
  "message": "title_identifier is required and must be a string",
  "status_code": 400,
  "timestamp": "2024-05-18T10:00:00Z"
}
```

### Not Found (404)
```json
{
  "error": "NotFoundError",
  "message": "Notification not found",
  "status_code": 404,
  "timestamp": "2024-05-18T10:00:00Z"
}
```

## Request Validation Rules

### title_identifier
- Required: Yes
- Type: string
- Max length: 100 characters

### locale_contents
- Required: Yes
- Type: object
- Each locale must have:
  - **title** (required, max 120 chars)
  - **body** (required, max 5000 chars)
  - **icon_url** (optional, max 512 chars)
  - **action_url** (optional, max 512 chars)

### scheduled_at
- Required: No
- Type: ISO 8601 date string
- Must be in the future

## Role-Based Access Control

| Endpoint | SUPER_ADMIN | ADMIN | Public |
|----------|---|---|---|
| POST /push/all | YES | NO | NO |
| POST /push/user/{user_id} | YES | YES | NO |
| POST /push/group/{segment_id} | YES | YES | NO |
| POST /push/topic/{topic_id} | YES | YES | NO |
| GET /health | YES | YES | NO |
| GET /notification/* | YES | YES | NO |
| PUT /notification/{id} | YES | NO | NO |
| DELETE /notification/{id} | YES | NO | NO |

## Example cURL Requests

### Create a campaign for all users
```bash
curl -X POST http://localhost:3000/push/all \
  -H "Authorization: Bearer admin-token-super" \
  -H "Content-Type: application/json" \
  -d '{
    "title_identifier": "test_campaign",
    "locale_contents": {
      "en": {
        "title": "Test Notification",
        "body": "This is a test"
      }
    }
  }'
```

### Get latest notifications
```bash
curl -X GET http://localhost:3000/notification/latest
```

### Get user notifications
```bash
curl -X GET http://localhost:3000/notification/user/user-123 \
  -H "Authorization: Bearer admin-token-regular"
```

### Update notification status
```bash
curl -X PUT http://localhost:3000/notification/uuid-here \
  -H "Authorization: Bearer admin-token-super" \
  -H "Content-Type: application/json" \
  -d '{"status": "SENDING"}'
```
