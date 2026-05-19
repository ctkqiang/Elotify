/**
 * @fileoverview API Routes Module
 *
 * Defines all HTTP endpoints for the Push Notification Service.
 * Handles:
 * - Request routing
 * - Input validation
 * - Response formatting
 * - Error handling and logging
 *
 * All endpoints are logged with comprehensive details for debugging.
 *
 * @module src/routes/notification.routes
 */

import { Elysia } from 'elysia';
import { PushNotificationViewModel } from '../viewmodels/push-notification.viewmodel';
import { CampaignModel } from '../models/campaign.model';
import { SubscriptionModel } from '../models/subscription.model';
import { NotificationLogModel } from '../models/notification-log.model';
import { AuthMiddleware } from '../middleware/auth.middleware';
import { ValidationMiddleware } from '../middleware/validation.middleware';
import { AdminRole } from '../types';
import { logger } from '../utils/logger';

/**
 * API Endpoint Contracts
 * These constants define the canonical endpoint paths
 */
export const PostPushNotifications = {
  pushToAll: '/push/all',
  pushToUser: '/push/user/:user_id',
  pushToSegment: '/push/group/:segment_id',
  pushToTopic: '/push/topic/:topic_id',
  createSubscription: '/admin/subscriptions',
};

export const GetPushNotifications = {
  getHealth: '/health',
  getLatestNotification: '/notification/latest',
  getUserNotifications: '/notification/user/:user_id',
  getSegmentNotifications: '/notification/group/:segment_id',
  getTopicNotifications: '/notification/topic/:topic_id',
};

export const PutPushNotifications = {
  updateNotificationById: '/notification/:notification_id',
};

export const DeletePushNotifications = {
  deleteNotificationById: '/notification/:notification_id',
};

/**
 * Setup all notification routes with middleware
 * Initializes models and handlers, registers all endpoints
 *
 * @param app - Elysia application instance
 * @returns The Elysia application with routes registered
 */
export function setupNotificationRoutes(app: Elysia) {
  logger.info('Setting up notification routes', 'ROUTES_SETUP');

  const campaignModel = new CampaignModel();
  const subscriptionModel = new SubscriptionModel();
  const logModel = new NotificationLogModel();
  const viewModel = new PushNotificationViewModel(
    campaignModel,
    subscriptionModel,
    logModel
  );

  /**
   * POST /push/all
   * Push notification to all active subscribers
   */
  app.post(PostPushNotifications.pushToAll, async ({ body, set, request }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Push to all request received', 'PUSH_ALL', { requestId });

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const validatedRequest = ValidationMiddleware.validatePushRequest(body);

      const result = await viewModel.pushToAll(
        validatedRequest.title_identifier,
        validatedRequest.locale_contents,
        auth,
        validatedRequest.scheduled_at
      );

      set.status = 201;
      logger.logOperation('PUSH_TO_ALL', result.campaign_id, 'SUCCESS', 0, {
        totalSubscribers: result.total_subscribers
      });

      return {
        campaign_id: result.campaign_id,
        total_subscribers: result.total_subscribers,
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      logger.error('Error in push/all endpoint', error as Error, 'PUSH_ALL', { requestId });
      set.status = 500;
      return {
        error: 'Error',
        message: (error as Error).message || 'Internal server error',
        status_code: 500,
        timestamp: new Date().toISOString()
      };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * POST /push/user/:user_id
   * Push notification to specific user
   */
  app.post(PostPushNotifications.pushToUser, async ({ body, set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Push to user request received', 'PUSH_USER', { requestId, userId: params.user_id });

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const validatedRequest = ValidationMiddleware.validatePushRequest(body);

      const result = await viewModel.pushToUser(
        params.user_id,
        validatedRequest.title_identifier,
        validatedRequest.locale_contents,
        auth,
        validatedRequest.scheduled_at
      );

      set.status = 201;
      logger.logOperation('PUSH_TO_USER', result.campaign_id, 'SUCCESS', 0, {
        userId: params.user_id,
        totalSubscribers: result.total_subscribers
      });

      return {
        campaign_id: result.campaign_id,
        user_id: params.user_id,
        total_subscribers: result.total_subscribers,
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      logger.error('Error in push/user endpoint', error as Error, 'PUSH_USER', { requestId });
      set.status = 500;
      return {
        error: 'Error',
        message: (error as Error).message || 'Internal server error',
        status_code: 500,
        timestamp: new Date().toISOString()
      };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * POST /push/group/:segment_id
   * Push notification to user segment
   */
  app.post(PostPushNotifications.pushToSegment, async ({ body, set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Push to segment request received', 'PUSH_SEGMENT', { requestId, segmentId: params.segment_id });

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const validatedRequest = ValidationMiddleware.validatePushRequest(body);

      const result = await viewModel.pushToSegment(
        params.segment_id,
        validatedRequest.title_identifier,
        validatedRequest.locale_contents,
        auth,
        validatedRequest.scheduled_at
      );

      set.status = 201;
      logger.logOperation('PUSH_TO_SEGMENT', result.campaign_id, 'SUCCESS', 0, {
        segmentId: params.segment_id,
        recipients: result.total_subscribers
      });

      return {
        campaign_id: result.campaign_id,
        segment_id: params.segment_id,
        estimated_recipients: result.total_subscribers,
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      logger.error('Error in push/segment endpoint', error as Error, 'PUSH_SEGMENT', { requestId });
      set.status = 500;
      return {
        error: 'Error',
        message: (error as Error).message || 'Internal server error',
        status_code: 500,
        timestamp: new Date().toISOString()
      };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * POST /push/topic/:topic_id
   * Push notification to topic subscribers
   */
  app.post(PostPushNotifications.pushToTopic, async ({ body, set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Push to topic request received', 'PUSH_TOPIC', { requestId, topicId: params.topic_id });

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const validatedRequest = ValidationMiddleware.validatePushRequest(body);

      const result = await viewModel.pushToTopic(
        params.topic_id,
        validatedRequest.title_identifier,
        validatedRequest.locale_contents,
        auth,
        validatedRequest.scheduled_at
      );

      set.status = 201;
      logger.logOperation('PUSH_TO_TOPIC', result.campaign_id, 'SUCCESS', 0, {
        topicId: params.topic_id,
        subscribers: result.total_subscribers
      });

      return {
        campaign_id: result.campaign_id,
        topic_id: params.topic_id,
        total_subscribers: result.total_subscribers,
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      logger.error('Error in push/topic endpoint', error as Error, 'PUSH_TOPIC', { requestId });
      set.status = 500;
      return {
        error: 'Error',
        message: (error as Error).message || 'Internal server error',
        status_code: 500,
        timestamp: new Date().toISOString()
      };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * POST /admin/subscriptions
   * Create a test web push subscription for a user
   */
  app.post(PostPushNotifications.createSubscription, async ({ body, set, request }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Create subscription request received', 'CREATE_SUB', { requestId });

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      if (auth.admin_role !== AdminRole.ADMIN && auth.admin_role !== AdminRole.SUPER_ADMIN) {
        set.status = 403;
        return {
          error: 'Forbidden',
          message: 'Only ADMIN role can create subscriptions',
          status_code: 403
        };
      }

      const { user_id, browser, endpoint, p256dh_key, auth_key } = body as any;

      if (!user_id || !browser || !endpoint || !p256dh_key || !auth_key) {
        set.status = 400;
        return {
          error: 'Bad Request',
          message: 'Missing required fields: user_id, browser, endpoint, p256dh_key, auth_key',
          status_code: 400
        };
      }

      const validBrowsers = ['CHROME', 'FIREFOX', 'SAFARI', 'EDGE'];
      if (!validBrowsers.includes(browser)) {
        set.status = 400;
        return {
          error: 'Bad Request',
          message: `Invalid browser: ${browser}. Must be one of: ${validBrowsers.join(', ')}`,
          status_code: 400
        };
      }

      logger.info('Creating subscription', 'CREATE_SUB', { user_id, browser, endpoint });

      const subscription = await subscriptionModel.create(
        user_id,
        browser as any,
        endpoint,
        p256dh_key,
        auth_key
      );

      logger.info('Subscription created successfully', 'CREATE_SUB', {
        id: subscription.id,
        userId: subscription.userId,
        browser: subscription.browser
      });

      logger.info('Subscription created', 'CREATE_SUB', {
        subscriptionId: subscription.id,
        userId: subscription.userId,
        browser: subscription.browser,
        isActive: subscription.isActive
      });

      set.status = 201;
      logger.logOperation('CREATE_SUB', subscription.id, 'SUCCESS', 0, {
        userId: user_id,
        browser
      });

      return {
        id: subscription.id,
        user_id: subscription.userId,
        browser: subscription.browser,
        endpoint: subscription.endpoint,
        is_active: subscription.isActive,
        created_at: subscription.createdAt,
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      logger.error('Error creating subscription', error as Error, 'CREATE_SUB', { requestId });
      set.status = 500;
      return {
        error: 'Error',
        message: (error as Error).message || 'Internal server error',
        status_code: 500,
        timestamp: new Date().toISOString()
      };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * GET /health
   * Health check endpoint
   */
  app.get(GetPushNotifications.getHealth, () => {
    logger.debug('Health check request', 'HEALTH_CHECK');
    return {
      status: 'ok',
      timestamp: new Date().toISOString(),
      version: '1.0.0'
    };
  });

  /**
   * GET /notification/latest
   * Get latest campaigns
   */
  app.get(GetPushNotifications.getLatestNotification, async ({ set }) => {
    try {
      logger.info('Latest notifications request', 'GET_LATEST');
      const notifications = await viewModel.getLatestNotifications(10);
      return {
        campaigns: notifications.map(n => ({
          id: n.id,
          title_identifier: n.title_identifier,
          status: n.status,
          created_at: n.created_at,
          created_by: n.created_by,
          total_sent: 0
        })),
        total_count: notifications.length,
        limit: 10
      };
    } catch (error) {
      const errorMsg = (error as Error).message || 'Internal server error';
      logger.error('Error fetching latest notifications', error as Error, 'GET_LATEST');
      set.status = 500;
      return {
        error: 'Internal server error',
        message: errorMsg,
        status_code: 500
      };
    }
  });

  /**
   * GET /notification/user/:user_id
   * Get notifications for specific user
   */
  app.get(GetPushNotifications.getUserNotifications, async ({ set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Get user notifications request', 'GET_USER_NOTIF', { userId: params.user_id });

      AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const notifications = await viewModel.getUserNotifications(params.user_id);

      return {
        user_id: params.user_id,
        notifications: notifications.map(n => ({
          id: n.id,
          campaign_id: n.campaign_id,
          sent_at: n.created_at,
          status: 'completed'
        })),
        total_count: notifications.length,
        unread_count: 0
      };
    } catch (error) {
      logger.error('Error fetching user notifications', error as Error, 'GET_USER_NOTIF');
      set.status = 500;
      return { error: 'Internal server error', status_code: 500 };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * GET /notification/group/:segment_id
   * Get notifications for user segment
   */
  app.get(GetPushNotifications.getSegmentNotifications, async ({ set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Get segment notifications request', 'GET_SEGMENT_NOTIF', { segmentId: params.segment_id });

      AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const notifications = await viewModel.getSegmentNotifications(params.segment_id);

      return {
        segment: params.segment_id,
        campaigns: notifications.map(n => ({
          id: n.id,
          title_identifier: n.title_identifier,
          sent_at: n.created_at,
          total_recipients: 0,
          successfully_sent: 0,
          failed_count: 0,
          success_rate: 100
        })),
        total_campaigns: notifications.length,
        aggregate_stats: {
          total_sent: 0,
          total_failed: 0,
          overall_success_rate: 100
        }
      };
    } catch (error) {
      logger.error('Error fetching segment notifications', error as Error, 'GET_SEGMENT_NOTIF');
      set.status = 500;
      return { error: 'Internal server error', status_code: 500 };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * GET /notification/topic/:topic_id
   * Get notifications for topic
   */
  app.get(GetPushNotifications.getTopicNotifications, async ({ set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Get topic notifications request', 'GET_TOPIC_NOTIF', { topicId: params.topic_id });

      AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const notifications = await viewModel.getTopicNotifications(params.topic_id);

      return {
        topic: params.topic_id,
        total_subscribers: 0,
        campaigns: notifications.map(n => ({
          id: n.id,
          title_identifier: n.title_identifier,
          sent_at: n.created_at,
          total_sent: 0,
          successfully_sent: 0,
          failed_count: 0,
          success_rate: 100
        })),
        total_campaigns: notifications.length,
        engagement_metrics: {
          avg_open_rate: 0,
          avg_click_rate: 0
        }
      };
    } catch (error) {
      logger.error('Error fetching topic notifications', error as Error, 'GET_TOPIC_NOTIF');
      set.status = 500;
      return { error: 'Internal server error', status_code: 500 };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * PUT /notification/:notification_id
   * Update notification status
   */
  app.put(PutPushNotifications.updateNotificationById, async ({ body, set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Update notification status request', 'UPDATE_NOTIF', { notificationId: params.notification_id });

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      const result = await viewModel.updateNotificationStatus(
        params.notification_id,
        (body as any).status,
        auth
      );

      set.status = 200;
      logger.logOperation('UPDATE_STATUS', params.notification_id, 'SUCCESS', 0, {
        newStatus: (body as any).status
      });

      return {
        id: params.notification_id,
        status: result.status,
        updated_at: new Date().toISOString(),
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      logger.error('Error updating notification status', error as Error, 'UPDATE_NOTIF');
      set.status = 500;
      return {
        error: 'Error',
        message: (error as Error).message || 'Internal server error',
        status_code: 500,
        timestamp: new Date().toISOString()
      };
    } finally {
      logger.endRequest();
    }
  });

  /**
   * DELETE /notification/:notification_id
   * Delete notification
   */
  app.delete(DeletePushNotifications.deleteNotificationById, async ({ set, request, params }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Delete notification request', 'DELETE_NOTIF', { notificationId: params.notification_id });

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: request.headers }
      } as any);

      await viewModel.deleteNotification(params.notification_id, auth);

      set.status = 200;
      logger.logOperation('DELETE_NOTIF', params.notification_id, 'SUCCESS', 0, {});

      return {
        id: params.notification_id,
        deleted_at: new Date().toISOString(),
        deleted_by: auth.admin_username,
        status: 'deleted',
        message: 'Campaign and associated data permanently removed'
      };
    } catch (error) {
      logger.error('Error deleting notification', error as Error, 'DELETE_NOTIF');
      set.status = 500;
      return {
        error: 'Error',
        message: (error as Error).message || 'Internal server error',
        status_code: 500,
        timestamp: new Date().toISOString()
      };
    } finally {
      logger.endRequest();
    }
  });

  logger.info('Notification routes setup complete', 'ROUTES_SETUP', {
    totalEndpoints: 12
  });

  return app;
}
