-- Admin users table
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('SUPER_ADMIN', 'ADMIN')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by INT REFERENCES admins(id)
);

-- Locales table
CREATE TABLE locales (
    id VARCHAR(5) PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

-- Web push subscriptions table
CREATE TABLE web_push_subscriptions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    browser VARCHAR(20) NOT NULL CHECK (browser IN ('CHROME', 'FIREFOX', 'SAFARI', 'EDGE')),
    endpoint TEXT NOT NULL,
    p256dh_key VARCHAR(255) NOT NULL,
    auth_key VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for subscriptions
CREATE INDEX idx_subscriptions_user_id ON web_push_subscriptions(user_id);
CREATE INDEX idx_subscriptions_is_active ON web_push_subscriptions(is_active);
CREATE INDEX idx_subscriptions_user_active ON web_push_subscriptions(user_id, is_active);

-- Notification campaigns table
CREATE TABLE notification_campaigns (
    id UUID PRIMARY KEY,
    title_identifier VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('DRAFT', 'QUEUED', 'SENDING', 'COMPLETED', 'FAILED')) DEFAULT 'DRAFT',
    created_by INT NOT NULL REFERENCES admins(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    scheduled_at TIMESTAMP
);

-- Create indexes for campaigns
CREATE INDEX idx_campaigns_status ON notification_campaigns(status);
CREATE INDEX idx_campaigns_created_at ON notification_campaigns(created_at);

-- Notification contents table
CREATE TABLE notification_contents (
    id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL REFERENCES notification_campaigns(id) ON DELETE CASCADE,
    locale_id VARCHAR(5) NOT NULL REFERENCES locales(id),
    title VARCHAR(120) NOT NULL,
    body TEXT NOT NULL,
    icon_url VARCHAR(512),
    action_url VARCHAR(512),
    UNIQUE(campaign_id, locale_id)
);

-- Create indexes for contents
CREATE INDEX idx_contents_campaign_id ON notification_contents(campaign_id);

-- Notification logs table
CREATE TABLE notification_logs (
    id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL REFERENCES notification_campaigns(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES web_push_subscriptions(id),
    user_id VARCHAR(64) NOT NULL,
    sent_by INT NOT NULL REFERENCES admins(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('QUEUED', 'SENT', 'FAILED_EXPIRED', 'FAILED_ERROR')),
    http_status_code INT,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for logs
CREATE INDEX idx_logs_campaign_id ON notification_logs(campaign_id);
CREATE INDEX idx_logs_subscription_id ON notification_logs(subscription_id);
CREATE INDEX idx_logs_user_id ON notification_logs(user_id);
CREATE INDEX idx_logs_status ON notification_logs(status);
CREATE INDEX idx_logs_created_at ON notification_logs(created_at);

-- User segments table (for push to group functionality)
CREATE TABLE user_segments (
    user_id VARCHAR(64) NOT NULL,
    segment_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, segment_id)
);

-- Create index for segments
CREATE INDEX idx_segments_segment_id ON user_segments(segment_id);

-- User topics table (for push to topic functionality)
CREATE TABLE user_topics (
    user_id VARCHAR(64) NOT NULL,
    topic_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, topic_id)
);

-- Create index for topics
CREATE INDEX idx_topics_topic_id ON user_topics(topic_id);

-- Insert sample locales
INSERT INTO locales (id, name) VALUES
    ('en', 'English'),
    ('zh', 'Chinese (Simplified)'),
    ('ms', 'Malay'),
    ('tm', 'Tamil');

-- Insert sample admins
INSERT INTO admins (name, username, role) VALUES
    ('Super Admin', 'superadmin', 'SUPER_ADMIN'),
    ('Regular Admin', 'admin', 'ADMIN');
