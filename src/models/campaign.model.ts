import { Campaign, CampaignStatus, NotificationContent } from '../types';
import db from '../config/database';

export class CampaignModel {
  async create(
    titleIdentifier: string,
    createdBy: number,
    scheduledAt?: Date
  ): Promise<Campaign> {
    return db.notificationCampaign.create({
      data: {
        titleIdentifier,
        status: CampaignStatus.DRAFT,
        createdBy,
        scheduledAt: scheduledAt || undefined
      }
    });
  }

  async getById(id: string): Promise<Campaign | null> {
    return db.notificationCampaign.findUnique({
      where: { id }
    });
  }

  async updateStatus(id: string, status: CampaignStatus): Promise<Campaign> {
    return db.notificationCampaign.update({
      where: { id },
      data: { status }
    });
  }

  async getAll(limit = 50, offset = 0): Promise<Campaign[]> {
    return db.notificationCampaign.findMany({
      take: limit,
      skip: offset,
      orderBy: { createdAt: 'desc' }
    }) as Promise<Campaign[]>;
  }

  async getContents(campaignId: string): Promise<NotificationContent[]> {
    return db.notificationContent.findMany({
      where: { campaignId }
    });
  }

  async addContent(
    campaignId: string,
    localeId: string,
    title: string,
    body: string,
    iconUrl?: string,
    actionUrl?: string
  ): Promise<NotificationContent> {
    return db.notificationContent.create({
      data: {
        campaignId,
        localeId,
        title,
        body,
        iconUrl: iconUrl || undefined,
        actionUrl: actionUrl || undefined
      }
    });
  }
}
