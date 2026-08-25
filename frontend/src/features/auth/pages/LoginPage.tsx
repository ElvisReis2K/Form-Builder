import { Alert, Button, Paper, Stack, TextField, Typography } from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import { FormEvent, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { getErrorMessage } from '../../../lib/api';
import { googleLoginURL, login, register } from '../api/authApi';

import { loginPageStyles } from '../styles/loginPage.styles';

type AuthMode = 'login' | 'register';

export default function LoginPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
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

  function startGoogleLogin() {
    window.location.assign(googleLoginURL());
  }

  const oauthError = searchParams.get('authError');

  return (
    <Paper sx={loginPageStyles.panel}>
      <Stack component="form" onSubmit={handleSubmit} sx={loginPageStyles.form}>
        <Stack sx={loginPageStyles.header}>
          <Typography variant="h5">{mode === 'login' ? 'Entrar na administracao' : 'Criar conta'}</Typography>
          <Typography color="text.secondary">Construtor de Formularios</Typography>
        </Stack>

        {authMutation.error ? <Alert severity="error">{getErrorMessage(authMutation.error)}</Alert> : null}
        {oauthError ? <Alert severity="error">Falha ao entrar com Google. Tente novamente.</Alert> : null}

        {mode === 'register' ? (
          <TextField
            label="Nome"
            value={name}
            onChange={(event) => setName(event.target.value)}
            autoComplete="name"
            required
          />
        ) : null}

        <TextField
          label="E-mail"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          autoComplete="email"
          required
        />
        <TextField
          label="Senha"
          type="password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          autoComplete={mode === 'register' ? 'new-password' : 'current-password'}
          required
        />

        <Stack sx={loginPageStyles.actions}>
          <Button type="submit" variant="contained" disabled={authMutation.isPending} sx={loginPageStyles.actionButton}>
            {mode === 'login' ? 'Entrar' : 'Criar'}
          </Button>
          <Button type="button" variant="text" onClick={toggleMode} sx={loginPageStyles.actionButton}>
            {mode === 'login' ? 'Criar conta' : 'Usar conta existente'}
          </Button>
        </Stack>

        <Button type="button" variant="outlined" onClick={startGoogleLogin} sx={loginPageStyles.googleButton}>
          Continuar com Google
        </Button>
      </Stack>
    </Paper>
  );
}
