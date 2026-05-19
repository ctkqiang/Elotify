import { Campaign, NotificationLog } from '../types';

// Request DTOs
export interface PushRequest {
  title_identifier: string;
  locale_contents: {
    [localeId: string]: {
      title: string;
      body: string;
      icon_url?: string;
      action_url?: string;
    };
  };
  scheduled_at?: string;
}

export interface UpdateNotificationRequest {
  status?: string;
}

// Response DTOs
export interface CampaignResponse {
  id: string;
  title_identifier: string;
  status: string;
  created_by: number;
  created_at: string;
  scheduled_at?: string;
  content_count?: number;
}

export interface NotificationLogResponse {
  id: string;
  campaign_id: string;
  subscription_id: string;
  user_id: string;
  sent_by: number;
  status: string;
  http_status_code?: number;
  error_message?: string;
  created_at: string;
}

export interface PushResponseDto {
  campaign_id: string;
  total_subscribers: number;
  timestamp: string;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
}

export interface ErrorResponse {
  error: string;
  message: string;
  status_code: number;
  timestamp: string;
}

export interface StatsResponse {
  total: number;
  sent: number;
  failed: number;
  queued: number;
}

// Mappers
export class PushNotificationViewMapper {
  static mapCampaignToResponse(campaign: Campaign): CampaignResponse {
    return {
      id: campaign.id,
      title_identifier: campaign.title_identifier,
      status: campaign.status,
      created_by: campaign.created_by,
      created_at: campaign.created_at.toISOString(),
      scheduled_at: campaign.scheduled_at?.toISOString()
    };
  }

  static mapLogToResponse(log: NotificationLog): NotificationLogResponse {
    return {
      id: log.id,
      campaign_id: log.campaign_id,
      subscription_id: log.subscription_id,
      user_id: log.user_id,
      sent_by: log.sent_by,
      status: log.status,
      http_status_code: log.http_status_code,
      error_message: log.error_message,
      created_at: log.created_at.toISOString()
    };
  }

  static mapPushResponse(campaignId: string, totalSubscribers: number): PushResponseDto {
    return {
      campaign_id: campaignId,
      total_subscribers: totalSubscribers,
      timestamp: new Date().toISOString()
    };
  }

  static mapHealthResponse(): HealthResponse {
    return {
      status: 'ok',
      timestamp: new Date().toISOString(),
      version: '1.0.0'
    };
  }

  static mapErrorResponse(error: Error, statusCode = 500): ErrorResponse {
    return {
      error: error.name || 'Error',
      message: error.message,
      status_code: statusCode,
      timestamp: new Date().toISOString()
    };
  }
}
