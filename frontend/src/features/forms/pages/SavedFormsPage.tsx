import { Alert, Button, Chip, Divider, LinearProgress, List, ListItemButton, ListItemText, Paper, Stack, Typography } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { Link as RouterLink } from 'react-router-dom';

import { getApiForms, getErrorMessage } from '../../../api/generated/client';
import { formsQueryKey } from '../queryKeys';
import { formPagesStyles } from '../styles/formPages.styles';
import type { FormStatus } from '../types';

const statusLabels: Record<FormStatus, string> = {
  draft: 'Rascunho',
  published: 'Publicado',
};

export default function SavedFormsPage() {
  const formsQuery = useQuery({
    queryKey: formsQueryKey,
    queryFn: () => getApiForms(),
  });
  const forms = formsQuery.data?.forms ?? [];

  return (
    <Stack sx={formPagesStyles.pageStack}>
      <Stack sx={formPagesStyles.header}>
        <Stack sx={formPagesStyles.titleBlock}>
          <Typography variant="h4">Formulários salvos</Typography>
          <Typography color="text.secondary">Acompanhe seus formulários antes de entrar na área administrativa.</Typography>
        </Stack>
        <Button component={RouterLink} to="/admin/workspace" variant="contained">
          Acessar área administrativa
        </Button>
      </Stack>

      {formsQuery.isFetching ? <LinearProgress sx={formPagesStyles.loadingBar} /> : null}
      {formsQuery.error ? <Alert severity="error">{getErrorMessage(formsQuery.error)}</Alert> : null}

      <Paper sx={formPagesStyles.savedFormsPanel}>
        <Stack sx={formPagesStyles.sidebarHeader}>
          <Typography variant="subtitle1">Formulários salvos</Typography>
          <Chip label={forms.length} size="small" />
        </Stack>
        <Divider />
        <List sx={formPagesStyles.savedFormsList}>
          {forms.length === 0 ? (
            <Stack sx={formPagesStyles.sidebarEmptyState}>
              <Typography variant="body2" color="text.secondary">
                Nenhum formulário salvo ainda.
              </Typography>
            </Stack>
          ) : null}
          {forms.map((form) => (
            <ListItemButton key={form.id} component={RouterLink} to="/admin/workspace" sx={formPagesStyles.savedFormItem}>
              <ListItemText
                primary={form.title}
                secondary={`${formatFieldCount(form.fields.length)} - ${form.description || 'Sem descrição'}`}
              />
              <Stack sx={formPagesStyles.savedFormMeta}>
                <Chip label={statusLabels[form.status]} size="small" color={form.status === 'published' ? 'success' : 'default'} />
              </Stack>
            </ListItemButton>
          ))}
        </List>
      </Paper>
    </Stack>
  );
}

function formatFieldCount(count: number) {
  return `${count} ${count === 1 ? 'campo' : 'campos'}`;
}
