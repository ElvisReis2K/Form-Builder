import { Button, Paper, Stack, TextField, Typography } from '@mui/material';
import { useParams } from 'react-router-dom';

export default function PublicFormPage() {
  const { slug } = useParams();

  return (
    <Paper sx={{ maxWidth: 640, mx: 'auto', p: 3 }}>
      <Stack spacing={2}>
        <Typography variant="h4">Published form</Typography>
        <Typography color="text.secondary">Slug: {slug}</Typography>
        <TextField label="Name" />
        <TextField label="Email" type="email" />
        <TextField label="Message" multiline rows={4} />
        <Button variant="contained">Submit response</Button>
      </Stack>
    </Paper>
  );
}
