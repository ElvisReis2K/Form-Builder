import { AppBar, Box, Button, Container, Stack, Toolbar, Typography } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { Link as RouterLink, useLocation } from 'react-router-dom';

import { getApiAuthMe } from '../api/generated/client';
import { authMeQueryKey } from '../features/auth/queryKeys';
import { appStyles } from './app.styles';
import AppRoutes from './AppRoutes';
import { authenticatedNavItems, guestNavItems } from './navigation';

export default function App() {
  const location = useLocation();
  const authQuery = useQuery({
    queryKey: authMeQueryKey,
    queryFn: () => getApiAuthMe(),
    retry: false,
  });
  const isLoginPage = location.pathname === '/';
  const navItems = authQuery.isSuccess && !isLoginPage ? authenticatedNavItems : guestNavItems;

  return (
    <Box sx={appStyles.root}>
      <AppBar position="static" color="default" elevation={0} sx={appStyles.appBar}>
        <Toolbar sx={appStyles.toolbar}>
          <Typography variant="h6" component="div" sx={appStyles.brand}>
            Construtor de Formularios
          </Typography>
          <Stack sx={appStyles.nav}>
            {navItems.map((item) => (
              <Button key={item.to} component={RouterLink} to={item.to} size="small" sx={appStyles.navButton}>
                {item.label}
              </Button>
            ))}
          </Stack>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={appStyles.content}>
        <AppRoutes />
      </Container>
    </Box>
  );
}
