import { CircularProgress, Paper, Stack, Typography } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';

import { getApiAuthMe } from '../../../api/generated/client';
import { authMeQueryKey } from '../queryKeys';
import { authGateStyles } from '../styles/authGate.styles';

type RequireAuthProps = {
  children: ReactNode;
};

export default function RequireAuth({ children }: RequireAuthProps) {
  const location = useLocation();
  const authQuery = useQuery({
    queryKey: authMeQueryKey,
    queryFn: () => getApiAuthMe(),
    retry: false,
  });

  if (authQuery.isPending) {
    return (
      <Paper sx={authGateStyles.panel}>
        <Stack sx={authGateStyles.content}>
          <CircularProgress size={28} />
          <Typography color="text.secondary">Verificando login...</Typography>
        </Stack>
      </Paper>
    );
  }

  if (authQuery.isError) {
    const redirectTo = `${location.pathname}${location.search}${location.hash}`;
    return <Navigate to={`/?redirectTo=${encodeURIComponent(redirectTo)}`} replace />;
  }

  return children;
}
