# Sequence Diagram: Push Notification Request Flow

![Sequence Diagram](../out/docs/sequence-diagram/Push%20Notification%20Service%20-%20Sequence%20Diagram.png)

## Overview

This diagram illustrates the complete execution flow of a single push notification request from client submission through database persistence. It demonstrates how each layer of the MVVM architecture interacts and the order of operations.

## Detailed Flow Breakdown

### 1. Client Request (Entry Point)

```
Client → API: POST /push/all
Headers: Content-Type: application/json
Body: {
  "title_identifier": "maintenance_alert",
  "locale_contents": {
    "en": {
      "title": "System Maintenance",
      "body": "System will be down for 2 hours",
      "icon_url": "...",
      "action_url": "..."
    },
    "zh": { ... }
  },
  "scheduled_at": "2026-05-20T10:00:00Z"
}
```

### 2. Authentication Middleware

**Purpose**: Validate request authentication and extract admin context

```
AuthMiddleware.extractAuthContext(context)
├─ Get Authorization header
├─ Parse Bearer token
├─ Validate token against VALID_TOKENS map
├─ Return AuthContext {
│   admin_id: number,
│   admin_username: string,
│   admin_role: AdminRole
│ }
└─ On failure: Throw UnauthorizedException (401)
```

**Default Behavior (Dev Mode):**
If no valid token provided, returns default SUPER_ADMIN context:
```typescript
{
  admin_id: 0,
  admin_username: 'test-user',
  admin_role: AdminRole.SUPER_ADMIN
}
```

### 3. Validation Middleware

**Purpose**: Validate request payload structure and content

```
ValidationMiddleware.validatePushRequest(body)
├─ Validate title_identifier
│  ├─ Type: string
│  ├─ Max length: 100 characters
│  └─ Not empty
├─ Validate locale_contents
│  ├─ Must be object
│  ├─ Must have at least one locale
│  └─ For each locale:
│     ├─ Validate title: string, max 256 chars
│     ├─ Validate body: string, max 1000 chars
│     ├─ Validate icon_url: optional, URL format
│     └─ Validate action_url: optional, URL format
└─ Validate scheduled_at
   ├─ Optional field
   ├─ Must be ISO 8601 timestamp if provided
   └─ Can be null for immediate send
```

**Validation Errors:**
- Missing required fields → 400 Bad Request
- Invalid field types → 400 Bad Request
- Field length exceeded → 400 Bad Request
- Invalid locale format → 400 Bad Request

### 4. ViewModel Orchestration

**Purpose**: Business logic execution and model coordination

```
PushNotificationViewModel.pushToAll(
  title_identifier,
  locale_contents,
  auth,
  scheduled_at
)
├─ Step 1: Permission Check
│  ├─ Verify auth.admin_role
│  ├─ Require SUPER_ADMIN role for pushToAll
│  └─ Throw ForbiddenException if insufficient
│
├─ Step 2: Campaign Creation
│  └─ CampaignModel.create()
│     ├─ Generate UUID for campaign_id
│     ├─ Set status = DRAFT (initial)
│     ├─ Set created_by = auth.admin_id
│     ├─ Set created_at = NOW()
│     ├─ Set scheduled_at = provided value or null
│     └─ Return Campaign object
│
├─ Step 3: Content Localization
│  └─ For each [locale_id, content] in locale_contents:
│     ├─ CampaignModel.createNotificationContent()
│     ├─ Generate UUID for content_id
│     ├─ Insert notification_content record:
│     │  ├─ campaign_id (FK to campaigns)
│     │  ├─ locale_id (FK to locales)
│     │  ├─ title
│     │  ├─ body
│     │  ├─ icon_url (optional)
│     │  └─ action_url (optional)
│     └─ Return NotificationContent object
│
├─ Step 4: Subscriber Lookup
│  └─ SubscriptionModel.getUserSubscriptions()
│     ├─ Query web_push_subscriptions
│     ├─ Filter by is_active = true
│     ├─ For pushToUser: Filter by user_id
│     ├─ For pushToSegment: Join with user_segments
│     ├─ For pushToTopic: Join with user_topics
│     ├─ For pushToAll: No filter
│     └─ Return array of subscriptions
│
└─ Step 5: Delivery Logging
   └─ For each subscription:
      ├─ NotificationLogModel.create()
      ├─ Generate UUID for log_id
      ├─ Insert notification_log record:
      │  ├─ campaign_id (FK)
      │  ├─ subscription_id (FK)
      │  ├─ user_id
      │  ├─ sent_by = auth.admin_id
      │  ├─ status = QUEUED/SENT/FAILED
      │  ├─ http_status_code (from delivery)
      │  ├─ error_message (if failed)
      │  └─ created_at = NOW()
      └─ Return NotificationLog object
```

### 5. Campaign Status Management

After creation, campaign status is determined:

```
IF scheduled_at is provided:
  status = QUEUED (waiting for scheduled time)
ELSE:
  status = SENDING (immediate delivery)
THEN:
  Update campaign.status in database
```

### 6. Database Persistence Layer

All operations use Prisma ORM with transaction support:

```
Database Operations:
├─ notificationCampaign.create()
│  └─ INSERT INTO notification_campaigns (...)
├─ notificationContent.create() [repeated for each locale]
│  └─ INSERT INTO notification_contents (...)
├─ webPushSubscription.findMany()
│  └─ SELECT FROM web_push_subscriptions WHERE ...
└─ notificationLog.create() [repeated for each subscription]
   └─ INSERT INTO notification_logs (...)

Constraints Enforced:
├─ Foreign Key: campaigns.created_by → admins.id
├─ Foreign Key: contents.campaign_id → campaigns.id
├─ Foreign Key: contents.locale_id → locales.id
├─ Foreign Key: logs.campaign_id → campaigns.id
├─ Foreign Key: logs.subscription_id → subscriptions.id
├─ Unique: (campaign_id, locale_id) in contents
└─ Cascade Delete: When campaign deleted, contents & logs deleted
```

### 7. Response Generation

```
HTTP 201 Created

Response Body: {
  "campaign_id": "550e8400-e29b-41d4-a716-446655440000",
  "total_subscribers": 1523,
  "timestamp": "2026-05-19T10:00:00.000Z"
}

Response Headers:
├─ Content-Type: application/json
├─ Content-Length: <calculated>
└─ Date: <server timestamp>
```

## Critical Decision Points

### 1. Authentication Decision
- **Has Bearer Token?** → Validate token
- **No Token?** → Use default SUPER_ADMIN (dev mode only)
- **Invalid Token?** → Return 401 Unauthorized

### 2. Role Check
- **SUPER_ADMIN?** → Allowed for pushToAll
- **ADMIN?** → Allowed for pushToUser, pushToSegment, pushToTopic
- **Insufficient Role?** → Return 403 Forbidden

### 3. Subscriber Count
- **Found subscribers?** → Create log entries for each
- **No subscribers?** → Valid scenario, continue (returns 0)
- **Database error?** → Return 500 with error details

### 4. Scheduling Decision
- **scheduled_at provided?** → Set status to QUEUED
- **scheduled_at null?** → Set status to SENDING
- **Invalid timestamp?** → Caught by validation middleware (400)

## Error Handling Flow

```
Request Processing
├─ AuthMiddleware throws UnauthorizedException
│  └─ ErrorHandler catches → 401 Response
├─ ValidationMiddleware throws ValidationException
│  └─ ErrorHandler catches → 400 Response
├─ ViewModel throws ForbiddenException
│  └─ ErrorHandler catches → 403 Response
├─ Database constraint violation
│  └─ ErrorHandler catches → 500 Response
├─ Unexpected errors
│  └─ ErrorHandler catches → 500 Response
└─ All errors logged with:
   ├─ Timestamp
   ├─ Request ID
   ├─ Stack trace
   ├─ Error message
   └─ Context information
```

## Performance Characteristics

| Operation | Time | Notes |
|-----------|------|-------|
| Auth validation | < 1ms | Token map lookup |
| Input validation | < 5ms | Field-by-field check |
| Campaign creation | 10-20ms | UUID generation + DB insert |
| Content creation | 5-10ms per locale | Multiple inserts (parallel possible) |
| Subscription lookup | 20-50ms | Depends on subscription count |
| Log creation | 2-5ms per subscription | Batched if possible |
| **Total per request** | **50-150ms** | Varies with subscriber count |

## Database Transactions

```
BEGIN TRANSACTION
├─ Insert campaign
├─ Insert content records (one per locale)
├─ Insert log records (one per subscription)
├─ Update campaign status
└─ COMMIT or ROLLBACK

Rollback Triggers:
├─ Foreign key constraint violation
├─ Unique constraint violation
├─ Database connection error
└─ Transaction timeout
```

## Logging Points

```
1. Request Start
   - Request ID generated
   - Timestamp recorded
   - Request parameters logged

2. Auth Check
   - Token validation result
   - Admin context extracted

3. Validation
   - Validation results
   - Any errors encountered

4. Business Logic
   - Campaign created
   - Content records created
   - Subscriptions queried
   - Logs created

5. Response
   - HTTP status code
   - Response body
   - Total execution time
   - Any warnings or errors
```

## Security Considerations

1. **Token Validation**: Bearer token checked before any processing
2. **Role Enforcement**: SUPER_ADMIN check before broadcast
3. **Input Sanitization**: Validation prevents injection attacks
4. **Audit Trail**: All operations logged for compliance
5. **Error Messages**: Production errors don't expose internals
6. **Database Security**: Prisma prevents SQL injection
7. **Foreign Keys**: Database enforces referential integrity
