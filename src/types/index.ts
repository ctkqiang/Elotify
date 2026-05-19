/**
 * @fileoverview Type Definitions Module
 *
 * Provides TypeScript interfaces and enums for all domain models.
 * Ensures type safety across the application.
 *
 * @module src/types
 */

/**
 * Admin user roles for role-based access control (RBAC)
 *
 * @enum {string}
 * @readonly
 */
export enum AdminRole {
  /** Full system access: create, update, delete, push to all */
  SUPER_ADMIN = 'SUPER_ADMIN',
  /** Limited access: push to user/group/topic, view data */
  ADMIN = 'ADMIN'
}

/**
 * Admin user information
 *
 * @interface Admin
 * @property {number} id - Unique identifier
 * @property {string} name - Full name
 * @property {string} username - Unique username
 * @property {AdminRole} role - User role
 * @property {Date} created_at - Creation timestamp
 * @property {number} [created_by] - Creator admin ID
 */
export interface Admin {
  id: number;
  name: string;
  username: string;
  role: AdminRole;
  created_at: Date;
  created_by?: number;
}

/**
 * Campaign lifecycle states
 *
 * Flow: DRAFT -> QUEUED -> SENDING -> COMPLETED/FAILED
 *
 * @enum {string}
 * @readonly
 */
export enum CampaignStatus {
  /** Initial state: campaign being prepared */
  DRAFT = 'DRAFT',
  /** Queued for sending: waiting for scheduler */
  QUEUED = 'QUEUED',
  /** Currently being sent to subscribers */
  SENDING = 'SENDING',
  /** Successfully delivered to all recipients */
  COMPLETED = 'COMPLETED',
  /** Failed: error during sending process */
  FAILED = 'FAILED'
}

/**
 * Push notification campaign batch
 *
 * @interface Campaign
 * @property {string} id - UUID identifier
 * @property {string} titleIdentifier - Internal campaign name
 * @property {CampaignStatus} status - Current state
 * @property {number} createdBy - Creator admin ID
 * @property {Date} createdAt - Creation timestamp
 * @property {Date} [scheduledAt] - Scheduled send time (optional)
 */
export interface Campaign {
  id: string;
  titleIdentifier: string;
  status: CampaignStatus;
  createdBy: number;
  createdAt: Date;
  scheduledAt?: Date | null;
}

/**
 * Translated notification content for a specific locale
 *
 * @interface NotificationContent
 * @property {string} id - UUID identifier
 * @property {string} campaignId - Associated campaign ID
 * @property {string} localeId - Language locale (e.g., 'en', 'zh')
 * @property {string} title - Notification title (max 120 chars)
 * @property {string} body - Notification message body
 * @property {string} [iconUrl] - Icon URL (HTTPS only in production)
 * @property {string} [actionUrl] - Click destination URL
 */
export interface NotificationContent {
  id: string;
  campaignId: string;
  localeId: string;
  title: string;
  body: string;
  iconUrl?: string | null;
  actionUrl?: string | null;
}

/**
 * Browser platforms for push subscription
 *
 * @enum {string}
 * @readonly
 */
export enum BrowserPlatform {
  CHROME = 'CHROME',
  FIREFOX = 'FIREFOX',
  SAFARI = 'SAFARI',
  EDGE = 'EDGE'
}

/**
 * User browser push subscription endpoint
 *
 * Stores Web Push API subscription details for a specific browser instance.
 * Can be active or inactive (410 Gone response handling).
 *
 * @interface WebPushSubscription
 * @property {string} id - UUID identifier
 * @property {string} userId - External user identifier
 * @property {BrowserPlatform} browser - Browser type
 * @property {string} endpoint - Push service URL
 * @property {string} p256dhKey - Client public key for encryption
 * @property {string} authKey - Client authentication secret
 * @property {boolean} isActive - Subscription active status
 * @property {Date} createdAt - Subscription creation time
 * @property {Date} updatedAt - Last update time
 */
export interface WebPushSubscription {
  id: string;
  userId: string;
  browser: BrowserPlatform;
  endpoint: string;
  p256dhKey: string;
  authKey: string;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * Notification delivery attempt status
 *
 * @enum {string}
 * @readonly
 */
export enum LogStatus {
  /** Queued for delivery */
  QUEUED = 'QUEUED',
  /** Successfully delivered */
  SENT = 'SENT',
  /** Failed: subscription expired (410 Gone) */
  FAILED_EXPIRED = 'FAILED_EXPIRED',
  /** Failed: other error (timeout, network, etc.) */
  FAILED_ERROR = 'FAILED_ERROR'
}

/**
 * Notification delivery attempt log
 *
 * Records every delivery attempt with status and HTTP response code.
 * Used for analytics, retry logic, and audit trails.
 *
 * @interface NotificationLog
 * @property {string} id - UUID identifier
 * @property {string} campaignId - Campaign being delivered
 * @property {string} subscriptionId - Target subscription
 * @property {string} userId - Target user (cached for queries)
 * @property {number} sentBy - Admin ID who triggered send
 * @property {LogStatus} status - Delivery status
 * @property {number} [httpStatusCode] - Push service HTTP response
 * @property {string} [errorMessage] - Error details if failed
 * @property {Date} createdAt - Delivery attempt timestamp
 */
export interface NotificationLog {
  id: string;
  campaignId: string;
  subscriptionId: string;
  userId: string;
  sentBy: number;
  status: LogStatus;
  httpStatusCode?: number | null;
  errorMessage?: string | null;
  createdAt: Date;
}

/**
 * Authenticated user context
 *
 * Contains extracted authentication and authorization information
 * for the current request. Populated from JWT or bearer token.
 *
 * @interface AuthContext
 * @property {number} adminId - Authenticated admin user ID
 * @property {string} adminUsername - Admin username
 * @property {AdminRole} adminRole - User role for RBAC
 */
export interface AuthContext {
  adminId: number;
  adminUsername: string;
  adminRole: AdminRole;
}

/**
 * Supported language/locale
 *
 * @interface Locale
 * @property {string} id - Locale code (e.g., 'en', 'zh', 'ms')
 * @property {string} name - Locale display name
 */
export interface Locale {
  id: string;
  name: string;
}
