import type { QueryClient } from '@tanstack/react-query';

import { postApiAuthLogout } from '../../api/generated/client';
import type { AuthResponse } from './types';
import { authMeQueryKey } from './queryKeys';

const pendingGoogleOAuthKey = 'form_builder_pending_google_oauth';

let authenticatedInCurrentPageLoad = false;

export function requireReauthentication() {
  authenticatedInCurrentPageLoad = false;
  window.sessionStorage.removeItem(pendingGoogleOAuthKey);
}

export function allowGoogleOAuthReturn() {
  window.sessionStorage.setItem(pendingGoogleOAuthKey, 'true');
}

export function canUseAuthenticatedSession() {
  return authenticatedInCurrentPageLoad;
}

export function consumeGoogleOAuthReturn() {
  const canReturnFromGoogle = window.sessionStorage.getItem(pendingGoogleOAuthKey) === 'true';
  window.sessionStorage.removeItem(pendingGoogleOAuthKey);

  if (canReturnFromGoogle) {
    authenticatedInCurrentPageLoad = true;
  }

  return canReturnFromGoogle;
}

export function confirmAuthenticatedSession() {
  authenticatedInCurrentPageLoad = true;
}

export function completeAuthentication(queryClient: QueryClient, authResponse: AuthResponse) {
  confirmAuthenticatedSession();
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
