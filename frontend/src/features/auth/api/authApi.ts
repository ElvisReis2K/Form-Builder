import { apiRequest, apiURL } from '../../../lib/api';
import type { LoginInput, RegisterInput } from '../types';

export function login(input: LoginInput) {
  return apiRequest('/api/auth/login', 'post', {
    body: input,
  });
}

export function register(input: RegisterInput) {
  return apiRequest('/api/auth/register', 'post', {
    body: input,
  });
}

export function getCurrentUser() {
  return apiRequest('/api/auth/me', 'get');
}

export function logout() {
  return apiRequest('/api/auth/logout', 'post');
}

export function googleLoginURL() {
  return apiURL('/api/auth/google');
}
