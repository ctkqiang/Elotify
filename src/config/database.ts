/**
 * @fileoverview Database Configuration Module
 *
 * Manages Prisma ORM initialization and connection lifecycle.
 * - Singleton pattern for database client
 * - Development logging enabled
 * - Graceful error handling
 * - Supports hot-reload in development
 * - Comprehensive verbose logging
 *
 * @module src/config/database
 */

import { PrismaClient } from '@prisma/client';
import { logger } from '../utils/logger';

let prisma: PrismaClient;

declare global {
  /**
   * Global Prisma instance to prevent multiple connections in dev mode
   */
  var prisma: PrismaClient | undefined;
}

/**
 * Initialize Prisma Client with environment-specific configuration
 *
 * In development: Enables query logging and uses global singleton
 * In production: Strict connection management with optimized settings
 */
if (process.env.NODE_ENV === 'production') {
  prisma = new PrismaClient({
    log: [
      { emit: 'stdout', level: 'error' },
      { emit: 'stdout', level: 'warn' }
    ]
  });
} else {
  if (!global.prisma) {
    global.prisma = new PrismaClient({
      log: [
        { emit: 'stdout', level: 'info' },
        { emit: 'stdout', level: 'warn' },
        { emit: 'stdout', level: 'error' }
      ]
    });
  }
  prisma = global.prisma;
}

export const db = prisma;
export type Database = PrismaClient;

/**
 * Establish database connection with health check
 *
 * Verifies that the database is accessible and operational.
 * Logs all connection attempts and results for debugging.
 *
 * @async
 * @throws {Error} If database connection fails
 * @returns {Promise<void>}
 */
export async function connectDatabase(): Promise<void> {
  const startTime = Date.now();
  logger.info('Attempting database connection', 'DATABASE_CONNECTION', {
    environment: process.env.NODE_ENV,
    databaseUrl: process.env.DATABASE_URL?.split('@')[1] || 'unknown'
  });

  try {
    await db.$connect();
    const duration = Date.now() - startTime;
    logger.logDatabase('CONNECT', 'N/A', duration);
    logger.info('Database connection established successfully', 'DATABASE_CONNECTION', {
      durationMs: duration
    });
  } catch (error) {
    const duration = Date.now() - startTime;
    logger.error(
      'Database connection failed after ' + duration + 'ms',
      error as Error,
      'DATABASE_CONNECTION',
      {
        durationMs: duration,
        environment: process.env.NODE_ENV
      }
    );
    throw error;
  }
}

/**
 * Gracefully close database connection
 *
 * Ensures all pending queries complete before disconnecting.
 * Called during application shutdown.
 * Logs connection closure for audit trail.
 *
 * @async
 * @returns {Promise<void>}
 */
export async function disconnectDatabase(): Promise<void> {
  const startTime = Date.now();
  logger.info('Closing database connection', 'DATABASE_CONNECTION');

  try {
    await db.$disconnect();
    const duration = Date.now() - startTime;
    logger.info('Database connection closed gracefully', 'DATABASE_CONNECTION', {
      durationMs: duration
    });
  } catch (error) {
    logger.error(
      'Error closing database connection',
      error as Error,
      'DATABASE_CONNECTION'
    );
    throw error;
  }
}

export default db;
