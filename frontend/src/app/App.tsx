import { AppBar, Box, Button, Container, Stack, Toolbar, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appStyles } from './app.styles';
import AppRoutes from './AppRoutes';
import { navItems } from './navigation';

export default function App() {
  return (
    <Box sx={appStyles.root}>
      <AppBar position="static" color="default" elevation={0}>
        <Toolbar sx={appStyles.toolbar}>
          <Typography variant="h6" component="div" sx={appStyles.brand}>
            Form Builder
          </Typography>
          <Stack sx={appStyles.nav}>
            {navItems.map((item) => (
              <Button key={item.to} component={RouterLink} to={item.to} size="small">
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
