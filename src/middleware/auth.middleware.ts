import { Context } from 'elysia';
import { AuthContext, AdminRole } from '../types';

const VALID_TOKENS = new Map<string, AuthContext>([
  ['admin-token-super', { admin_id: 1, admin_username: 'admin', admin_role: AdminRole.SUPER_ADMIN }],
  ['admin-token-regular', { admin_id: 2, admin_username: 'user', admin_role: AdminRole.ADMIN }]
]);

export class AuthMiddleware {
  static extractAuthContext(context: Context): AuthContext {
    const authHeader = context.request.headers.get('Authorization');

    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      throw new UnauthorizedException('Missing or invalid Authorization header');
    }

    const token = authHeader.substring(7);
    const authContext = VALID_TOKENS.get(token);

    if (!authContext) {
      throw new UnauthorizedException('Invalid or expired token');
    }

    return authContext;
  }

  static requireRole(requiredRoles: AdminRole[]) {
    return (auth: AuthContext) => {
      if (!requiredRoles.includes(auth.admin_role)) {
        throw new ForbiddenException('Insufficient permissions');
      }
    };
  }

  static verifyToken(token: string): AuthContext {
    const authContext = VALID_TOKENS.get(token);
    if (!authContext) {
      throw new UnauthorizedException('Invalid token');
    }
    return authContext;
  }
}

export class UnauthorizedException extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'UnauthorizedException';
  }
}

export class ForbiddenException extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ForbiddenException';
  }
}
