import {
  Alert,
  LinearProgress,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';

import { listFormResponses } from '../../forms/api/formsApi';
import type { FormSubmission } from '../../forms/types';
import { getErrorMessage } from '../../../lib/api';
import { responsesPageStyles } from '../styles/responsesPage.styles';

export default function FormResponsesPage() {
  const { formId } = useParams();

  const responsesQuery = useQuery({
    queryKey: ['form-responses', formId],
    queryFn: () => listFormResponses(formId ?? ''),
    enabled: Boolean(formId),
  });

  const form = responsesQuery.data?.form;
  const responses = responsesQuery.data?.responses ?? [];

  return (
    <Stack sx={responsesPageStyles.pageStack}>
      <Stack sx={responsesPageStyles.header}>
        <Typography variant="h4">Responses</Typography>
        {form ? <Typography color="text.secondary">{form.title}</Typography> : null}
      </Stack>

      {responsesQuery.isFetching ? <LinearProgress sx={responsesPageStyles.loadingBar} /> : null}
      {responsesQuery.error ? <Alert severity="error">{getErrorMessage(responsesQuery.error)}</Alert> : null}

      <Paper sx={responsesPageStyles.panel}>
        {form && responses.length > 0 ? (
          <TableContainer sx={responsesPageStyles.tableContainer}>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Submitted at</TableCell>
                  {form.fields.map((field) => (
                    <TableCell key={field.id}>{field.label}</TableCell>
                  ))}
                </TableRow>
              </TableHead>
              <TableBody>
                {responses.map((response) => (
                  <TableRow key={response.id}>
                    <TableCell>{formatDate(response.submittedAt)}</TableCell>
                    {form.fields.map((field) => (
                      <TableCell key={field.id}>{formatAnswer(response, field.id)}</TableCell>
                    ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        ) : (
          <Stack sx={responsesPageStyles.emptyState}>
            <Typography variant="h6">No responses yet</Typography>
            <Typography color="text.secondary">Published forms will show submissions here.</Typography>
          </Stack>
        )}
      </Paper>
    </Stack>
  );
}

function formatAnswer(response: FormSubmission, fieldId: string) {
  const value = response.answers[fieldId];
  if (value === undefined || value === null || value === '') {
    return '-';
  }

  if (typeof value === 'boolean') {
    return value ? 'Yes' : 'No';
  }

  return String(value);
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(new Date(value));
}
