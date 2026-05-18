/**
 * SECURITY GUIDELINES & CVE PREVENTION
 * 
 * This document outlines security practices and CVE prevention strategies
 * for the Push Notification Service. Follow these guidelines to maintain
 * the security and integrity of the application.
 * 
 * Last Updated: 2024-05-18
 * Version: 1.0.0
 */

# Security Guidelines

## Authentication & Authorization

### JWT Implementation (Production)
```typescript
/**
 * Replace development tokens with JWT-based authentication.
 * Use RS256 (asymmetric) algorithm for token signing.
 */

// Development: Bearer token validation
// Production: JWT with RS256 signature verification
const token = authHeader.substring(7);
const decoded = jwt.verify(token, PUBLIC_KEY, { algorithm: 'RS256' });
```

### Token Management
- [YES] Tokens must expire (max 24 hours)
- [YES] Implement refresh token rotation
- [YES] Revoke tokens on logout
- [YES] Use HttpOnly cookies for token storage
- [YES] Implement rate limiting on auth endpoints
- [NO] Never store tokens in localStorage
- [NO] Never hardcode tokens in code

### Role-Based Access Control (RBAC)
```typescript
/**
 * RBAC ensures users can only perform authorized actions.
 * Implement principle of least privilege.
 */

const requiredRoles = [AdminRole.SUPER_ADMIN];
if (!requiredRoles.includes(auth.admin_role)) {
  throw new ForbiddenException('Insufficient permissions');
}
```

## Input Validation & Sanitization

### Validation Strategy
```typescript
/**
 * Always validate at system boundaries (user input, external APIs).
 * Never trust user input. Validate:
 * - Type (string, number, boolean, etc.)
 * - Length (min, max boundaries)
 * - Format (email, URL, UUID, ISO date)
 * - Enum values (whitelist allowed values)
 */

// [YES] CORRECT: Validate before use
ValidationMiddleware.validatePushRequest(body);
Validators.validateStringLength(title, 'title', 0, 120);

// [NO] WRONG: Using unvalidated input
const query = `WHERE title = '${userInput}'`; // SQL injection!
```

### Sanitization
- [YES] Trim whitespace: `value.trim()`
- [YES] Escape special characters for context
- [YES] Use parameterized queries (Prisma handles this)
- [YES] Validate URL schemes (https only in production)
- [NO] Never concatenate SQL strings
- [NO] Never use `eval()` or `Function()`

### Field Constraints
| Field | Min | Max | Pattern |
|-------|-----|-----|---------|
| title_identifier | 1 | 100 | alphanumeric + underscore |
| title | 1 | 120 | any unicode |
| body | 1 | 5000 | any unicode |
| icon_url | N/A | 512 | valid HTTPS URL |
| action_url | N/A | 512 | valid HTTPS URL |
| user_id | 1 | 64 | alphanumeric |

## Common Vulnerabilities (CVE)

### 1. SQL Injection (CWE-89)
**Risk:** Database compromise, data theft

**Prevention:**
- [YES] Use Prisma (parameterized queries automatically)
- [YES] Never concatenate user input into SQL
- [YES] Use prepared statements

**Example:**
```typescript
// [YES] SAFE: Prisma handles parameterization
const log = await db.notificationLog.findMany({
  where: { userId: userInput }
});

// [NO] UNSAFE: SQL injection vulnerability
const result = await db.$queryRaw(`SELECT * WHERE user_id = '${userInput}'`);
```

### 2. Authentication Bypass (CWE-287)
**Risk:** Unauthorized access to protected resources

**Prevention:**
- [YES] Validate every endpoint
- [YES] Verify token signature
- [YES] Check token expiration
- [YES] Verify role permissions
- [YES] Implement session timeout
- [NO] Don't skip auth on certain endpoints

**Example:**
```typescript
// [YES] CORRECT: All endpoints require auth
app.post('/push/all', async (context) => {
  const auth = AuthMiddleware.extractAuthContext(context);
  AuthMiddleware.requireRole([AdminRole.SUPER_ADMIN])(auth);
  // ... proceed
});

// [NO] WRONG: Exposed endpoint
app.post('/push/all', async (context) => {
  // No auth check!
});
```

### 3. Sensitive Data Exposure (CWE-200)
**Risk:** Passwords, keys, PII exposed in logs or errors

**Prevention:**
- [YES] Never log passwords, tokens, keys
- [YES] Sanitize error messages
- [YES] Use HTTPS in production
- [YES] Mask PII in logs
- [YES] Encrypt sensitive data at rest
- [NO] Don't expose stack traces to users
- [NO] Don't log full request/response bodies

**Example:**
```typescript
// [YES] CORRECT: Sanitized error
logger.error('Auth failed', { userId, attempt: 1 });

// [NO] WRONG: Exposed sensitive data
logger.error('Auth failed', { token, password });
```

### 4. Cross-Site Request Forgery (CSRF) (CWE-352)
**Risk:** Unauthorized actions from other domains

**Prevention:**
- [YES] Implement CSRF tokens for state-changing operations
- [YES] Use SameSite cookie attribute
- [YES] Validate Origin/Referer headers
- [YES] Use POST for state changes (not GET)

### 5. Rate Limiting (CWE-770)
**Risk:** Brute force attacks, DoS attacks

**Prevention:**
```typescript
/**
 * Implement rate limiting on sensitive endpoints.
 * Per-user and per-IP limits recommended.
 */

const RATE_LIMIT = {
  auth: '5 attempts per 15 minutes',
  push: '100 requests per minute per admin',
  default: '1000 requests per hour'
};
```

### 6. Insecure Deserialization (CWE-502)
**Risk:** Arbitrary code execution

**Prevention:**
- [YES] Use JSON (safe)
- [YES] Validate all input
- [YES] Use `JSON.parse()` safely
- [NO] Never use `eval()` or `Function()`

### 7. Broken Access Control (CWE-284)
**Risk:** Users accessing resources they shouldn't

**Prevention:**
- [YES] Check permissions on EVERY operation
- [YES] Validate ownership of resources
- [YES] Use resource-level permissions
- [YES] Log authorization failures

**Example:**
```typescript
// [YES] CORRECT: Check before returning data
async getUserNotifications(userId: string, auth: AuthContext) {
  // Only return current user's notifications
  if (auth.user_id !== userId && auth.role !== 'ADMIN') {
    throw new ForbiddenException();
  }
  return db.notificationLog.findMany({
    where: { userId }
  });
}
```

### 8. Injection Attacks (CWE-94)
**Risk:** Code execution through user input

**Prevention:**
- [YES] Use templating engines safely
- [YES] Escape output for context
- [YES] Never use `eval()` or `Function()`
- [YES] Use HTML/URL encoding for output

## Data Protection

### Encryption
```env
/**
 * At-Rest Encryption:
 * - Database: Use PostgreSQL encryption (pgcrypto)
 * - Sensitive fields: Encrypt before storage
 * 
 * In-Transit Encryption:
 * - Always use HTTPS
 * - TLS 1.2 minimum
 */

ENABLE_DATABASE_ENCRYPTION=true
ENCRYPTION_KEY=your-256-bit-hex-key
TLS_MIN_VERSION=1.2
```

### Password Hashing
```typescript
/**
 * If storing passwords, use bcrypt with proper cost factor.
 * Cost: 12+ (2^12 iterations)
 */

import bcrypt from 'bcrypt';

const hashedPassword = await bcrypt.hash(password, 12);
const isValid = await bcrypt.compare(inputPassword, hashedPassword);
```

### Token Storage
- [YES] Access token: In-memory or HttpOnly cookie
- [YES] Refresh token: HttpOnly, Secure, SameSite=Strict
- [YES] Private keys: Never in code, use environment variables
- [NO] Never store in localStorage or sessionStorage

## Deployment Security

### Environment Variables
```bash
/**
 * Never commit secrets to version control.
 * Use environment variables for all secrets.
 */

# .env (development only, never commit)
DATABASE_URL=postgresql://...
JWT_SECRET=your-secret-key
VAPID_PRIVATE_KEY=...

# Production: Use secrets management service
# AWS Secrets Manager, HashiCorp Vault, Azure Key Vault, etc.
```

### Headers
```typescript
/**
 * Security headers prevent various attacks.
 * Configure your reverse proxy or middleware.
 */

app.use((context) => {
  context.set('X-Content-Type-Options', 'nosniff');
  context.set('X-Frame-Options', 'DENY');
  context.set('X-XSS-Protection', '1; mode=block');
  context.set('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
  context.set('Content-Security-Policy', "default-src 'self'");
});
```

### HTTPS Only
```env
# Production
NODE_ENV=production
FORCE_HTTPS=true
TLS_CERT=/path/to/cert.pem
TLS_KEY=/path/to/key.pem
```

## Logging & Monitoring

### What to Log
```typescript
/**
 * Log security-relevant events for audit trails.
 * Implementation: Use structured logging with correlationId.
 */

// ✅ Log these events
logger.info('Admin created campaign', { 
  adminId, 
  campaignId, 
  timestamp 
});

logger.warn('Failed auth attempt', { 
  username, 
  ipAddress, 
  timestamp 
});

logger.error('Suspicious activity detected', { 
  userId, 
  action, 
  details 
});

// ❌ Don't log these
// Passwords, tokens, keys, full request bodies
```

### Monitoring
```typescript
/**
 * Monitor for security anomalies:
 * - High number of failed auth attempts
 * - Unusual API call patterns
 * - Rate limit violations
 * - Database query errors
 */

const ALERTS = {
  failed_auth_threshold: 5,        // Failed attempts
  request_rate_threshold: 1000,    // Per minute
  large_response_threshold: 1048576 // 1MB
};
```

### Audit Logging
```sql
/**
 * Log all state-changing operations for compliance.
 * Implement immutable audit logs.
 */

CREATE TABLE audit_logs (
  id UUID PRIMARY KEY,
  admin_id INT NOT NULL,
  action VARCHAR(100) NOT NULL,
  resource_id UUID,
  resource_type VARCHAR(50),
  changes JSONB,
  ip_address INET,
  timestamp TIMESTAMP DEFAULT NOW()
);

-- Create index for queries
CREATE INDEX idx_audit_logs_admin_timestamp 
ON audit_logs(admin_id, timestamp);
```

## API Security

### CORS Configuration
```typescript
/**
 * Restrict cross-origin requests to trusted domains.
 * Never use * in production.
 */

const ALLOWED_ORIGINS = [
  'https://app.example.com',
  'https://admin.example.com'
];

app.use((context) => {
  const origin = context.request.headers.get('origin');
  if (ALLOWED_ORIGINS.includes(origin)) {
    context.set('Access-Control-Allow-Origin', origin);
  }
});
```

### Request Size Limits
```typescript
/**
 * Prevent large payload attacks.
 */

app.use(limitJson({ limit: '1mb' }));
app.use(limitForm({ limit: '1mb' }));
```

## Dependency Security

### Package Management
```bash
/**
 * Keep dependencies updated and secure.
 */

# Audit dependencies
npm audit
npm audit fix

# Update packages
npm update
npm install --save-exact @prisma/client

# Lock dependencies
npm ci  # Use package-lock.json

# Check for vulnerabilities
npm audit --production
```

### Dependency Scanning
```bash
# Use tools like:
# - npm audit
# - OWASP Dependency-Check
# - Snyk
# - GitHub Dependabot

# In CI/CD pipeline:
npm audit --audit-level=moderate || exit 1
```

## Security Checklist

### Pre-Deployment
- [ ] All user input validated
- [ ] No hardcoded secrets in code
- [ ] HTTPS enabled
- [ ] CORS properly configured
- [ ] Rate limiting implemented
- [ ] Security headers set
- [ ] Authentication required on all endpoints
- [ ] Authorization checks in place
- [ ] Error messages sanitized
- [ ] Logging configured (no secrets)
- [ ] Database encrypted
- [ ] Backups tested
- [ ] Incident response plan ready
- [ ] Security audit completed
- [ ] Dependencies up-to-date
- [ ] No known vulnerabilities

### After Deployment
- [ ] Monitor error logs
- [ ] Monitor auth failures
- [ ] Monitor API usage patterns
- [ ] Check for anomalies
- [ ] Review audit logs
- [ ] Update security patches
- [ ] Regular backups
- [ ] Security training

## Incident Response

### If Breach Detected
1. **STOP**: Isolate affected systems
2. **ASSESS**: Determine scope and impact
3. **NOTIFY**: Alert stakeholders
4. **SECURE**: Apply patches/fixes
5. **VERIFY**: Confirm issue resolved
6. **REVIEW**: Post-mortem analysis
7. **IMPROVE**: Implement preventions

### Contact Information
```
Security Team: security@example.com
Incident Hotline: +1-XXX-XXX-XXXX
On-Call Rotation: [link to runbook]
```

## References

- [OWASP Top 10 2023](https://owasp.org/www-project-top-ten/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework/)
- [CWE Top 25](https://cwe.mitre.org/top25/)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/sql-syntax.html#SQL-SYNTAX-LEXICAL)
- [Node.js Security Best Practices](https://nodejs.org/en/docs/guides/security/)

## Compliance

This application implements controls for:
- [YES] OWASP Top 10
- [YES] GDPR (data protection)
- [YES] SOC 2 (security, availability, processing integrity)
- [YES] PCI DSS (payment card security, if applicable)
- [YES] HIPAA (health data, if applicable)

## Questions?

For security questions or to report vulnerabilities:
- Email: security@example.com
- Do not open public issues for security vulnerabilities
- Include steps to reproduce and potential impact

---

**Remember:** Security is a process, not a product. Continuously improve and stay vigilant.
