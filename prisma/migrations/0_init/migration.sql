-- CreateEnum: AdminRole
CREATE TYPE "AdminRole" AS ENUM ('SUPER_ADMIN', 'ADMIN');

-- CreateEnum: CampaignStatus
CREATE TYPE "CampaignStatus" AS ENUM ('DRAFT', 'QUEUED', 'SENDING', 'COMPLETED', 'FAILED');

-- CreateEnum: LogStatus
CREATE TYPE "LogStatus" AS ENUM ('QUEUED', 'SENT', 'FAILED_EXPIRED', 'FAILED_ERROR');

-- CreateEnum: BrowserPlatform
CREATE TYPE "BrowserPlatform" AS ENUM ('CHROME', 'FIREFOX', 'SAFARI', 'EDGE');

-- CreateTable: admins
CREATE TABLE "admins" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "username" VARCHAR(50) NOT NULL UNIQUE,
    "role" "AdminRole" NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" INTEGER
);

-- CreateTable: locales
CREATE TABLE "locales" (
    "id" VARCHAR(5) NOT NULL PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL
);

-- CreateTable: web_push_subscriptions
CREATE TABLE "web_push_subscriptions" (
    "id" UUID NOT NULL PRIMARY KEY,
    "user_id" VARCHAR(64) NOT NULL,
    "browser" "BrowserPlatform" NOT NULL,
    "endpoint" TEXT NOT NULL,
    "p256dh_key" VARCHAR(255) NOT NULL,
    "auth_key" VARCHAR(255) NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL
);

-- CreateTable: notification_campaigns
CREATE TABLE "notification_campaigns" (
    "id" UUID NOT NULL PRIMARY KEY,
    "title_identifier" VARCHAR(100) NOT NULL,
    "status" "CampaignStatus" NOT NULL DEFAULT 'DRAFT',
    "created_by" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "scheduled_at" TIMESTAMP(3)
);

-- CreateTable: notification_contents
CREATE TABLE "notification_contents" (
    "id" UUID NOT NULL PRIMARY KEY,
    "campaign_id" UUID NOT NULL,
    "locale_id" VARCHAR(5) NOT NULL,
    "title" VARCHAR(120) NOT NULL,
    "body" TEXT NOT NULL,
    "icon_url" VARCHAR(512),
    "action_url" VARCHAR(512)
);

-- CreateTable: notification_logs
CREATE TABLE "notification_logs" (
    "id" UUID NOT NULL PRIMARY KEY,
    "campaign_id" UUID NOT NULL,
    "subscription_id" UUID NOT NULL,
    "user_id" VARCHAR(64) NOT NULL,
    "sent_by" INTEGER NOT NULL,
    "status" "LogStatus" NOT NULL,
    "http_status_code" INTEGER,
    "error_message" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateTable: user_segments
CREATE TABLE "user_segments" (
    "user_id" VARCHAR(64) NOT NULL,
    "segment_id" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("user_id", "segment_id")
);

-- CreateTable: user_topics
CREATE TABLE "user_topics" (
    "user_id" VARCHAR(64) NOT NULL,
    "topic_id" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("user_id", "topic_id")
);

-- CreateIndex
CREATE INDEX "web_push_subscriptions_user_id_idx" ON "web_push_subscriptions"("user_id");

-- CreateIndex
CREATE INDEX "web_push_subscriptions_is_active_idx" ON "web_push_subscriptions"("is_active");

-- CreateIndex
CREATE INDEX "web_push_subscriptions_user_id_is_active_idx" ON "web_push_subscriptions"("user_id", "is_active");

-- CreateIndex
CREATE INDEX "notification_campaigns_status_idx" ON "notification_campaigns"("status");

-- CreateIndex
CREATE INDEX "notification_campaigns_created_at_idx" ON "notification_campaigns"("created_at");

-- CreateIndex
CREATE INDEX "notification_contents_campaign_id_idx" ON "notification_contents"("campaign_id");

-- CreateIndex
CREATE UNIQUE INDEX "notification_contents_campaign_id_locale_id_key" ON "notification_contents"("campaign_id", "locale_id");

-- CreateIndex
CREATE INDEX "notification_logs_campaign_id_idx" ON "notification_logs"("campaign_id");

-- CreateIndex
CREATE INDEX "notification_logs_subscription_id_idx" ON "notification_logs"("subscription_id");

-- CreateIndex
CREATE INDEX "notification_logs_user_id_idx" ON "notification_logs"("user_id");

-- CreateIndex
CREATE INDEX "notification_logs_status_idx" ON "notification_logs"("status");

-- CreateIndex
CREATE INDEX "notification_logs_created_at_idx" ON "notification_logs"("created_at");

-- CreateIndex
CREATE INDEX "user_segments_segment_id_idx" ON "user_segments"("segment_id");

-- CreateIndex
CREATE INDEX "user_topics_topic_id_idx" ON "user_topics"("topic_id");

-- AddForeignKey
ALTER TABLE "admins" ADD CONSTRAINT "admins_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "admins"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "notification_campaigns" ADD CONSTRAINT "notification_campaigns_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "admins"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "notification_contents" ADD CONSTRAINT "notification_contents_campaign_id_fkey" FOREIGN KEY ("campaign_id") REFERENCES "notification_campaigns"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "notification_contents" ADD CONSTRAINT "notification_contents_locale_id_fkey" FOREIGN KEY ("locale_id") REFERENCES "locales"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "notification_logs" ADD CONSTRAINT "notification_logs_campaign_id_fkey" FOREIGN KEY ("campaign_id") REFERENCES "notification_campaigns"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "notification_logs" ADD CONSTRAINT "notification_logs_subscription_id_fkey" FOREIGN KEY ("subscription_id") REFERENCES "web_push_subscriptions"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "notification_logs" ADD CONSTRAINT "notification_logs_sent_by_fkey" FOREIGN KEY ("sent_by") REFERENCES "admins"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- Insert seed data: Locales
INSERT INTO "locales" ("id", "name") VALUES ('en', 'English');
INSERT INTO "locales" ("id", "name") VALUES ('zh', 'Chinese (Simplified)');
INSERT INTO "locales" ("id", "name") VALUES ('ms', 'Malay');
INSERT INTO "locales" ("id", "name") VALUES ('tm', 'Tamil');

-- Insert seed data: Sample admins
INSERT INTO "admins" ("name", "username", "role") VALUES ('Super Admin', 'superadmin', 'SUPER_ADMIN');
INSERT INTO "admins" ("name", "username", "role") VALUES ('Regular Admin', 'admin', 'ADMIN');
