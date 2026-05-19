import { ValidationException } from '../middleware/validation.middleware';

export class Validators {
  static isValidUUID(value: string): boolean {
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    return uuidRegex.test(value);
  }

  static isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  static isValidUrl(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }

  static isValidISODate(date: string): boolean {
    const parsed = new Date(date);
    return !isNaN(parsed.getTime()) && parsed.toISOString().startsWith(date.split('T')[0]);
  }

  static isBase64(str: string): boolean {
    try {
      return btoa(atob(str)) === str;
    } catch {
      return false;
    }
  }

  static sanitizeString(value: string, maxLength?: number): string {
    let sanitized = value.trim();
    if (maxLength && sanitized.length > maxLength) {
      sanitized = sanitized.substring(0, maxLength);
    }
    return sanitized;
  }

  static validateStringLength(
    value: string,
    fieldName: string,
    minLength?: number,
    maxLength?: number
  ): void {
    if (minLength && value.length < minLength) {
      throw new ValidationException(
        `${fieldName} must be at least ${minLength} characters long`
      );
    }
    if (maxLength && value.length > maxLength) {
      throw new ValidationException(
        `${fieldName} must not exceed ${maxLength} characters`
      );
    }
  }

  static validateRequired(value: any, fieldName: string): void {
    if (value === null || value === undefined || value === '') {
      throw new ValidationException(`${fieldName} is required`);
    }
  }

  static validateEnum(value: string, allowedValues: string[], fieldName: string): void {
    if (!allowedValues.includes(value)) {
      throw new ValidationException(
        `${fieldName} must be one of: ${allowedValues.join(', ')}`
      );
    }
  }

  static validateObject(obj: any, requiredFields: string[]): void {
    if (!obj || typeof obj !== 'object') {
      throw new ValidationException('Invalid object');
    }

    for (const field of requiredFields) {
      if (!(field in obj) || obj[field] === null || obj[field] === undefined) {
        throw new ValidationException(`Required field missing: ${field}`);
      }
    }
  }

  static validatePageSize(page?: number, limit?: number): { page: number; limit: number } {
    const validatedPage = Math.max(1, page || 1);
    const validatedLimit = Math.min(Math.max(1, limit || 50), 1000);
    return { page: validatedPage, limit: validatedLimit };
  }
}
