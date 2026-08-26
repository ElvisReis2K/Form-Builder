import { AppBar, Box, Button, Container, Stack, Toolbar, Typography } from '@mui/material';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Link as RouterLink, useLocation, useNavigate } from 'react-router-dom';

import { getApiAuthMe } from '../api/generated/client';
import { authMeQueryKey } from '../features/auth/queryKeys';
import { canUseAuthenticatedSession, endAuthenticatedSession } from '../features/auth/session';
import { appStyles } from './app.styles';
import AppRoutes from './AppRoutes';
import { authenticatedNavItems, guestNavItems } from './navigation';

export default function App() {
  const location = useLocation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const authQuery = useQuery({
    queryKey: authMeQueryKey,
    queryFn: () => getApiAuthMe(),
    retry: false,
    enabled: canUseAuthenticatedSession(),
  });
  const isLoginPage = location.pathname === '/';
  const canShowAuthenticatedActions = authQuery.isSuccess && !isLoginPage;
  const navItems = canShowAuthenticatedActions ? authenticatedNavItems : guestNavItems;
  const logoutMutation = useMutation({
    mutationFn: () => endAuthenticatedSession(queryClient),
    onSettled: () => {
      navigate('/', { replace: true });
    },
  });

  return (
    <Box sx={appStyles.root}>
      <AppBar position="static" color="default" elevation={0} sx={appStyles.appBar}>
        <Toolbar sx={appStyles.toolbar}>
          <Typography variant="h6" component="div" sx={appStyles.brand}>
            Construtor de Formulários
          </Typography>
          <Stack sx={appStyles.nav}>
            {navItems.map((item) => (
              <Button key={item.to} component={RouterLink} to={item.to} size="small" sx={appStyles.navButton}>
                {item.label}
              </Button>
            ))}
            {canShowAuthenticatedActions ? (
              <Button
                type="button"
                size="small"
                sx={appStyles.navButton}
                onClick={() => logoutMutation.mutate()}
                disabled={logoutMutation.isPending}
              >
                {logoutMutation.isPending ? 'Saindo...' : 'Sair'}
              </Button>
            ) : null}
          </Stack>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={appStyles.content}>
        <AppRoutes />
      </Container>
    </Box>
  );
}
