import { Button, Paper, Stack, TextField, Typography } from '@mui/material';

export default function LoginPage() {
  return (
    <Paper sx={{ maxWidth: 420, mx: 'auto', p: 3 }}>
      <Stack spacing={2}>
        <Typography variant="h5">Admin login</Typography>
        <TextField label="Email" type="email" autoComplete="email" />
        <TextField label="Password" type="password" autoComplete="current-password" />
        <Button variant="contained">Sign in</Button>
        <Button variant="outlined">Continue with Google</Button>
      </Stack>
    </Paper>
  );
}
