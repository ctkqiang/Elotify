import { PushRequest, UpdateNotificationRequest } from '../views/push-notification.view';

export class ValidationMiddleware {
  static validatePushRequest(body: any): PushRequest {
    if (!body.title_identifier || typeof body.title_identifier !== 'string') {
      throw new ValidationException('title_identifier is required and must be a string');
    }

    if (body.title_identifier.length > 100) {
      throw new ValidationException('title_identifier must not exceed 100 characters');
    }

    if (!body.locale_contents || typeof body.locale_contents !== 'object') {
      throw new ValidationException('locale_contents is required and must be an object');
    }

    const localeContents = body.locale_contents;
    for (const [localeId, content] of Object.entries(localeContents)) {
      if (typeof localeId !== 'string' || localeId.length > 5) {
        throw new ValidationException(`Invalid locale ID: ${localeId}`);
      }

      if (!content.title || typeof content.title !== 'string') {
        throw new ValidationException(`title is required for locale ${localeId}`);
      }

      if (content.title.length > 120) {
        throw new ValidationException(`title must not exceed 120 characters for locale ${localeId}`);
      }

      if (!content.body || typeof content.body !== 'string') {
        throw new ValidationException(`body is required for locale ${localeId}`);
      }

      if (content.icon_url && typeof content.icon_url !== 'string') {
        throw new ValidationException(`icon_url must be a string for locale ${localeId}`);
      }

      if (content.icon_url && content.icon_url.length > 512) {
        throw new ValidationException(`icon_url must not exceed 512 characters for locale ${localeId}`);
      }

      if (content.action_url && typeof content.action_url !== 'string') {
        throw new ValidationException(`action_url must be a string for locale ${localeId}`);
      }

      if (content.action_url && content.action_url.length > 512) {
        throw new ValidationException(`action_url must not exceed 512 characters for locale ${localeId}`);
      }
    }

    if (body.scheduled_at) {
      const date = new Date(body.scheduled_at);
      if (isNaN(date.getTime())) {
        throw new ValidationException('scheduled_at must be a valid ISO 8601 date string');
      }
      if (date < new Date()) {
        throw new ValidationException('scheduled_at cannot be in the past');
      }
    }

    return {
      title_identifier: body.title_identifier,
      locale_contents: body.locale_contents,
      scheduled_at: body.scheduled_at
    };
  }

  static validateUserId(userId: string): string {
    if (!userId || typeof userId !== 'string') {
      throw new ValidationException('user_id is required and must be a string');
    }

    if (userId.length > 64) {
      throw new ValidationException('user_id must not exceed 64 characters');
    }

    return userId;
  }

  static validateSegmentId(segmentId: string): string {
    if (!segmentId || typeof segmentId !== 'string') {
      throw new ValidationException('segment_id is required and must be a string');
    }

    return segmentId;
  }

  static validateTopicId(topicId: string): string {
    if (!topicId || typeof topicId !== 'string') {
      throw new ValidationException('topic_id is required and must be a string');
    }

    return topicId;
  }

  static validateNotificationId(notificationId: string): string {
    if (!notificationId || typeof notificationId !== 'string') {
      throw new ValidationException('notification_id is required and must be a valid UUID');
    }

    return notificationId;
  }

  static validateUpdateRequest(body: any): UpdateNotificationRequest {
    if (body.status) {
      const validStatuses = ['DRAFT', 'QUEUED', 'SENDING', 'COMPLETED', 'FAILED'];
      if (!validStatuses.includes(body.status)) {
        throw new ValidationException(`status must be one of: ${validStatuses.join(', ')}`);
      }
    }

    return {
      status: body.status
    };
  }
}

export class ValidationException extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ValidationException';
  }
}
