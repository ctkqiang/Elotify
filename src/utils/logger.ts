/**
 * @fileoverview Comprehensive Logging Module
 *
 * Provides verbose, structured logging for debugging and monitoring.
 * All operations are logged with full context for easy issue identification.
 *
 * Features:
 * - Multi-level logging (DEBUG, INFO, WARN, ERROR)
 * - JSON structured logging for production
 * - Text format for development
 * - Request tracking and correlation IDs
 * - Performance metrics logging
 * - Error stack traces with full context
 * - Request/response logging
 * - Database query logging
 * - Authentication event logging
 * - Business operation logging
 * - Permission tracking
 *
 * Usage:
 *   import { logger } from '@/utils/logger';
 *   logger.info('Operation started', 'MODULE_NAME', { userId: 123 });
 *   logger.error('Operation failed', error, 'MODULE_NAME', { userId: 123 });
 *
 * Environment Variables:
 *   LOG_LEVEL: DEBUG|INFO|WARN|ERROR (default: INFO)
 *   LOG_FORMAT: json|text (default: text)
 *
 * @module src/utils/logger
 */

export enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARN = 'WARN',
  ERROR = 'ERROR'
}

/**
 * Log entry structure
 */
interface LogEntry {
  timestamp: string;
  level: LogLevel;
  module?: string;
  message: string;
  context?: Record<string, any>;
  error?: {
    name: string;
    message: string;
    stack?: string;
  };
  duration?: number;
  requestId?: string;
  userId?: string;
}

/**
 * Request context for logging
 */
interface RequestLogContext {
  method: string;
  url: string;
  statusCode?: number;
  duration?: number;
  requestId: string;
  userId?: string;
  ip?: string;
  userAgent?: string;
}

/**
 * Comprehensive logger with verbose output for debugging
 * Logs all operations with full context for issue identification
 */
class Logger {
  private logLevel: LogLevel = LogLevel.INFO;
  private logFormat: string = 'text';
  private requestIdStack: string[] = [];

  constructor() {
    const level = process.env.LOG_LEVEL?.toUpperCase();
    if (level && Object.values(LogLevel).includes(level as LogLevel)) {
      this.logLevel = level as LogLevel;
    }

    this.logFormat = process.env.LOG_FORMAT || 'text';
  }

  /**
   * Check if log level should be output based on configured level
   * Log levels in order: DEBUG < INFO < WARN < ERROR
   */
  private shouldLog(level: LogLevel): boolean {
    const levels = [LogLevel.DEBUG, LogLevel.INFO, LogLevel.WARN, LogLevel.ERROR];
    return levels.indexOf(level) >= levels.indexOf(this.logLevel);
  }

  /**
   * Format log entry for output
   * Supports both JSON format (production) and text format (development)
   * Includes comprehensive context information for debugging
   */
  private format(entry: LogEntry): string {
    /* JSON format for structured logging and log aggregation */
    if (this.logFormat === 'json') {
      return JSON.stringify({
        timestamp: entry.timestamp,
        level: entry.level,
        module: entry.module,
        message: entry.message,
        context: entry.context,
        duration: entry.duration,
        requestId: entry.requestId,
        userId: entry.userId,
        error: entry.error
      });
    }

    /**
     * Text format with comprehensive information for development
     * Format includes all available context for easy debugging
     */
    const parts: string[] = [];

    /* Timestamp and level */
    parts.push(`[${entry.timestamp}]`);
    parts.push(`[${entry.level}]`);

    /* Module context */
    if (entry.module) {
      parts.push(`[${entry.module}]`);
    }

    /* Main message */
    parts.push(entry.message);

    /* Request ID for tracing */
    if (entry.requestId) {
      parts.push(`| RequestID: ${entry.requestId}`);
    }

    /* User ID for tracking user actions */
    if (entry.userId) {
      parts.push(`| UserID: ${entry.userId}`);
    }

    /* Additional context data */
    if (entry.context && Object.keys(entry.context).length > 0) {
      const contextStr = JSON.stringify(entry.context, null, 2);
      parts.push(`| Context:\n${contextStr}`);
    }

    /* Operation duration */
    if (entry.duration !== undefined) {
      parts.push(`| Duration: ${entry.duration}ms`);
    }

    /* Error information with stack trace */
    if (entry.error) {
      parts.push(`| Error Name: ${entry.error.name}`);
      parts.push(`| Error Message: ${entry.error.message}`);
      if (entry.error.stack) {
        parts.push(`| Stack Trace:\n${entry.error.stack}`);
      }
    }

    return parts.join(' ');
  }

  /**
   * Push request ID to stack for request tracking
   * Allows nested request tracking
   */
  private pushRequestId(requestId: string): void {
    this.requestIdStack.push(requestId);
  }

  /**
   * Pop request ID from stack
   * Called when request completes
   */
  private popRequestId(): string | undefined {
    return this.requestIdStack.pop();
  }

  /**
   * Get current request ID from stack
   * Used to include requestId in all logs for current request
   */
  private getCurrentRequestId(): string | undefined {
    return this.requestIdStack.length > 0
      ? this.requestIdStack[this.requestIdStack.length - 1]
      : undefined;
  }

  /**
   * Core logging method with full context
   * Internal method used by all public logging methods
   */
  private log(
    level: LogLevel,
    module: string | undefined,
    message: string,
    context?: any,
    error?: Error,
    duration?: number
  ): void {
    if (!this.shouldLog(level)) return;

    const entry: LogEntry = {
      timestamp: new Date().toISOString(),
      level,
      module,
      message,
      context,
      duration,
      requestId: this.getCurrentRequestId(),
      error: error
        ? {
            name: error.name,
            message: error.message,
            stack: error.stack
          }
        : undefined
    };

    const formatted = this.format(entry);
    const consoleMethod = level === LogLevel.ERROR ? console.error : console.log;
    consoleMethod(formatted);
  }

  /**
   * Log debug message (development only, not shown in production)
   * Use for detailed debugging information
   *
   * @param message - Debug message describing the operation
   * @param module - Module name for context (e.g., 'AUTH', 'DATABASE')
   * @param context - Additional context data as key-value object
   */
  debug(message: string, module?: string, context?: any): void {
    this.log(LogLevel.DEBUG, module, message, context);
  }

  /**
   * Log info message for normal operations
   * Use for successful operations and important milestones
   *
   * @param message - Information message
   * @param module - Module name for context
   * @param context - Additional context data
   */
  info(message: string, module?: string, context?: any): void {
    this.log(LogLevel.INFO, module, message, context);
  }

  /**
   * Log warning message for potentially problematic situations
   * Use for unexpected behavior that doesn't stop execution
   *
   * @param message - Warning message
   * @param module - Module name for context
   * @param context - Additional context data
   */
  warn(message: string, module?: string, context?: any): void {
    this.log(LogLevel.WARN, module, message, context);
  }

  /**
   * Log error with full stack trace and context
   * Use for failures and exceptions
   *
   * @param message - Error message
   * @param error - Error object or error message string
   * @param module - Module name for context
   * @param context - Additional context data
   */
  error(message: string, error?: Error | string, module?: string, context?: any): void {
    const err = typeof error === 'string' ? new Error(error) : error;
    this.log(LogLevel.ERROR, module, message, context, err);
  }

  /**
   * Log HTTP request with comprehensive details
   * Called for each API request to track traffic and performance
   *
   * @param context - Request context containing method, URL, status, duration
   */
  logRequest(context: RequestLogContext): void {
    const message = `HTTP Request: ${context.method} ${context.url}`;
    const logContext = {
      statusCode: context.statusCode,
      duration: context.duration,
      ip: context.ip,
      userAgent: context.userAgent ? context.userAgent.substring(0, 100) : undefined
    };

    /* Log failed requests as warnings */
    if (context.statusCode && context.statusCode >= 400) {
      this.warn(message, 'HTTP_REQUEST', logContext);
    } else {
      this.info(message, 'HTTP_REQUEST', logContext);
    }
  }

  /**
   * Log database operation with query details
   * Called for all database operations to track performance
   *
   * @param operation - Operation type (SELECT, INSERT, UPDATE, DELETE)
   * @param table - Table name being queried
   * @param duration - Query execution duration in milliseconds
   * @param rowsAffected - Number of rows affected by operation
   */
  logDatabase(
    operation: string,
    table: string,
    duration: number,
    rowsAffected?: number
  ): void {
    const message = `Database Operation: ${operation} on table '${table}'`;
    const context = {
      operation,
      table,
      durationMs: duration,
      rowsAffected: rowsAffected || 0
    };
    this.debug(message, 'DATABASE', context);
  }

  /**
   * Log authentication events (attempts, success, failure)
   * Called for all authentication-related operations
   *
   * @param event - Event type (AUTH_ATTEMPT, AUTH_SUCCESS, AUTH_FAILURE)
   * @param username - Username attempting authentication
   * @param success - Whether authentication was successful
   * @param ip - Client IP address
   * @param reason - Failure reason if authentication failed
   */
  logAuth(
    event: 'AUTH_ATTEMPT' | 'AUTH_SUCCESS' | 'AUTH_FAILURE',
    username: string,
    success: boolean,
    ip?: string,
    reason?: string
  ): void {
    const message = `Authentication Event: ${event} for user '${username}'`;
    const context = {
      event,
      username,
      success,
      ip: ip || 'unknown',
      reason: reason || 'N/A',
      timestamp: new Date().toISOString()
    };

    /* Log authentication failures as warnings */
    if (event === 'AUTH_FAILURE' || !success) {
      this.warn(message, 'AUTHENTICATION', context);
    } else {
      this.info(message, 'AUTHENTICATION', context);
    }
  }

  /**
   * Log business operation (campaign creation, push send, etc.)
   * Called for all significant business operations
   *
   * @param operation - Operation name (e.g., 'CREATE_CAMPAIGN', 'SEND_NOTIFICATION')
   * @param resourceId - ID of resource being operated on
   * @param status - Operation status (SUCCESS or FAILURE)
   * @param duration - Operation execution duration in milliseconds
   * @param details - Additional operation details and results
   */
  logOperation(
    operation: string,
    resourceId: string,
    status: 'SUCCESS' | 'FAILURE',
    duration: number,
    details?: Record<string, any>
  ): void {
    const message = `Business Operation: ${operation} on resource '${resourceId}' - Status: ${status}`;
    const context = {
      operation,
      resourceId,
      status,
      durationMs: duration,
      ...details
    };

    /* Log failed operations as errors */
    if (status === 'FAILURE') {
      this.error(message, undefined, 'BUSINESS_OPERATION', context);
    } else {
      this.info(message, 'BUSINESS_OPERATION', context);
    }
  }

  /**
   * Log validation failure with detailed information
   * Called when input validation fails
   *
   * @param field - Field name that failed validation
   * @param value - Value that failed validation
   * @param rule - Validation rule that was violated
   */
  logValidation(field: string, value: any, rule: string): void {
    const message = `Validation Failed: field '${field}' violates rule '${rule}'`;
    const valueStr = typeof value === 'object' ? JSON.stringify(value) : String(value);
    const context = {
      field,
      value: valueStr,
      rule,
      timestamp: new Date().toISOString()
    };
    this.warn(message, 'VALIDATION', context);
  }

  /**
   * Log permission denied event
   * Called when user attempts unauthorized action
   *
   * @param userId - User ID that was denied access
   * @param action - Action that was attempted
   * @param resource - Resource being accessed
   */
  logPermissionDenied(userId: number | string, action: string, resource: string): void {
    const message = `Permission Denied: user '${userId}' attempted unauthorized action '${action}' on resource '${resource}'`;
    const context = {
      userId,
      action,
      resource,
      timestamp: new Date().toISOString()
    };
    this.warn(message, 'AUTHORIZATION', context);
  }

  /**
   * Start request tracking
   * Called at the beginning of request processing
   *
   * @param requestId - Unique request identifier
   */
  startRequest(requestId: string): void {
    this.pushRequestId(requestId);
    this.debug(`Request Started: ${requestId}`, 'REQUEST_LIFECYCLE', { requestId });
  }

  /**
   * End request tracking
   * Called at the end of request processing
   */
  endRequest(): void {
    const requestId = this.popRequestId();
    if (requestId) {
      this.debug(`Request Completed: ${requestId}`, 'REQUEST_LIFECYCLE', { requestId });
    }
  }

  /**
   * Log campaign event with details
   * Called for campaign creation, sending, completion
   */
  logCampaign(
    event: 'CREATED' | 'QUEUED' | 'SENDING' | 'COMPLETED' | 'FAILED',
    campaignId: string,
    details?: Record<string, any>
  ): void {
    const message = `Campaign Event: ${event} for campaign '${campaignId}'`;
    const context = {
      event,
      campaignId,
      ...details
    };
    this.info(message, 'CAMPAIGN', context);
  }

  /**
   * Log notification delivery attempt
   * Called for each delivery attempt
   */
  logDelivery(
    campaignId: string,
    subscriptionId: string,
    status: 'QUEUED' | 'SENT' | 'FAILED_EXPIRED' | 'FAILED_ERROR',
    httpStatusCode?: number,
    errorMessage?: string
  ): void {
    const message = `Notification Delivery: status '${status}' for subscription '${subscriptionId}'`;
    const context = {
      campaignId,
      subscriptionId,
      status,
      httpStatusCode,
      errorMessage
    };
    this.debug(message, 'DELIVERY', context);
  }
}

export const logger = new Logger();
