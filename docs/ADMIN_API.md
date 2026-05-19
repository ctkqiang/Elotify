# Admin API Documentation

## Overview

The Admin API provides campaign management endpoints for users with ADMIN or SUPER_ADMIN roles. These endpoints follow REST standards and include comprehensive validation, pagination, and filtering.

## Authentication

All admin endpoints require Bearer token authentication with ADMIN or SUPER_ADMIN role:

```
Authorization: Bearer admin-token-regular  (ADMIN role)
Authorization: Bearer admin-token-super    (SUPER_ADMIN role)
```

## Industry Standards Applied

- RESTful architecture (HTTP methods: GET, POST, PUT, DELETE)
- Resource-based URL structure
- Consistent response format (JSON)
- Proper HTTP status codes (200, 201, 400, 401, 403, 404, 500)
- Pagination for list endpoints
- Filtering and sorting capabilities
- Comprehensive error messages
- Request/Response validation
- ACID transaction support

---

## Test Data Management

### Create Web Push Subscription

**Endpoint:** `POST /admin/subscriptions`

**Role Required:** ADMIN, SUPER_ADMIN

**Purpose:** Create a test web push subscription for a user (useful for testing notification delivery)

#### Request

```bash
curl -X POST http://localhost:3000/admin/subscriptions \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "browser": "CHROME",
    "endpoint": "https://push.example.com/endpoint/user-123",
    "p256dh_key": "base64-encoded-public-key",
    "auth_key": "base64-encoded-auth-secret"
  }'
```

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| user_id | string | Yes | User identifier |
| browser | string | Yes | Browser type: CHROME, FIREFOX, SAFARI, EDGE |
| endpoint | string | Yes | Web push endpoint URL (from browser) |
| p256dh_key | string | Yes | ECDH public key (base64) |
| auth_key | string | Yes | Authentication secret (base64) |

#### Response - 201 Created

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440010",
  "user_id": "user-123",
  "browser": "CHROME",
  "endpoint": "https://fcm.googleapis.com/fcm/send/example-endpoint",
  "is_active": true,
  "created_at": "2026-05-19T10:00:00.000Z",
  "timestamp": "2026-05-19T10:00:00.000Z"
}
```

#### Example Workflow: Create subscription then push notification

```bash
# 1. Create a test subscription for user-123
curl -X POST http://localhost:3000/admin/subscriptions \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "browser": "CHROME",
    "endpoint": "https://push.example.com/endpoint/abc123",
    "p256dh_key": "test-key",
    "auth_key": "test-auth"
  }'

# 2. Now push a notification to that user
curl -X POST http://localhost:3000/push/user/user-123 \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "title_identifier": "welcome_msg",
    "locale_contents": {
      "en": {
        "title": "Welcome!",
        "body": "Welcome to our service"
      }
    }
  }'

# 3. Retrieve notification logs for user-123
curl -X GET http://localhost:3000/notification/user/user-123 \
  -H "Authorization: Bearer admin-token-regular"
```

---

## 1. Create Campaign

**Endpoint:** `POST /admin/campaigns`

**Role Required:** ADMIN, SUPER_ADMIN

**Purpose:** Create a new notification campaign

### Request

```bash
curl -X POST http://localhost:3000/admin/campaigns \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "title_identifier": "summer_promotion",
    "description": "Summer sale campaign for all users",
    "locale_contents": {
      "en": {
        "title": "Summer Sale",
        "body": "Get 30% off on all products",
        "icon_url": "https://example.com/summer.png",
        "action_url": "https://example.com/sale"
      },
      "zh": {
        "title": "夏季促销",
        "body": "所有产品享受30% 折扣",
        "icon_url": "https://example.com/summer.png",
        "action_url": "https://example.com/sale"
      }
    },
    "scheduled_at": null
  }'
```

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| title_identifier | string | Yes | 1-100 chars | Unique campaign identifier |
| description | string | No | Max 500 chars | Internal campaign description |
| locale_contents | object | Yes | Min 1 locale | Multi-language content map |
| scheduled_at | ISO 8601 | No | Future date | Scheduled delivery time |

### Locale Content Schema

```json
{
  "locale_id": {
    "title": "string (1-256 chars)",
    "body": "string (1-1000 chars)",
    "icon_url": "string (URL, optional)",
    "action_url": "string (URL, optional)"
  }
}
```

### Response - 201 Created

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title_identifier": "summer_promotion",
  "description": "Summer sale campaign for all users",
  "status": "DRAFT",
  "created_by": 2,
  "created_at": "2026-05-19T10:00:00.000Z",
  "scheduled_at": null,
  "locales": ["en", "zh"],
  "total_content_items": 2
}
```

### Error Responses

**400 Bad Request** - Invalid input
```json
{
  "error": "Bad Request",
  "message": "Validation failed: title_identifier is required",
  "status_code": 400,
  "details": {
    "field": "title_identifier",
    "reason": "required"
  }
}
```

**401 Unauthorized** - Missing/invalid token
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing Authorization header",
  "status_code": 401
}
```

**403 Forbidden** - Insufficient permissions
```json
{
  "error": "Forbidden",
  "message": "User does not have ADMIN role",
  "status_code": 403
}
```

---

## 2. List Campaigns

**Endpoint:** `GET /admin/campaigns`

**Role Required:** ADMIN, SUPER_ADMIN

**Purpose:** Retrieve paginated list of campaigns with filtering and sorting

### Request

```bash
curl -X GET "http://localhost:3000/admin/campaigns?page=1&limit=10&status=DRAFT&sort_by=created_at" \
  -H "Authorization: Bearer admin-token-regular"
```

### Query Parameters

| Parameter | Type | Default | Constraints | Description |
|-----------|------|---------|-------------|-------------|
| page | integer | 1 | Min 1 | Page number (1-indexed) |
| limit | integer | 10 | 1-100 | Records per page |
| status | string | None | DRAFT, QUEUED, SENDING, COMPLETED, FAILED | Filter by status |
| sort_by | string | created_at | created_at, title_identifier, status | Sort field |
| sort_order | string | DESC | ASC, DESC | Sort direction |
| search | string | None | Any string | Search title_identifier or description |

### Response - 200 OK

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title_identifier": "summer_promotion",
      "status": "DRAFT",
      "created_by": 2,
      "created_at": "2026-05-19T10:00:00.000Z",
      "scheduled_at": null,
      "locales": ["en", "zh"],
      "total_content_items": 2
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "title_identifier": "flash_sale",
      "status": "QUEUED",
      "created_by": 2,
      "created_at": "2026-05-19T09:00:00.000Z",
      "scheduled_at": "2026-05-20T14:00:00.000Z",
      "locales": ["en"],
      "total_content_items": 1
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_records": 47,
    "per_page": 10,
    "has_next": true,
    "has_previous": false
  },
  "meta": {
    "timestamp": "2026-05-19T10:00:00.000Z",
    "request_id": "req-abc123"
  }
}
```

### Pagination Example

```
Page 1 (default): ?page=1&limit=10
Page 2: ?page=2&limit=10
Next page: Check has_next in response
```

### Filtering Examples

```
Draft campaigns only: ?status=DRAFT
Scheduled campaigns: ?status=QUEUED
Completed campaigns: ?status=COMPLETED
Search by title: ?search=promotion
Sort by status: ?sort_by=status&sort_order=ASC
```

---

## 3. Get Campaign Details

**Endpoint:** `GET /admin/campaigns/:id`

**Role Required:** ADMIN, SUPER_ADMIN

**Purpose:** Retrieve detailed information about a specific campaign

### Request

```bash
curl -X GET http://localhost:3000/admin/campaigns/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer admin-token-regular"
```

### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Campaign ID |

### Response - 200 OK

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title_identifier": "summer_promotion",
  "description": "Summer sale campaign for all users",
  "status": "DRAFT",
  "created_by": 2,
  "created_by_username": "admin",
  "created_at": "2026-05-19T10:00:00.000Z",
  "updated_at": "2026-05-19T10:00:00.000Z",
  "scheduled_at": null,
  "contents": [
    {
      "locale_id": "en",
      "title": "Summer Sale",
      "body": "Get 30% off on all products",
      "icon_url": "https://example.com/summer.png",
      "action_url": "https://example.com/sale"
    },
    {
      "locale_id": "zh",
      "title": "夏季促销",
      "body": "所有产品享受30% 折扣",
      "icon_url": "https://example.com/summer.png",
      "action_url": "https://example.com/sale"
    }
  ],
  "metrics": {
    "total_content_items": 2,
    "supported_locales": ["en", "zh"],
    "total_subscribers_target": 0,
    "delivery_status": "not_sent"
  }
}
```

### Error Response - 404 Not Found

```json
{
  "error": "Not Found",
  "message": "Campaign with ID 550e8400-e29b-41d4-a716-446655440000 not found",
  "status_code": 404
}
```

---

## 4. Update Campaign

**Endpoint:** `PUT /admin/campaigns/:id`

**Role Required:** ADMIN, SUPER_ADMIN

**Purpose:** Update a campaign (only DRAFT status campaigns can be modified)

### Request

```bash
curl -X PUT http://localhost:3000/admin/campaigns/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated campaign description",
    "locale_contents": {
      "en": {
        "title": "Updated Title",
        "body": "Updated message",
        "icon_url": "https://example.com/new-icon.png",
        "action_url": "https://example.com/updated"
      }
    }
  }'
```

### Request Body (All Optional)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| description | string | Max 500 chars | Campaign description |
| locale_contents | object | Min 1 locale | Updated content per locale |
| scheduled_at | ISO 8601 | Future date | Update scheduled time |

### Response - 200 OK

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title_identifier": "summer_promotion",
  "description": "Updated campaign description",
  "status": "DRAFT",
  "created_by": 2,
  "created_at": "2026-05-19T10:00:00.000Z",
  "updated_at": "2026-05-19T10:05:00.000Z",
  "scheduled_at": null,
  "contents": [
    {
      "locale_id": "en",
      "title": "Updated Title",
      "body": "Updated message",
      "icon_url": "https://example.com/new-icon.png",
      "action_url": "https://example.com/updated"
    }
  ]
}
```

### Error Response - 409 Conflict

```json
{
  "error": "Conflict",
  "message": "Cannot update campaign - status is SENDING. Only DRAFT campaigns can be modified.",
  "status_code": 409
}
```

---

## Campaign Status Machine

```
CREATE
  ↓
DRAFT (Editable)
  ├─ Can: Update, Delete, Send
  ├─ Cannot: Edit content once sent
  └─ Transition to: QUEUED or SENDING
    ↓
QUEUED (Scheduled)
  ├─ Waiting for scheduled_at time
  ├─ Cannot: Modify
  └─ Transition to: SENDING (automatic)
    ↓
SENDING (In Progress)
  ├─ Currently delivering
  ├─ Cannot: Modify
  └─ Transition to: COMPLETED or FAILED
    ↓
COMPLETED or FAILED (Terminal)
  ├─ Cannot: Modify or Send
  └─ View: Full delivery metrics
```

---

## Error Code Reference

| Code | Meaning | Action |
|------|---------|--------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Fix request body/parameters |
| 401 | Unauthorized | Check auth token |
| 403 | Forbidden | User lacks required role |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Operation conflicts with current state |
| 500 | Server Error | Contact support |

---

## Rate Limiting

No rate limiting in development. Production should implement:
- 100 requests/minute per user
- 1000 requests/minute per IP

---

## Best Practices

### 1. Always Use Pagination
```bash
# Bad: Get all campaigns
GET /admin/campaigns

# Good: Use pagination
GET /admin/campaigns?page=1&limit=20
```

### 2. Filter by Status
```bash
# Find only draft campaigns
GET /admin/campaigns?status=DRAFT&limit=50
```

### 3. Handle Errors Gracefully
Always check HTTP status code and error details before retrying.

### 4. Validate Before Submit
Validate locale_contents and title_identifier before API call to reduce roundtrips.

### 5. Check Campaign Status Before Modifying
Only DRAFT campaigns can be updated. Check status first:
```bash
GET /admin/campaigns/:id  # Check status field
PUT /admin/campaigns/:id  # Update only if status=DRAFT
```

---

## Example Workflows

### Workflow 1: Create and Schedule Campaign

```bash
# 1. Create draft campaign
POST /admin/campaigns
{
  "title_identifier": "flash_sale",
  "locale_contents": { "en": {...} }
}
Response: status = DRAFT, id = abc123

# 2. Review campaign details
GET /admin/campaigns/abc123

# 3. Update with scheduled time
PUT /admin/campaigns/abc123
{
  "scheduled_at": "2026-05-20T14:00:00Z"
}
Response: status = QUEUED (auto-transitioned)

# 4. Monitor delivery
GET /admin/campaigns/abc123  # Check status updates
```

### Workflow 2: Bulk List and Filter

```bash
# List all draft campaigns
GET /admin/campaigns?status=DRAFT&limit=50

# Search for specific campaign
GET /admin/campaigns?search=promotion&limit=20

# Sort by creation date (newest first)
GET /admin/campaigns?sort_by=created_at&sort_order=DESC&limit=10
```

---

## Security Considerations

- All requests require ADMIN+ role
- Campaign IDs are UUIDs (not sequential)
- Timestamps are immutable after creation
- Admin actions are logged for audit trail
- Sensitive data is never returned in error messages
