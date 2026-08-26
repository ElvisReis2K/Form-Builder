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

import {
  deleteApiFormsFormIdResponsesResponseId,
  getApiFormsFormIdResponses,
  getApiFormsFormIdResponsesExport,
  getErrorMessage,
} from '../../../api/generated/client';
import { responsesPageStyles } from '../styles/responsesPage.styles';
import {
  downloadResponsesExcel,
  downloadResponsesJSON,
  downloadResponsesPDF,
  formatAnswer,
  formatDate,
} from '../utils/responseExports';

export default function FormResponsesPage() {
  const { formId } = useParams();
  const queryClient = useQueryClient();

  const responsesQuery = useQuery({
    queryKey: ['form-responses', formId],
    queryFn: () =>
      getApiFormsFormIdResponses({
        path: { formId: formId ?? '' },
      }),
    enabled: Boolean(formId),
  });

  const form = responsesQuery.data?.form;
  const responses = responsesQuery.data?.responses ?? [];
  const exportMutation = useMutation({
    mutationFn: () =>
      getApiFormsFormIdResponsesExport({
        path: { formId: formId ?? '' },
      }),
    onSuccess: (payload) => {
      downloadResponsesJSON(payload);
    },
  });
  const deleteMutation = useMutation({
    mutationFn: (responseId: string) =>
      deleteApiFormsFormIdResponsesResponseId({
        path: { formId: formId ?? '', responseId },
      }),
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
        <Stack sx={responsesPageStyles.exportActions}>
          <Button
            variant="outlined"
            onClick={() => exportMutation.mutate()}
            disabled={!form || responses.length === 0 || isBusy}
          >
            Exportar JSON
          </Button>
          <Button
            variant="outlined"
            onClick={() => form && downloadResponsesPDF(form, responses)}
            disabled={!form || responses.length === 0 || isBusy}
          >
            Exportar PDF
          </Button>
          <Button
            variant="outlined"
            onClick={() => form && downloadResponsesExcel(form, responses)}
            disabled={!form || responses.length === 0 || isBusy}
          >
            Exportar Excel
          </Button>
        </Stack>
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
                  <TableCell>Ciência LGPD</TableCell>
                  {form.fields.map((field) => (
                    <TableCell key={field.id}>{field.label}</TableCell>
                  ))}
                  <TableCell>Ações</TableCell>
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
            <Typography color="text.secondary">As respostas de formulários publicados aparecerão aqui.</Typography>
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
