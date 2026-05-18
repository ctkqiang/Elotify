# Push Notification Service

A production-grade push notification service built with **ElysiaJS**, **TypeScript**, **Prisma ORM**, and **MVVM architecture** for delivering web push notifications at scale.

## Features

### Push Notification Capabilities

| Feature | Description |
|---------|-------------|
| Push to All Users | Broadcast notifications to entire user base |
| Push to Specific User | Target individual users by ID |
| Push to Segments | Send to user groups/segments |
| Push to Topics | Topic-based subscription delivery |

### Internationalization (i18n)

| Aspect | Details |
|--------|---------|
| Locale Support | Multi-language content delivery |
| Built-in Locales | en (English), zh (Chinese), ms (Malay), tm (Tamil) |
| Content Versioning | Campaign content per locale |
| Unicode Support | Full UTF-8 character support |

### Comprehensive Logging & Debugging

| Capability | Implementation |
|------------|-----------------|
| Operation Logging | All operations logged with full context |
| Request Tracing | Correlation IDs for request tracking |
| Performance Metrics | Duration tracking on every operation |
| Error Information | Complete stack traces in logs |
| Audit Trail | User action and database query logging |
| Easy Debugging | Contextual information for quick issue identification |

### Security Features

| Security Layer | Implementation |
|---|---|
| Authentication | Bearer token validation |
| Authorization | Role-based access control (SUPER_ADMIN, ADMIN) |
| Input Validation | Comprehensive validation and sanitization |
| Error Handling | Centralized error handling with safe messages |
| SQL Injection Prevention | Prisma parameterized queries |
| Permission Tracking | Logging of all authorization checks |
| Compliance | Audit trails for regulatory compliance |

### Campaign Management

| Feature | Details |
|---------|---------|
| Workflow | DRAFT → QUEUED → SENDING → COMPLETED/FAILED |
| Scheduling | Support for time-based campaign delivery |
| Localization | Content versioning per locale |
| Analytics | Delivery tracking and statistics |
| Logging | Comprehensive delivery logs for debugging |

### Scalable Architecture

| Component | Benefit |
|-----------|---------|
| MVVM Pattern | Clear separation of concerns |
| Layered Design | Independent testing of each layer |
| Prisma ORM | Type-safe database queries |
| Abstraction Layer | Database independence |
| Extensibility | Easy to add features |
| Error Handling | Production-ready resilience |

## Project Structure

### Directory Layout

```
src/
├── config/              # Configuration (database, env)
│   └── database.ts      # Database connection
├── middleware/          # Cross-cutting concerns
│   ├── auth.middleware.ts        # Authentication & authorization
│   ├── validation.middleware.ts  # Input validation
│   └── error.middleware.ts       # Error handling
├── models/              # Data layer (M in MVVM)
│   ├── campaign.model.ts           # Campaign operations
│   ├── subscription.model.ts       # Subscription management
│   └── notification-log.model.ts   # Delivery logs
├── viewmodels/          # Business logic (VM in MVVM)
│   └── push-notification.viewmodel.ts  # Core business logic
├── views/               # Presentation layer (V in MVVM)
│   └── push-notification.view.ts       # DTOs & mappers
├── types/               # TypeScript types
│   └── index.ts         # Shared types
├── routes/              # API endpoints
│   └── notification.routes.ts  # Route handlers
└── index.ts             # Application entry point
```

### Directory Responsibilities

| Directory | Purpose | Key Files |
|-----------|---------|-----------|
| `config/` | Application configuration | `database.ts` (Prisma setup) |
| `middleware/` | Cross-cutting concerns | Auth, validation, error handling |
| `models/` | Database operations (M) | Campaign, subscription, notification log |
| `viewmodels/` | Business logic (VM) | Core notification orchestration |
| `views/` | Response formatting (V) | DTOs and response mappers |
| `types/` | Type definitions | Enums and interfaces |
| `routes/` | API endpoints | All HTTP route handlers |
| `prisma/` | ORM configuration | Schema and migrations |

## Architecture: MVVM Pattern

### MVVM Layer Comparison

| Layer | Purpose | Responsibilities | Key Files |
|-------|---------|------------------|-----------|
| **Model** | Data Access | CRUD ops, persistence, query optimization, transactions | `models/*.ts` |
| **ViewModel** | Business Logic | Rule enforcement, permission checks, orchestration, error handling | `viewmodels/*.ts` |
| **View** | API Contracts | Request/response formatting, DTOs, data mapping, serialization | `views/*.ts` |

### Model Layer (`src/models/`)
Handles direct database operations without business logic.

```typescript
// Example: CampaignModel
class CampaignModel {
  async create(title, createdBy): Campaign { }
  async getById(id): Campaign { }
  async updateStatus(id, status): Campaign { }
}
```

**Responsibilities:**
- Database CRUD operations
- Data persistence
- Query optimization
- Transaction management

### ViewModel Layer (`src/viewmodels/`)
Contains all business logic and orchestrates models.

```typescript
// Example: PushNotificationViewModel
class PushNotificationViewModel {
  async pushToAll(title, contents, auth): Promise<{ campaign_id, total }> {
    // 1. Validate permissions
    // 2. Create campaign
    // 3. Add translated content
    // 4. Query subscriptions
    // 5. Queue notifications
    // 6. Return response
  }
}
```

**Responsibilities:**
- Business rule enforcement
- Permission validation
- Model orchestration
- Error handling
- Data transformation

### View Layer (`src/views/`)
Defines DTOs and response formatting.

```typescript
// Request DTO
interface PushRequest {
  title_identifier: string;
  locale_contents: { [localeId]: { title, body, icon_url?, action_url? } };
  scheduled_at?: string;
}

// Response DTO
interface PushResponseDto {
  campaign_id: string;
  total_subscribers: number;
  timestamp: string;
}
```

**Responsibilities:**
- Request/response contracts
- Data mapping
- Serialization/deserialization
- API documentation

## Installation

### Prerequisites
- Node.js 18+
- PostgreSQL 12+
- npm or bun

### Setup

1. **Clone and install**
```bash
git clone <repository>
cd ValuFarm.PushNoification
npm install
# or
bun install
```

2. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your database credentials
# DATABASE_URL=postgresql://user:password@localhost:5432/push_notifications
```

3. **Initialize database with Prisma**
```bash
# Generate Prisma types
npm run prisma:generate
# or
bun run prisma:generate

# Create database schema and apply migrations
npm run prisma:migrate
# or
bun run prisma:migrate
```

4. **Start development server**
```bash
npm run dev
# or
bun run dev
```

Server runs at `http://localhost:3000`

### Viewing Logs

When the server starts, you'll see comprehensive logs with full context:

```
[2024-05-18T10:30:45.123Z] [INFO] [APP_STARTUP] Application started successfully | port: 3000 | startupTimeMscompletion: 331ms
[2024-05-18T10:30:46.234Z] [INFO] [HTTP_REQUEST] HTTP Request: POST /push/all | statusCode: 201 | duration: 145ms
[2024-05-18T10:30:47.345Z] [INFO] [CAMPAIGN] Campaign Event: CREATED for campaign 'uuid-123'
```

All operations are logged with:
- Timestamp (ISO 8601)
- Log level (DEBUG, INFO, WARN, ERROR)
- Module name
- Full context data
- Request/User IDs
- Performance metrics

## Quick Start

For detailed setup instructions, see [START_HERE.md](START_HERE.md)

### 1. Health Check
```bash
curl http://localhost:3000/health
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2024-05-18T10:30:45.123Z",
  "version": "1.0.0"
}
```

### 2. Push to All Users
```bash
curl -X POST http://localhost:3000/push/all \
  -H "Authorization: Bearer admin-token-super" \
  -H "Content-Type: application/json" \
  -d '{
    "title_identifier": "welcome_campaign",
    "locale_contents": {
      "en": {
        "title": "Welcome!",
        "body": "Thanks for joining us",
        "action_url": "https://example.com/welcome"
      },
      "zh": {
        "title": "欢迎!",
        "body": "感谢加入我们"
      }
    }
  }'
```

### 3. Get Latest Notifications
```bash
curl http://localhost:3000/notification/latest
```

## Logging System

This service includes comprehensive logging on every operation for easy debugging.

Each log entry includes:
- Timestamp (ISO 8601 format)
- Log level (DEBUG, INFO, WARN, ERROR)
- Module name for context
- Detailed message
- Request ID for tracing
- User ID for tracking
- Operation context (all relevant data)
- Duration metrics (in milliseconds)
- Error stack traces (when applicable)

This makes identifying and fixing issues straightforward - all relevant context is logged automatically.

## API Documentation

See [docs/API.md](docs/API.md) for complete API reference.

### Key Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/push/all` | Push to all users |
| POST | `/push/user/{user_id}` | Push to specific user |
| POST | `/push/group/{segment_id}` | Push to user segment |
| POST | `/push/topic/{topic_id}` | Push to topic subscribers |
| GET | `/notification/latest` | Get latest campaigns |
| GET | `/notification/user/{user_id}` | Get user notifications |
| GET | `/notification/group/{segment_id}` | Get segment notifications |
| PUT | `/notification/{id}` | Update notification status |
| DELETE | `/notification/{id}` | Delete notification |

## Security

### Authentication
All endpoints require Bearer token authentication (except `/health` and `/`).

**Development tokens:**
- `admin-token-super` → SUPER_ADMIN role
- `admin-token-regular` → ADMIN role

**For production:** Implement JWT with RS256 or OAuth2

### Authorization
Role-based access control:

| Role | Capabilities |
|------|--------------|
| SUPER_ADMIN | Create, update, delete, push to all users |
| ADMIN | Push to user/group/topic, view data only |

### Input Validation
- Type checking
- Field length validation
- Format validation (ISO dates, URLs)
- Sanitization

### Error Handling
- Structured error responses
- HTTP status codes
- Request logging
- Stack traces in development

## Database Schema

### Key Tables

**admins**: Admin users with roles
```sql
CREATE TABLE admins (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE,
  role VARCHAR(20) CHECK (role IN ('SUPER_ADMIN', 'ADMIN')),
  created_at TIMESTAMP
);
```

**notification_campaigns**: Campaign lifecycle
```sql
CREATE TABLE notification_campaigns (
  id UUID PRIMARY KEY,
  title_identifier VARCHAR(100),
  status VARCHAR(20) CHECK (status IN ('DRAFT', 'QUEUED', 'SENDING', 'COMPLETED', 'FAILED')),
  created_by INT REFERENCES admins(id),
  scheduled_at TIMESTAMP
);
```

**web_push_subscriptions**: User subscriptions
```sql
CREATE TABLE web_push_subscriptions (
  id UUID PRIMARY KEY,
  user_id VARCHAR(64),
  browser VARCHAR(20),
  endpoint TEXT,
  p256dh_key VARCHAR(255),
  auth_key VARCHAR(255),
  is_active BOOLEAN
);
```

**notification_logs**: Delivery tracking
```sql
CREATE TABLE notification_logs (
  id UUID PRIMARY KEY,
  campaign_id UUID REFERENCES notification_campaigns(id),
  subscription_id UUID REFERENCES web_push_subscriptions(id),
  status VARCHAR(20) CHECK (status IN ('QUEUED', 'SENT', 'FAILED_EXPIRED', 'FAILED_ERROR')),
  http_status_code INT,
  error_message TEXT
);
```

See [docs/schema.sql](docs/schema.sql) for complete schema.

## Logging & Debugging

All operations are logged with comprehensive context for easy debugging.

### Logging Levels

Set `LOG_LEVEL` environment variable:

| Level | Purpose | Use Case |
|-------|---------|----------|
| DEBUG | Detailed information | Development and troubleshooting |
| INFO | Normal operation | Default, general information |
| WARN | Warnings and alerts | Potential issues to monitor |
| ERROR | Error events only | Production error tracking |

### Logging Format

Set `LOG_FORMAT` environment variable:

| Format | Description | Use Case |
|--------|-------------|----------|
| text | Human-readable output | Development, easier to read |
| json | Structured JSON logs | Production, log aggregation systems |

### Example Debug Session

```bash
# Enable debug logging
LOG_LEVEL=DEBUG bun run dev

# View database operations
LOG_LEVEL=DEBUG npm run dev | grep DATABASE

# View authentication events
LOG_LEVEL=DEBUG npm run dev | grep AUTHENTICATION
```

### Log Output Examples

Database Operation:
```
[2024-05-18T10:30:49.567Z] [DEBUG] [DATABASE] Database Operation: SELECT on table 'notification_campaigns' | durationMs: 12 | rowsAffected: 25
```

Business Operation:
```
[2024-05-18T10:30:48.456Z] [INFO] [BUSINESS_OPERATION] Business Operation: PUSH_TO_ALL on resource 'campaign-uuid' - Status: SUCCESS | Duration: 234ms
```

Error with Stack Trace:
```
[2024-05-18T10:30:50.678Z] [ERROR] [BUSINESS_OPERATION] Operation failed: Database error | Error Name: PrismaClientKnownRequestError | Stack Trace: [...]
```

## Development

### Database GUI
Open Prisma Studio to explore and edit database:
```bash
npm run prisma:studio
# or
bun run prisma:studio
```

Opens at: `http://localhost:5555`

### Database Migrations
Create and apply database migrations:
```bash
# Create migration
npm run prisma:migrate

# View migration status
npm run prisma:migrate status
```

### Running Tests
```bash
npm test
# or
bun test
```

### Building for Production
```bash
npm run build
# or
bun run build
```

### Environment Variables
See [.env.example](.env.example) for all available configuration options.

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| DATABASE_URL | - | YES | PostgreSQL connection string |
| PORT | 3000 | NO | Server port number |
| NODE_ENV | development | NO | Environment (development/production) |
| LOG_LEVEL | INFO | NO | Logging level (DEBUG/INFO/WARN/ERROR) |
| LOG_FORMAT | text | NO | Log format (text/json) |

## Architecture Decisions

### Why MVVM?

| Benefit | Explanation |
|---------|-------------|
| Separation of Concerns | Models, ViewModels, and Views are independent and focused |
| Testability | Each layer can be tested in isolation without dependencies |
| Scalability | Easy to add new features without affecting existing code |
| Maintainability | Clear responsibilities for each layer, easy to understand |
| Reusability | ViewModels can be reused across different view implementations |
| Debugging | Issues can be isolated to specific layers quickly |

### Why Elysia?

| Advantage | Description |
|-----------|-------------|
| Performance | Fast TypeScript web framework with minimal overhead |
| Type Safety | Full TypeScript support with excellent developer experience |
| Simplicity | Minimal boilerplate, easy to learn and use |
| Extensibility | Plugin system for middleware and custom features |
| Modern | Built for modern Node.js patterns and best practices |
| Community | Active community with good documentation |

## What's Included

| Feature | Status | Implementation |
|---------|--------|-----------------|
| Comprehensive Logging | [YES] | All operations logged with full context |
| MVVM Architecture | [YES] | Model-View-ViewModel with Prisma ORM |
| Push Notifications | [YES] | Multi-targeted (all, user, segment, topic) |
| Role-Based Access Control | [YES] | SUPER_ADMIN and ADMIN roles |
| Input Validation | [YES] | Comprehensive validation and sanitization |
| Error Handling | [YES] | Centralized error handling and recovery |
| Database Migrations | [YES] | Prisma migrations with version control |
| Type Safety | [YES] | Type-safe queries throughout |
| Request Tracing | [YES] | Correlation IDs for request tracking |
| Performance Metrics | [YES] | Duration tracking on all operations |
| Authentication Logging | [YES] | Complete auth event tracking |
| Delivery Tracking | [YES] | Campaign and delivery analytics |
| Complete Documentation | [YES] | 15+ comprehensive guide documents |

## Future Enhancements

- [ ] Batch notification processing
- [ ] Real-time delivery status updates via WebSocket
- [ ] Advanced segmentation/targeting
- [ ] A/B testing for notification content
- [ ] Analytics dashboard
- [ ] Retry logic for failed deliveries
- [ ] Rate limiting per admin
- [ ] Redis caching layer
- [ ] GraphQL API (in addition to REST)
- [ ] Webhook support for external integrations

## Contributing

1. Follow the MVVM pattern
2. Add input validation for new endpoints
3. Include error handling
4. Update API documentation
5. Write tests for new features

## License

MIT

## Documentation

Complete documentation is provided:

- [START_HERE.md](START_HERE.md) - Quick start (5 minutes)
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Developer cheat sheet
- [docs/API.md](docs/API.md) - Complete API reference
- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) - Developer guide and MVVM patterns
- [docs/ORM_SETUP.md](docs/ORM_SETUP.md) - Prisma ORM usage guide
- [SECURITY.md](SECURITY.md) - Security practices and guidelines
- [TECHNICAL_SPECIFICATIONS.md](TECHNICAL_SPECIFICATIONS.md) - System specs and deployment
- [GIT_COMMITS.md](GIT_COMMITS.md) - Git workflow and commit standards
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - Database schema and design
- [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md) - Complete documentation map

## Support

For issues and questions:
1. Check [START_HERE.md](START_HERE.md) for setup issues
2. Check [QUICK_REFERENCE.md](QUICK_REFERENCE.md) for common tasks
3. Check [docs/API.md](docs/API.md) for API questions
4. Check [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md) for all documentation
5. Create an issue in the repository with logs (see Logging & Debugging section)
