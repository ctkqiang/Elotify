import { NotificationLog, LogStatus } from '../types';
import db from '../config/database';

export class NotificationLogModel {
  async create(
    campaignId: string,
    subscriptionId: string,
    userId: string,
    sentBy: number,
    status: LogStatus = LogStatus.QUEUED,
    httpStatusCode?: number,
    errorMessage?: string
  ): Promise<NotificationLog> {
    return db.notificationLog.create({
      data: {
        campaignId,
        subscriptionId,
        userId,
        sentBy,
        status,
        httpStatusCode: httpStatusCode || undefined,
        errorMessage: errorMessage || undefined
      }
    });
  }

  async getById(id: string): Promise<NotificationLog | null> {
    return db.notificationLog.findUnique({
      where: { id }
    });
  }

  async updateStatus(
    id: string,
    status: LogStatus,
    httpStatusCode?: number,
    errorMessage?: string
  ): Promise<NotificationLog> {
    return db.notificationLog.update({
      where: { id },
      data: {
        status,
        httpStatusCode: httpStatusCode || undefined,
        errorMessage: errorMessage || undefined
      }
    });
  }

  async getByCampaignId(campaignId: string, limit = 100): Promise<NotificationLog[]> {
    return db.notificationLog.findMany({
      where: { campaignId },
      take: limit,
      orderBy: { createdAt: 'desc' }
    });
  }

  async getByUserId(userId: string, limit = 50): Promise<NotificationLog[]> {
    return db.notificationLog.findMany({
      where: { userId },
      take: limit,
      orderBy: { createdAt: 'desc' }
    });
  }

  async getStatsByCampaign(campaignId: string): Promise<{
    total: number;
    sent: number;
    failed: number;
    queued: number;
  }> {
    const logs = await db.notificationLog.findMany({
      where: { campaignId },
      select: { status: true }
    });

    const stats = {
      total: logs.length,
      sent: logs.filter((l: { status: LogStatus }) => l.status === LogStatus.SENT).length,
      failed: logs.filter((l: { status: LogStatus }) => l.status === LogStatus.FAILED_EXPIRED || l.status === LogStatus.FAILED_ERROR).length,
      queued: logs.filter((l: { status: LogStatus }) => l.status === LogStatus.QUEUED).length
    };

    return stats;
  }
}
