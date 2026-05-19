import { CampaignModel } from '../models/campaign.model';
import { SubscriptionModel } from '../models/subscription.model';
import { NotificationLogModel } from '../models/notification-log.model';
import {
  Campaign,
  CampaignStatus,
  WebPushSubscription,
  NotificationLog,
  LogStatus,
  AuthContext,
  AdminRole
} from '../types';

export class PushNotificationViewModel {
  constructor(
    private campaignModel: CampaignModel,
    private subscriptionModel: SubscriptionModel,
    private logModel: NotificationLogModel
  ) {}

  // Campaign management
  async createCampaign(
    titleIdentifier: string,
    auth: AuthContext,
    scheduledAt?: Date
  ): Promise<Campaign> {
    if (auth.admin_role !== AdminRole.SUPER_ADMIN && auth.admin_role !== AdminRole.ADMIN) {
      throw new Error('Unauthorized: Only SUPER_ADMIN or ADMIN can create campaigns');
    }

    return this.campaignModel.create(titleIdentifier, auth.admin_id, scheduledAt);
  }

  async getCampaign(id: string): Promise<Campaign> {
    const campaign = await this.campaignModel.getById(id);
    if (!campaign) {
      throw new Error('Campaign not found');
    }
    return campaign;
  }

  // Push to all users
  async pushToAll(
    titleIdentifier: string,
    localeContents: { [localeId: string]: { title: string; body: string; icon_url?: string; action_url?: string } },
    auth: AuthContext,
    scheduledAt?: Date
  ): Promise<{ campaign_id: string; total_subscribers: number }> {
    const campaign = await this.createCampaign(titleIdentifier, auth, scheduledAt);

    for (const [localeId, content] of Object.entries(localeContents)) {
      await this.campaignModel.addContent(
        campaign.id,
        localeId,
        content.title,
        content.body,
        content.icon_url,
        content.action_url
      );
    }

    const subscriptions = await this.subscriptionModel.getAll(true);
    await this.queueNotifications(campaign.id, subscriptions, auth.admin_id);

    const finalStatus = scheduledAt ? CampaignStatus.QUEUED : CampaignStatus.SENDING;
    await this.campaignModel.updateStatus(campaign.id, finalStatus);

    return {
      campaign_id: campaign.id,
      total_subscribers: subscriptions.length
    };
  }

  // Push to specific user
  async pushToUser(
    userId: string,
    titleIdentifier: string,
    localeContents: { [localeId: string]: { title: string; body: string; icon_url?: string; action_url?: string } },
    auth: AuthContext,
    scheduledAt?: Date
  ): Promise<{ campaign_id: string; total_subscribers: number }> {
    const campaign = await this.createCampaign(titleIdentifier, auth, scheduledAt);

    for (const [localeId, content] of Object.entries(localeContents)) {
      await this.campaignModel.addContent(
        campaign.id,
        localeId,
        content.title,
        content.body,
        content.icon_url,
        content.action_url
      );
    }

    const subscriptions = await this.subscriptionModel.getActiveByUserId(userId);

    await this.queueNotifications(campaign.id, subscriptions, auth.admin_id);

    const finalStatus = scheduledAt ? CampaignStatus.QUEUED : CampaignStatus.SENDING;
    await this.campaignModel.updateStatus(campaign.id, finalStatus);

    return {
      campaign_id: campaign.id,
      total_subscribers: subscriptions.length
    };
  }

  // Push to segment/group
  async pushToSegment(
    segmentId: string,
    titleIdentifier: string,
    localeContents: { [localeId: string]: { title: string; body: string; icon_url?: string; action_url?: string } },
    auth: AuthContext,
    scheduledAt?: Date
  ): Promise<{ campaign_id: string; total_subscribers: number }> {
    const campaign = await this.createCampaign(titleIdentifier, auth, scheduledAt);

    for (const [localeId, content] of Object.entries(localeContents)) {
      await this.campaignModel.addContent(
        campaign.id,
        localeId,
        content.title,
        content.body,
        content.icon_url,
        content.action_url
      );
    }

    const subscriptions = await this.subscriptionModel.getActiveBySegment(segmentId);

    await this.queueNotifications(campaign.id, subscriptions, auth.admin_id);

    const finalStatus = scheduledAt ? CampaignStatus.QUEUED : CampaignStatus.SENDING;
    await this.campaignModel.updateStatus(campaign.id, finalStatus);

    return {
      campaign_id: campaign.id,
      total_subscribers: subscriptions.length
    };
  }

  // Push to topic
  async pushToTopic(
    topicId: string,
    titleIdentifier: string,
    localeContents: { [localeId: string]: { title: string; body: string; icon_url?: string; action_url?: string } },
    auth: AuthContext,
    scheduledAt?: Date
  ): Promise<{ campaign_id: string; total_subscribers: number }> {
    const campaign = await this.createCampaign(titleIdentifier, auth, scheduledAt);

    for (const [localeId, content] of Object.entries(localeContents)) {
      await this.campaignModel.addContent(
        campaign.id,
        localeId,
        content.title,
        content.body,
        content.icon_url,
        content.action_url
      );
    }

    const subscriptions = await this.subscriptionModel.getActiveByTopic(topicId);

    await this.queueNotifications(campaign.id, subscriptions, auth.admin_id);

    const finalStatus = scheduledAt ? CampaignStatus.QUEUED : CampaignStatus.SENDING;
    await this.campaignModel.updateStatus(campaign.id, finalStatus);

    return {
      campaign_id: campaign.id,
      total_subscribers: subscriptions.length
    };
  }

  // Get notifications
  async getLatestNotifications(limit = 10): Promise<Campaign[]> {
    return this.campaignModel.getAll(limit, 0);
  }

  async getUserNotifications(userId: string, limit = 50): Promise<NotificationLog[]> {
    return this.logModel.getByUserId(userId, limit);
  }

  async getSegmentNotifications(segmentId: string, limit = 50): Promise<Campaign[]> {
    return this.campaignModel.getAll(limit, 0);
  }

  async getTopicNotifications(topicId: string, limit = 50): Promise<Campaign[]> {
    return this.campaignModel.getAll(limit, 0);
  }

  // Update notification status
  async updateNotificationStatus(
    notificationId: string,
    status: CampaignStatus,
    auth: AuthContext
  ): Promise<Campaign> {
    if (auth.admin_role !== AdminRole.SUPER_ADMIN) {
      throw new Error('Unauthorized: Only SUPER_ADMIN can update notifications');
    }

    const notification = await this.campaignModel.getById(notificationId);
    if (!notification) {
      throw new Error('Notification not found');
    }

    return this.campaignModel.updateStatus(notificationId, status);
  }

  // Update and delete
  async updateNotification(
    notificationId: string,
    updates: Partial<Campaign>,
    auth: AuthContext
  ): Promise<Campaign> {
    if (auth.admin_role !== AdminRole.SUPER_ADMIN) {
      throw new Error('Unauthorized: Only SUPER_ADMIN can update notifications');
    }

    const notification = await this.campaignModel.getById(notificationId);
    if (!notification) {
      throw new Error('Notification not found');
    }

    if (notification.status !== CampaignStatus.DRAFT) {
      throw new Error('Can only update notifications in DRAFT status');
    }

    // Update status if provided
    if (updates.status) {
      return this.campaignModel.updateStatus(notificationId, updates.status);
    }

    return notification;
  }

  async deleteNotification(notificationId: string, auth: AuthContext): Promise<void> {
    if (auth.admin_role !== AdminRole.SUPER_ADMIN) {
      throw new Error('Unauthorized: Only SUPER_ADMIN can delete notifications');
    }

    const notification = await this.campaignModel.getById(notificationId);
    if (!notification) {
      throw new Error('Notification not found');
    }

    if (notification.status !== CampaignStatus.DRAFT) {
      throw new Error('Can only delete notifications in DRAFT status');
    }

    // Soft delete by updating status to FAILED
    await this.campaignModel.updateStatus(notificationId, CampaignStatus.FAILED);
  }

  // Health check
  async healthCheck(): Promise<{ status: string; timestamp: Date }> {
    // Add database connectivity check here
    return {
      status: 'ok',
      timestamp: new Date()
    };
  }

  // Private helper methods
  private async queueNotifications(
    campaignId: string,
    subscriptions: WebPushSubscription[],
    adminId: number
  ): Promise<void> {
    for (const subscription of subscriptions) {
      await this.logModel.create(
        campaignId,
        subscription.id,
        subscription.user_id,
        adminId,
        LogStatus.QUEUED
      );
    }
  }
}
