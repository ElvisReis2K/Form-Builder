import { Button, Paper, Stack, TextField, Typography } from '@mui/material';

export default function AdminHomePage() {
  return (
    <Stack spacing={3}>
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} justifyContent="space-between">
        <Typography variant="h4">Forms</Typography>
        <Button variant="contained">New form</Button>
      </Stack>

      <Paper sx={{ p: 3 }}>
        <Stack spacing={2}>
          <TextField label="Form title" defaultValue="Customer feedback" />
          <TextField label="Description" defaultValue="Simple published form draft" multiline rows={3} />
          <Button variant="outlined" sx={{ alignSelf: 'flex-start' }}>
            Publish
          </Button>
        </Stack>
      </Paper>
    </Stack>
  );
}
