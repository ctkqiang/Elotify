# Push Notification Service Architecture

## Data Structure

```plantuml
@startuml
!theme plain
skinparam linestyle orthographic
skinparam roundcorner 8
skinparam shadowing false

enum AdminRole {
    SUPER_ADMIN
    ADMIN
}

enum CampaignStatus {
    DRAFT
    QUEUED
    SENDING
    COMPLETED
    FAILED
}

enum LogStatus {
    QUEUED
    SENT
    FAILED_EXPIRED
    FAILED_ERROR
}

enum BrowserPlatform {
    CHROME
    FIREFOX
    SAFARI
    EDGE
}

entity "admins" as admins {
    * id : INT <<PK>>
    --
    * name : VARCHAR(100)
    * username : VARCHAR(50) <<UNIQUE>>
    * role : AdminRole
    * created_at : TIMESTAMP
    created_by : INT <<FK to admins.id>>
}

entity "locales" as locales {
    * id : VARCHAR(5) <<PK>> -- e.g., 'en', 'zh', 'ms', 'tm'
    --
    * name : VARCHAR(50)     -- e.g., 'English', 'Chinese'
}

entity "web_push_subscriptions" as subscriptions {
    * id : UUID <<PK>>
    --
    * user_id : VARCHAR(64) <<FK/Index>> -- Maps to your external User service
    * browser : BrowserPlatform
    * endpoint : TEXT                    -- Browser push service gateway URL
    * p256dh_key : VARCHAR(255)         -- Client public crypto key
    * auth_key : VARCHAR(255)           -- Client auth secret
    * is_active : BOOLEAN                -- Housekeeping flag (False on 410 Gone)
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
}

entity "notification_campaigns" as campaigns {
    * id : UUID <<PK>>
    --
    * title_identifier : VARCHAR(100)   -- Internal name e.g. "maintenance_alert"
    * status : CampaignStatus
    * created_by : INT <<FK to admins.id>>
    * created_at : TIMESTAMP
    scheduled_at : TIMESTAMP
}

entity "notification_contents" as contents {
    * id : UUID <<PK>>
    --
    * campaign_id : UUID <<FK to campaigns.id>>
    * locale_id : VARCHAR(5) <<FK to locales.id>>
    --
    * title : VARCHAR(120)               -- Visible notification title
    * body : TEXT                        -- Visible notification message
    icon_url : VARCHAR(512)              -- Image asset displayed in browser UI
    action_url : VARCHAR(512)            -- Target URL when user clicks notification
}

entity "notification_logs" as logs {
    * id : UUID <<PK>>
    --
    * campaign_id : UUID <<FK to campaigns.id>>
    * subscription_id : UUID <<FK to subscriptions.id>>
    * user_id : VARCHAR(64) <<Index>>    -- Cached for quick querying
    * sent_by : INT <<FK to admins.id>>  -- Audit context
    --
    * status : LogStatus
    http_status_code : INT               -- Direct response from browser vendor (e.g., 201, 404, 410)
    error_message : TEXT
    * created_at : TIMESTAMP
}

' Relationships
admins ||--o{ admins : "creates"
admins ||--o{ campaigns : "creates"
admins ||--o{ logs : "triggers (audited)"

campaigns ||--|{ contents : "has translations"
locales ||--o{ contents : "utilizes"

campaigns ||--o{ logs : "generates executions"
subscriptions ||--o{ logs : "receives delivery attempts"

@enduml
```

## Entities

### admins
- Stores admin users with role-based access control
- Supports audit trail via `created_by` foreign key
- Roles: SUPER_ADMIN, ADMIN

### locales
- Enables multi-language support for notification content
- Standard locale identifiers (e.g., 'en', 'zh', 'ms', 'tm')

### web_push_subscriptions
- Stores browser push service endpoints for each user
- Includes encryption keys required by Web Push API (p256dh, auth)
- `is_active` flag for handling 410 Gone responses from push services

### notification_campaigns
- Represents a batch of notifications to be sent
- Lifecycle: DRAFT → QUEUED → SENDING → COMPLETED/FAILED
- Supports scheduled delivery

### notification_contents
- Translated content for each campaign and locale
- Includes title, body, icon, and action URL

### notification_logs
- Tracks every delivery attempt
- Records HTTP status from browser vendors
- Supports user-level and campaign-level analytics
- Cached `user_id` for query efficiency
