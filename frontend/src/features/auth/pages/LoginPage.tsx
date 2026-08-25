import { Button, Paper, Stack, TextField, Typography } from '@mui/material';

import { loginPageStyles } from '../styles/loginPage.styles';

export default function LoginPage() {
  return (
    <Paper sx={loginPageStyles.panel}>
      <Stack sx={loginPageStyles.form}>
        <Typography variant="h5">Admin login</Typography>
        <TextField label="Email" type="email" autoComplete="email" />
        <TextField label="Password" type="password" autoComplete="current-password" />
        <Button variant="contained">Sign in</Button>
        <Button variant="outlined">Continue with Google</Button>
      </Stack>
    </Paper>
  );
}
