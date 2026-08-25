import { apiRequest } from '../../../lib/api';
import type { AuthResponse, LoginInput, RegisterInput } from '../types';

export function login(input: LoginInput) {
  return apiRequest<AuthResponse>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export function register(input: RegisterInput) {
  return apiRequest<AuthResponse>('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export function getCurrentUser() {
  return apiRequest<AuthResponse>('/api/auth/me');
}

export function logout() {
  return apiRequest<void>('/api/auth/logout', {
    method: 'POST',
  });
}
