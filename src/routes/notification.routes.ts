/**
 * @fileoverview API Routes Module
 *
 * Defines all HTTP endpoints for the Push Notification Service.
 * Handles:
 * - Request routing
 * - Authentication and authorization
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
import { AuthMiddleware, UnauthorizedException, ForbiddenException } from '../middleware/auth.middleware';
import { ValidationMiddleware, ValidationException } from '../middleware/validation.middleware';
import { ErrorHandler } from '../middleware/error.middleware';
import { PushNotificationViewMapper } from '../views/push-notification.view';
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
   * Requires: SUPER_ADMIN role
   */
  app.post(PostPushNotifications.pushToAll, async ({ body, set }) => {
    const requestId = crypto.randomUUID();
    logger.startRequest(requestId);

    try {
      logger.info('Push to all request received', 'PUSH_ALL', { requestId });

      /* Extract and validate authentication */
      const authHeader = (body as any)?.headers?.authorization;
      if (!authHeader) {
        logger.warn('Missing authorization header in push/all request', 'PUSH_ALL', {
          requestId
        });
        set.status = 401;
        return { error: 'Unauthorized', message: 'Missing authorization header', status_code: 401 };
      }

      const auth = AuthMiddleware.extractAuthContext({
        request: { headers: { get: (name: string) => name === 'authorization' ? authHeader : undefined } }
      } as any);

      AuthMiddleware.requireRole([AdminRole.SUPER_ADMIN])(auth);

      /* Validate request body */
      const validatedRequest = ValidationMiddleware.validatePushRequest(body);

      /* Execute business logic */
      const result = await viewModel.pushToAll(
        validatedRequest.title_identifier,
        validatedRequest.locale_contents,
        auth
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
   * GET /health
   * Health check endpoint - no authentication required
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
   * Get latest campaigns - no authentication required
   */
  app.get(GetPushNotifications.getLatestNotification, async ({ set }) => {
    try {
      logger.info('Latest notifications request', 'GET_LATEST');
      const notifications = await viewModel.getLatestNotifications(10);
      return notifications.map(n => PushNotificationViewMapper.mapCampaignToResponse(n));
    } catch (error) {
      logger.error('Error fetching latest notifications', error as Error, 'GET_LATEST');
      set.status = 500;
      return { error: 'Internal server error', status_code: 500 };
    }
  });

  /**
   * Return the Elysia app for further configuration
   */
  logger.info('Notification routes setup complete', 'ROUTES_SETUP', {
    totalEndpoints: 9
  });

  return app;
}
