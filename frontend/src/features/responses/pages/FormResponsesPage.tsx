import { Paper, Stack, Typography } from '@mui/material';
import { useParams } from 'react-router-dom';

import { responsesPageStyles } from '../styles/responsesPage.styles';

export default function FormResponsesPage() {
  const { formId } = useParams();

  return (
    <Stack sx={responsesPageStyles.pageStack}>
      <Typography variant="h4">Responses</Typography>
      <Paper sx={responsesPageStyles.panel}>
        <Typography color="text.secondary">Form ID: {formId}</Typography>
      </Paper>
    </Stack>
  );
}
