import { Paper, Stack, Typography } from '@mui/material';
import { useParams } from 'react-router-dom';

export default function FormResponsesPage() {
  const { formId } = useParams();

  return (
    <Stack spacing={3}>
      <Typography variant="h4">Responses</Typography>
      <Paper sx={{ p: 3 }}>
        <Typography color="text.secondary">Form ID: {formId}</Typography>
      </Paper>
    </Stack>
  );
}
