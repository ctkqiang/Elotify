# Database: Real PostgreSQL Data

## Critical Notice

**This service uses a REAL PostgreSQL database.** All API responses contain actual data persisted in the database. There is NO mock data hardcoded in the application. Every request reads from and writes to PostgreSQL.

## What This Means

| Operation | Result |
|-----------|--------|
| POST `/push/all` | Inserts campaign into `notification_campaigns` table |
| GET `/notification/latest` | Selects from `notification_campaigns` table |
| PUT `/notification/:id` | Updates `notification_campaigns` table |
| DELETE `/notification/:id` | Deletes from `notification_campaigns` table |
| Any GET request | Queries PostgreSQL database |
| Any POST request | Persists to PostgreSQL database |

## Data Flow

```
Postman/Client Request
         ↓
    API Endpoint
         ↓
  ViewModel (Business Logic)
         ↓
    Model Layer (Data Access)
         ↓
  Prisma ORM (Type-safe queries)
         ↓
PostgreSQL Database ← REAL DATA
         ↓
Response (Actual data from DB)
         ↓
  Postman/Client
```

## Verification

### 1. Via Postman (Easy)

```
GET /notification/latest
```

This returns actual campaigns you created via POST requests. The response contains:
- Real campaign UUIDs
- Real timestamps from `created_at`
- Real status values from database
- Real subscriber counts from queries

### 2. Via Direct Database Query (Advanced)

```powershell
# Connect to PostgreSQL
docker exec valufarmpushnoification-postgres-1 psql -U postgres -d push_notifications

# View all campaigns
SELECT id, title_identifier, status, created_at FROM notification_campaigns;

# View campaign content
SELECT * FROM notification_contents;

# View delivery logs
SELECT * FROM notification_logs;
```

All data visible in Postman responses matches exactly what's in the database.

## Database Schema

**8 Tables - All Contains Real Data:**

1. **admins** - Admin user accounts
   - Default: `admin_id: 0, username: 'test-user', role: 'SUPER_ADMIN'`
   - Data: Inserted when you start (if needed)

2. **locales** - Supported languages
   - Data: Pre-inserted (en, zh, ms, tm)
   - Real reference data

3. **notification_campaigns** - Campaign master records
   - Data: Created every time you POST to `/push/all`, `/push/user`, `/push/group`, `/push/topic`
   - Real UUIDs, real timestamps, real status values

4. **notification_contents** - Localized content per campaign
   - Data: Created for each locale in your request
   - Real content you provided in request body

5. **web_push_subscriptions** - Browser push endpoints
   - Data: Empty by default (you can add via SQL)
   - Represents real user devices

6. **notification_logs** - Delivery audit trail
   - Data: Created for each subscription when campaign sent
   - Real HTTP status codes and error messages

7. **user_segments** - User group membership
   - Data: Empty by default
   - Can be populated for segment-based delivery

8. **user_topics** - Topic subscriptions
   - Data: Empty by default
   - Can be populated for topic-based delivery

## Data Persistence

**Your data is PERSISTENT:**

```
[YES] Data survives container restarts
[YES] Data survives application redeployment
[YES] Data persists across Postman sessions
[YES] PostgreSQL volume is mounted: postgres_data
```

When you restart docker-compose:
```powershell
docker-compose down
docker-compose up
```

All your campaigns, content, and logs still exist in the database.

## Working with Real Data

### Creating Real Data

```bash
# POST /push/all
{
  "title_identifier": "maintenance_alert",
  "locale_contents": {
    "en": {"title": "Maintenance", "body": "2 hours downtime"},
    "zh": {"title": "维护", "body": "停机2小时"}
  }
}
```

This creates:
- 1 row in `notification_campaigns`
- 2 rows in `notification_contents` (one per locale)
- N rows in `notification_logs` (one per subscriber)

### Reading Real Data

```bash
# GET /notification/latest
```

Returns the actual campaigns you created, with real UUIDs and timestamps.

### Updating Real Data

```bash
# PUT /notification/{id}
{"status": "SENDING"}
```

Updates the actual campaign status in the database.

### Deleting Real Data

```bash
# DELETE /notification/{id}
```

Deletes the actual campaign and all related records (cascade delete).

## Prisma ORM Integration

The application uses **Prisma ORM** which provides:

```typescript
// Type-safe queries
const campaigns = await db.notificationCampaign.findMany();

// Parameterized queries (SQL injection prevention)
const campaign = await db.notificationCampaign.findUnique({
  where: { id: campaignId }
});

// Transactions (ACID compliance)
await db.$transaction([
  db.notificationCampaign.create({ data: {...} }),
  db.notificationContent.create({ data: {...} })
]);

// Relationships
const campaignWithContent = await db.notificationCampaign.findUnique({
  where: { id: campaignId },
  include: { contents: true, logs: true }
});
```

All real database operations with zero mock data.

## No Mock Data Anywhere

### What's NOT in the code:

```typescript
// [NO] NOT hardcoded responses
return {
  campaign_id: "fake-uuid-12345",
  total_subscribers: 999
};

// [NO] NOT fake data arrays
const mockCampaigns = [
  { id: "1", title: "Test", status: "DRAFT" }
];

// [NO] NOT mock databases
const fakeDb = { campaigns: [] };
```

### What IS in the code:

```typescript
// [YES] Real database queries
const campaigns = await db.notificationCampaign.findMany();

// [YES] Real inserts
await db.notificationCampaign.create({ data: {...} });

// [YES] Real updates
await db.notificationCampaign.update({
  where: { id },
  data: { status }
});
```

## Performance Characteristics

Since you're reading REAL data from PostgreSQL:

| Operation | Time | Notes |
|-----------|------|-------|
| GET /notification/latest | 10-50ms | Database query with limit |
| POST /push/all | 50-200ms | Inserts campaign + content + logs |
| PUT /notification/:id | 5-20ms | Single update |
| GET with many results | 50-500ms | Depends on data size |

Times vary based on actual data volume in your database.

## Data Consistency

```
Every API response contains:
[YES] Real data from PostgreSQL
[YES] Current timestamp (not mocked)
[YES] Actual UUIDs (not generated in code)
[YES] Real foreign key relationships
[YES] Actual aggregate counts (not hardcoded)
```

## Resetting Data

To delete all real data and start fresh:

```powershell
# Delete database volume
docker-compose down -v

# Restart (creates empty database)
docker-compose up

# Re-insert locales and admin
docker exec valufarmpushnoification-postgres-1 psql -U postgres -d push_notifications << 'EOF'
INSERT INTO locales (id, name) VALUES ('en', 'English'), ('zh', 'Chinese'), ('ms', 'Malay'), ('tm', 'Tamil');
INSERT INTO admins (id, username, name, role) VALUES (0, 'test-user', 'Test Admin', 'SUPER_ADMIN');
EOF
```

Now you have a fresh database with empty campaigns.

## Testing with Real Data

### Workflow:

1. Create campaign via Postman POST → Stored in database
2. Query campaign via Postman GET → Retrieved from database
3. Update campaign via Postman PUT → Updated in database
4. Delete campaign via Postman DELETE → Removed from database
5. Verify via direct database query → Same data

All data is persistent and real. No mocks anywhere.
