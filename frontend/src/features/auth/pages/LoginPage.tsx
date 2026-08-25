import { Alert, Button, Paper, Stack, TextField, Typography } from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { getErrorMessage } from '../../../lib/api';
import { login, register } from '../api/authApi';

import { loginPageStyles } from '../styles/loginPage.styles';

type AuthMode = 'login' | 'register';

export default function LoginPage() {
  const navigate = useNavigate();
  const [mode, setMode] = useState<AuthMode>('login');
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const authMutation = useMutation({
    mutationFn: () => {
      if (mode === 'register') {
        return register({ name, email, password });
      }

      return login({ email, password });
    },
    onSuccess: () => {
      navigate('/admin');
    },
  });

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    authMutation.mutate();
  }

  function toggleMode() {
    authMutation.reset();
    setMode((currentMode) => (currentMode === 'login' ? 'register' : 'login'));
  }

  return (
    <Paper sx={loginPageStyles.panel}>
      <Stack component="form" onSubmit={handleSubmit} sx={loginPageStyles.form}>
        <Stack sx={loginPageStyles.header}>
          <Typography variant="h5">{mode === 'login' ? 'Admin login' : 'Create account'}</Typography>
          <Typography color="text.secondary">Form Builder</Typography>
        </Stack>

        {authMutation.error ? <Alert severity="error">{getErrorMessage(authMutation.error)}</Alert> : null}

        {mode === 'register' ? (
          <TextField
            label="Name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            autoComplete="name"
            required
          />
        ) : null}

        <TextField
          label="Email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          autoComplete="email"
          required
        />
        <TextField
          label="Password"
          type="password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          autoComplete={mode === 'register' ? 'new-password' : 'current-password'}
          required
        />

        <Stack sx={loginPageStyles.actions}>
          <Button type="submit" variant="contained" disabled={authMutation.isPending}>
            {mode === 'login' ? 'Sign in' : 'Create'}
          </Button>
          <Button type="button" variant="text" onClick={toggleMode}>
            {mode === 'login' ? 'Create account' : 'Use existing account'}
          </Button>
        </Stack>
      </Stack>
    </Paper>
  );
}
