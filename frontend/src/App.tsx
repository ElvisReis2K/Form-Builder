import { AppBar, Box, Button, Container, Stack, Toolbar, Typography } from '@mui/material';
import { Link as RouterLink, Route, Routes } from 'react-router-dom';

import AdminHomePage from './pages/AdminHomePage';
import FormResponsesPage from './pages/FormResponsesPage';
import LoginPage from './pages/LoginPage';
import PublicFormPage from './pages/PublicFormPage';

const navItems = [
  { label: 'Admin', to: '/admin' },
  { label: 'Responses', to: '/admin/forms/demo/responses' },
  { label: 'Public form', to: '/f/demo' },
];

export default function App() {
  return (
    <Box minHeight="100vh">
      <AppBar position="static" color="default" elevation={0}>
        <Toolbar sx={{ gap: 2 }}>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Form Builder
          </Typography>
          <Stack direction="row" spacing={1}>
            {navItems.map((item) => (
              <Button key={item.to} component={RouterLink} to={item.to} size="small">
                {item.label}
              </Button>
            ))}
          </Stack>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Routes>
          <Route path="/" element={<LoginPage />} />
          <Route path="/admin" element={<AdminHomePage />} />
          <Route path="/admin/forms/:formId/responses" element={<FormResponsesPage />} />
          <Route path="/f/:slug" element={<PublicFormPage />} />
        </Routes>
      </Container>
    </Box>
  );
}
