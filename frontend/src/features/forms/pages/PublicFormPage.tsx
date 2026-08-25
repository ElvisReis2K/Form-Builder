import { Button, Paper, Stack, TextField, Typography } from '@mui/material';
import { useParams } from 'react-router-dom';

import { formPagesStyles } from '../styles/formPages.styles';

export default function PublicFormPage() {
  const { slug } = useParams();

  return (
    <Paper sx={formPagesStyles.publicPanel}>
      <Stack sx={formPagesStyles.fieldStack}>
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
