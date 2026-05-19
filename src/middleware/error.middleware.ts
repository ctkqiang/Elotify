import { Context } from 'elysia';
import { PushNotificationViewMapper } from '../views/push-notification.view';
import { UnauthorizedException, ForbiddenException } from './auth.middleware';
import { ValidationException } from './validation.middleware';

export class ErrorHandler {
  static handle(error: any, context: any) {
    console.error('[Error]', {
      name: error.name,
      message: error.message,
      path: context.request.url,
      method: context.request.method,
      timestamp: new Date().toISOString()
    });

    let statusCode = 500;
    let message = 'Internal server error';

    if (error instanceof UnauthorizedException) {
      statusCode = 401;
      message = error.message;
    } else if (error instanceof ForbiddenException) {
      statusCode = 403;
      message = error.message;
    } else if (error instanceof ValidationException) {
      statusCode = 400;
      message = error.message;
    } else if (error instanceof NotFoundError) {
      statusCode = 404;
      message = error.message;
    } else if (error instanceof ConflictError) {
      statusCode = 409;
      message = error.message;
    }

    context.set.status = statusCode;
    return PushNotificationViewMapper.mapErrorResponse(
      new Error(message),
      statusCode
    );
  }
}

export class NotFoundError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'NotFoundError';
  }
}

export class ConflictError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ConflictError';
  }
}
