# Testing Guide: Push Notifications

This guide shows how to test the push notification system end-to-end.

## Complete Workflow

### Step 1: Create a Test Subscription

```bash
curl -X POST http://localhost:3000/admin/subscriptions \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "browser": "CHROME",
    "endpoint": "https://push.example.com/endpoint/sub-123",
    "p256dh_key": "test-public-key-base64",
    "auth_key": "test-auth-secret-base64"
  }'
```

**Expected Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440010",
  "user_id": "user-123",
  "browser": "CHROME",
  "endpoint": "https://push.example.com/endpoint/sub-123",
  "is_active": true,
  "created_at": "2026-05-19T10:00:00.000Z",
  "timestamp": "2026-05-19T10:00:00.000Z"
}
```

### Step 2: Push Notification to That User

```bash
curl -X POST http://localhost:3000/push/user/user-123 \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "title_identifier": "test_notification",
    "locale_contents": {
      "en": {
        "title": "Test Title",
        "body": "This is a test notification"
      }
    }
  }'
```

**Expected Response (201 Created):**
```json
{
  "campaign_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user-123",
  "timestamp": "2026-05-19T10:00:00.000Z"
}
```

### Step 3: Retrieve Notifications for That User

```bash
curl -X GET http://localhost:3000/notification/user/user-123 \
  -H "Authorization: Bearer admin-token-regular"
```

**Expected Response (200 OK):**
```json
{
  "user_id": "user-123",
  "notifications": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "campaign_id": "550e8400-e29b-41d4-a716-446655440000",
      "sent_at": "2026-05-19T10:00:00.000Z",
      "status": "completed"
    }
  ],
  "total_count": 1,
  "unread_count": 0
}
```

## What Changed

The key issue was that `GET /notification/user/user-123` returned empty results because:

1. **Before**: No subscription existed for user-123
2. When pushing → subscription lookup found 0 subscriptions
3. No notification logs were created (logs are only created for real subscriptions)
4. GET endpoint returns empty

Now with the new `POST /admin/subscriptions` endpoint:

1. Create a subscription for user-123
2. Push notification to user-123
3. Logs are created and returned by GET endpoint

## Testing Multiple Users

Create subscriptions for multiple users:

```bash
# User 1
curl -X POST http://localhost:3000/admin/subscriptions \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-alice",
    "browser": "FIREFOX",
    "endpoint": "https://push.example.com/endpoint/alice-001",
    "p256dh_key": "key1",
    "auth_key": "auth1"
  }'

# User 2
curl -X POST http://localhost:3000/admin/subscriptions \
  -H "Authorization: Bearer admin-token-regular" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-bob",
    "browser": "CHROME",
    "endpoint": "https://push.example.com/endpoint/bob-001",
    "p256dh_key": "key2",
    "auth_key": "auth2"
  }'

# Push to all subscriptions
curl -X POST http://localhost:3000/push/all \
  -H "Authorization: Bearer admin-token-super" \
  -H "Content-Type: application/json" \
  -d '{
    "title_identifier": "broadcast_message",
    "locale_contents": {
      "en": {
        "title": "Broadcast",
        "body": "Message to all users"
      }
    }
  }'
```

## Database Setup

If you need to insert test data directly into PostgreSQL:

```sql
-- Insert test subscription
INSERT INTO web_push_subscriptions (
  id, user_id, browser, endpoint, p256dh_key, auth_key, is_active, created_at, updated_at
) VALUES (
  gen_random_uuid(),
  'user-123',
  'CHROME',
  'https://push.example.com/endpoint/example',
  'test-key',
  'test-auth',
  true,
  NOW(),
  NOW()
);

-- Verify subscription exists
SELECT id, user_id, browser, is_active FROM web_push_subscriptions WHERE user_id = 'user-123';
```

## Postman Collection

A Postman collection file is available at `docs/postman.json` with pre-configured requests.

Import and use:
1. Set `base_url` variable to `http://localhost:3000`
2. Set `token_admin` variable to `admin-token-regular`
3. Run requests in order:
   - Create Subscription
   - Push Notification
   - Get User Notifications
