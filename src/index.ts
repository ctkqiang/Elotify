/**
 * @fileoverview Application Entry Point
 *
 * Initializes the Push Notification Service with:
 * - Environment configuration loading
 * - Database connection establishment
 * - Route setup
 * - Graceful shutdown handling
 * - Comprehensive logging of all startup steps
 *
 * @module src/index
 */

import 'dotenv/config';
import { Elysia } from 'elysia';
import { setupNotificationRoutes } from './routes/notification.routes';
import { connectDatabase, disconnectDatabase, db } from './config/database';
import { logger } from './utils/logger';

/**
 * Bootstrap the application
 * - Load configuration
 * - Verify database connectivity
 * - Start HTTP server
 * - Setup graceful shutdown
 */
async function bootstrap() {
  const startTime = Date.now();

  try {
    /**
     * Log application startup with environment information
     */
    logger.info('Push Notification Service starting', 'APP_STARTUP', {
      environment: process.env.NODE_ENV || 'development',
      nodeVersion: process.version,
      port: process.env.PORT || 3000
    });

    /**
     * Initialize database connection
     * This will log connection attempts and results
     */
    logger.debug('Initializing database connection', 'APP_STARTUP');
    await connectDatabase();

    /**
     * Verify database health with test query
     * Ensures database is fully operational
     */
    logger.debug('Verifying database health', 'APP_STARTUP');
    const testStart = Date.now();
    await db.$queryRaw`SELECT 1`;
    const testDuration = Date.now() - testStart;
    logger.info('Database health check passed', 'APP_STARTUP', {
      healthCheckDurationMs: testDuration
    });

    /**
     * Initialize HTTP server
     */
    const port = parseInt(process.env.PORT || '3000', 10);
    logger.debug('Creating Elysia application instance', 'APP_STARTUP', { port });
    const app = new Elysia();

    /**
     * Setup all notification routes with middleware
     */
    logger.debug('Setting up API routes', 'APP_STARTUP');
    setupNotificationRoutes(app);

    /**
     * Root endpoint for health and status
     */
    app.get('/', () => {
      logger.debug('Root endpoint called', 'HTTP_REQUEST');
      return {
        message: 'Push Notification Service',
        version: '1.0.0',
        status: 'running',
        timestamp: new Date().toISOString()
      };
    });

    /**
     * Start HTTP server
     */
    logger.info(`🦊 Elysia is running at http://localhost:${port}`); 
    app.listen(port);

    const totalStartupTime = Date.now() - startTime;
    logger.info('Application started successfully', 'APP_STARTUP', {
      port,
      startupTimeMscompletion: totalStartupTime,
      timestamp: new Date().toISOString()
    });

    /**
     * Handle graceful shutdown on SIGTERM
     * Logged for audit trail and debugging
     */
    process.on('SIGTERM', async () => {
      logger.warn('SIGTERM signal received: starting graceful shutdown', 'SHUTDOWN');
      try {
        logger.info('Closing database connection', 'SHUTDOWN');
        await disconnectDatabase();
        logger.info('Graceful shutdown completed successfully', 'SHUTDOWN', {
          signal: 'SIGTERM'
        });
        process.exit(0);
      } catch (error) {
        logger.error('Error during graceful shutdown', error as Error, 'SHUTDOWN');
        process.exit(1);
      }
    });

    /**
     * Handle graceful shutdown on SIGINT (Ctrl+C)
     * Logged for audit trail and debugging
     */
    process.on('SIGINT', async () => {
      logger.warn('SIGINT signal received: starting graceful shutdown', 'SHUTDOWN');
      try {
        logger.info('Closing database connection', 'SHUTDOWN');
        await disconnectDatabase();
        logger.info('Graceful shutdown completed successfully', 'SHUTDOWN', {
          signal: 'SIGINT'
        });
        process.exit(0);
      } catch (error) {
        logger.error('Error during graceful shutdown', error as Error, 'SHUTDOWN');
        process.exit(1);
      }
    });

    /**
     * Handle uncaught exceptions
     * Log for debugging and monitoring
     */
    process.on('uncaughtException', (error) => {
      logger.error('Uncaught exception occurred', error as Error, 'UNCAUGHT_EXCEPTION');
      process.exit(1);
    });

    /**
     * Handle unhandled promise rejections
     * Log for debugging and monitoring
     */
    process.on('unhandledRejection', (reason, promise) => {
      logger.error(
        'Unhandled promise rejection',
        new Error(String(reason)),
        'UNHANDLED_REJECTION',
        {
          promise: String(promise)
        }
      );
    });
  } catch (error) {
    const totalTime = Date.now() - startTime;
    logger.error(
      'Failed to start application after ' + totalTime + 'ms',
      error as Error,
      'APP_STARTUP'
    );
    process.exit(1);
  }
}

/**
 * Start the application
 * Wrapped in IIFE to use async/await at module level
 */
bootstrap();

