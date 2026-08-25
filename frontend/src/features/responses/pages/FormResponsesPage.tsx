import {
  Alert,
  Button,
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
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';

import { deleteFormResponse, exportFormResponses, listFormResponses } from '../../forms/api/formsApi';
import type { FormSubmission, FormSubmissionExportResponse } from '../../forms/types';
import { getErrorMessage } from '../../../lib/api';
import { responsesPageStyles } from '../styles/responsesPage.styles';

export default function FormResponsesPage() {
  const { formId } = useParams();
  const queryClient = useQueryClient();

  const responsesQuery = useQuery({
    queryKey: ['form-responses', formId],
    queryFn: () => listFormResponses(formId ?? ''),
    enabled: Boolean(formId),
  });

  const form = responsesQuery.data?.form;
  const responses = responsesQuery.data?.responses ?? [];
  const exportMutation = useMutation({
    mutationFn: () => exportFormResponses(formId ?? ''),
    onSuccess: (payload) => {
      downloadJSON(payload, responseExportFilename(form?.title ?? formId ?? 'formulario'));
    },
  });
  const deleteMutation = useMutation({
    mutationFn: (responseId: string) => deleteFormResponse(formId ?? '', responseId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['form-responses', formId] });
    },
  });
  const activeError = responsesQuery.error ?? exportMutation.error ?? deleteMutation.error;
  const isBusy = responsesQuery.isFetching || exportMutation.isPending || deleteMutation.isPending;

  return (
    <Stack sx={responsesPageStyles.pageStack}>
      <Stack sx={responsesPageStyles.header}>
        <Stack sx={responsesPageStyles.titleBlock}>
          <Typography variant="h4">Respostas</Typography>
          {form ? <Typography color="text.secondary">{form.title}</Typography> : null}
        </Stack>
        <Button variant="outlined" onClick={() => exportMutation.mutate()} disabled={!form || responses.length === 0 || isBusy}>
          Exportar JSON
        </Button>
      </Stack>

      {isBusy ? <LinearProgress sx={responsesPageStyles.loadingBar} /> : null}
      {activeError ? <Alert severity="error">{getErrorMessage(activeError)}</Alert> : null}

      <Paper sx={responsesPageStyles.panel}>
        {form && responses.length > 0 ? (
          <TableContainer sx={responsesPageStyles.tableContainer}>
            <Table size="small" sx={responsesPageStyles.table}>
              <TableHead sx={responsesPageStyles.tableHead}>
                <TableRow>
                  <TableCell>Enviada em</TableCell>
                  <TableCell>Ciencia LGPD</TableCell>
                  {form.fields.map((field) => (
                    <TableCell key={field.id}>{field.label}</TableCell>
                  ))}
                  <TableCell>Acoes</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {responses.map((response) => (
                  <TableRow key={response.id} hover>
                    <TableCell>{formatDate(response.submittedAt)}</TableCell>
                    <TableCell>{formatDate(response.privacyAcknowledgedAt)}</TableCell>
                    {form.fields.map((field) => (
                      <TableCell key={field.id}>{formatAnswer(response, field.id)}</TableCell>
                    ))}
                    <TableCell>
                      <Button
                        variant="text"
                        color="error"
                        size="small"
                        onClick={() => deleteSelectedResponse(response.id, deleteMutation.mutate)}
                        disabled={deleteMutation.isPending}
                      >
                        Excluir
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        ) : (
          <Stack sx={responsesPageStyles.emptyState}>
            <Typography variant="h6">Nenhuma resposta ainda</Typography>
            <Typography color="text.secondary">As respostas de formularios publicados aparecerao aqui.</Typography>
          </Stack>
        )}
      </Paper>
    </Stack>
  );
}

function deleteSelectedResponse(responseId: string, deleteResponse: (responseId: string) => void) {
  if (window.confirm('Excluir esta resposta?')) {
    deleteResponse(responseId);
  }
}

function downloadJSON(payload: FormSubmissionExportResponse, filename: string) {
  const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}

function responseExportFilename(title: string) {
  const slug = title
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');

  return `respostas-${slug || 'formulario'}.json`;
}

function formatAnswer(response: FormSubmission, fieldId: string) {
  const value = response.answers[fieldId];
  if (value === undefined || value === null || value === '') {
    return '-';
  }

  if (typeof value === 'boolean') {
    return value ? 'Sim' : 'Nao';
  }

  return String(value);
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(new Date(value));
}
