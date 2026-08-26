import { CircularProgress, Paper, Stack, Typography } from '@mui/material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { useEffect, useState } from 'react';
import { Navigate, useLocation, useNavigate } from 'react-router-dom';

import { getApiAuthMe } from '../../../api/generated/client';
import { authMeQueryKey } from '../queryKeys';
import {
  canUseAuthenticatedSession,
  confirmAuthenticatedSession,
  consumeGoogleOAuthReturn,
  endAuthenticatedSession,
} from '../session';
import { authGateStyles } from '../styles/authGate.styles';

type RequireAuthProps = {
  children: ReactNode;
};

export default function RequireAuth({ children }: RequireAuthProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const redirectTo = `${location.pathname}${location.search}${location.hash}`;
  const [canUseSession] = useState(() => canUseAuthenticatedSession() || consumeGoogleOAuthReturn());
  const mustReauthenticate = !canUseSession;
  const authQuery = useQuery({
    queryKey: authMeQueryKey,
    queryFn: () => getApiAuthMe(),
    retry: false,
    enabled: canUseSession,
  });

  useEffect(() => {
    if (!mustReauthenticate) {
      return;
    }

    void endAuthenticatedSession(queryClient).finally(() => {
      navigate(`/?redirectTo=${encodeURIComponent(redirectTo)}`, { replace: true });
    });
  }, [mustReauthenticate, navigate, queryClient, redirectTo]);

  useEffect(() => {
    if (authQuery.isSuccess) {
      confirmAuthenticatedSession();
    }
  }, [authQuery.isSuccess]);

  if (mustReauthenticate || authQuery.isPending) {
    return (
      <Paper sx={authGateStyles.panel}>
        <Stack sx={authGateStyles.content}>
          <CircularProgress size={28} />
          <Typography color="text.secondary">Faça login novamente para continuar...</Typography>
        </Stack>
      </Paper>
    );
  }

  if (authQuery.isError) {
    return <Navigate to={`/?redirectTo=${encodeURIComponent(redirectTo)}`} replace />;
  }

  return children;
}
