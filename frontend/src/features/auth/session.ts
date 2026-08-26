import type { QueryClient } from '@tanstack/react-query';

import { postApiAuthLogout } from '../../api/generated/client';
import type { AuthResponse } from './types';
import { authMeQueryKey } from './queryKeys';

const reauthenticationRequiredKey = 'form_builder_reauthentication_required';

export function requireReauthentication() {
  window.sessionStorage.setItem(reauthenticationRequiredKey, 'true');
}

export function isReauthenticationRequired() {
  return window.sessionStorage.getItem(reauthenticationRequiredKey) === 'true';
}

export function clearReauthenticationRequirement() {
  window.sessionStorage.removeItem(reauthenticationRequiredKey);
}

export function completeAuthentication(queryClient: QueryClient, authResponse: AuthResponse) {
  clearReauthenticationRequirement();
  queryClient.setQueryData(authMeQueryKey, authResponse);
}

export async function endAuthenticatedSession(queryClient: QueryClient) {
  requireReauthentication();

  try {
    await postApiAuthLogout();
  } finally {
    queryClient.removeQueries({ queryKey: authMeQueryKey });
  }
}
