import { WebPushSubscription, BrowserPlatform } from '../types';
import db from '../config/database';

export class SubscriptionModel {
  async create(
    userId: string,
    browser: BrowserPlatform,
    endpoint: string,
    p256dhKey: string,
    authKey: string
  ): Promise<WebPushSubscription> {
    console.log('[SubscriptionModel.create] Creating subscription:', { userId, browser, endpoint });
    const result = await db.webPushSubscription.create({
      data: {
        userId,
        browser,
        endpoint,
        p256dhKey,
        authKey,
        isActive: true
      }
    });
    console.log('[SubscriptionModel.create] Created:', { id: result.id, userId: result.userId });
    return result;
  }

  async getById(id: string): Promise<WebPushSubscription | null> {
    return db.webPushSubscription.findUnique({
      where: { id }
    });
  }

  async getActiveByUserId(userId: string): Promise<WebPushSubscription[]> {
    return db.webPushSubscription.findMany({
      where: {
        userId,
        isActive: true
      }
    });
  }

  async getActiveBySegment(segmentId: string): Promise<WebPushSubscription[]> {
    return db.webPushSubscription.findMany({
      where: {
        isActive: true,
        NOT: {
          userId: ''
        }
      },
      include: {
        logs: true
      }
    }).then(subs => {
      return db.userSegment.findMany({
        where: { segmentId },
        select: { userId: true }
      }).then(segments => {
        const segmentUserIds = new Set(segments.map(s => s.userId));
        return subs.filter(sub => segmentUserIds.has(sub.userId));
      });
    });
  }

  async getActiveByTopic(topicId: string): Promise<WebPushSubscription[]> {
    return db.userTopic.findMany({
      where: { topicId },
      select: { userId: true }
    }).then(topics => {
      const topicUserIds = new Set(topics.map(t => t.userId));
      return db.webPushSubscription.findMany({
        where: {
          isActive: true,
          userId: { in: Array.from(topicUserIds) }
        }
      });
    });
  }

  async getAll(active = true): Promise<WebPushSubscription[]> {
    return db.webPushSubscription.findMany({
      where: active ? { isActive: true } : undefined
    });
  }

  async markInactive(id: string): Promise<void> {
    await db.webPushSubscription.update({
      where: { id },
      data: { isActive: false }
    });
  }

  async update(id: string, updates: Partial<WebPushSubscription>): Promise<WebPushSubscription> {
    return db.webPushSubscription.update({
      where: { id },
      data: updates
    });
  }
}
