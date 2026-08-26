import {
  Alert,
  Box,
  Button,
  Chip,
  Divider,
  FormControlLabel,
  LinearProgress,
  List,
  ListItemButton,
  ListItemText,
  MenuItem,
  Paper,
  Stack,
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useMemo, useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';

import {
  deleteApiFormsFormId,
  getApiForms,
  getErrorMessage,
  postApiForms,
  postApiFormsFormIdPublish,
  postApiFormsFormIdUnpublish,
  putApiFormsFormId,
} from '../../../api/generated/client';
import { formsQueryKey } from '../queryKeys';
import { formPagesStyles } from '../styles/formPages.styles';
import type { FieldType, FormStatus } from '../types';
import {
  createBlankDraft,
  createFieldDraft,
  draftToRequest,
  fieldTypeLabels,
  fieldTypes,
  formToDraft,
  type FieldDraft,
  type FormDraft,
} from '../utils/formDraft';

const statusLabels: Record<FormStatus, string> = {
  draft: 'Rascunho',
  published: 'Publicado',
};

export default function AdminHomePage() {
  const queryClient = useQueryClient();
  const [selectedFormId, setSelectedFormId] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(true);
  const [draft, setDraft] = useState<FormDraft>(() => createBlankDraft());

  const formsQuery = useQuery({
    queryKey: formsQueryKey,
    queryFn: () => getApiForms(),
  });

  const forms = formsQuery.data?.forms ?? [];
  const selectedForm = useMemo(
    () => forms.find((form) => form.id === selectedFormId) ?? null,
    [forms, selectedFormId],
  );

  useEffect(() => {
    if (!isCreating && selectedForm) {
      setDraft(formToDraft(selectedForm));
    }
  }, [isCreating, selectedForm]);

  const createMutation = useMutation({
    mutationFn: (input: ReturnType<typeof draftToRequest>) => postApiForms({ body: input }),
    onSuccess: (form) => {
      setIsCreating(false);
      setSelectedFormId(form.id);
      setDraft(formToDraft(form));
      void queryClient.invalidateQueries({ queryKey: formsQueryKey });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ formId, input }: { formId: string; input: ReturnType<typeof draftToRequest> }) =>
      putApiFormsFormId({
        path: { formId },
        body: input,
      }),
    onSuccess: (form) => {
      setDraft(formToDraft(form));
      void queryClient.invalidateQueries({ queryKey: formsQueryKey });
    },
  });

  const publishMutation = useMutation({
    mutationFn: (formId: string) =>
      postApiFormsFormIdPublish({
        path: { formId },
      }),
    onSuccess: (form) => {
      setDraft(formToDraft(form));
      void queryClient.invalidateQueries({ queryKey: formsQueryKey });
    },
  });

  const unpublishMutation = useMutation({
    mutationFn: (formId: string) =>
      postApiFormsFormIdUnpublish({
        path: { formId },
      }),
    onSuccess: (form) => {
      setDraft(formToDraft(form));
      void queryClient.invalidateQueries({ queryKey: formsQueryKey });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (formId: string) =>
      deleteApiFormsFormId({
        path: { formId },
      }),
    onSuccess: () => {
      setIsCreating(true);
      setSelectedFormId(null);
      setDraft(createBlankDraft());
      void queryClient.invalidateQueries({ queryKey: formsQueryKey });
    },
  });

  const activeError =
    formsQuery.error ??
    createMutation.error ??
    updateMutation.error ??
    publishMutation.error ??
    unpublishMutation.error ??
    deleteMutation.error;
  const isBusy =
    formsQuery.isFetching ||
    createMutation.isPending ||
    updateMutation.isPending ||
    publishMutation.isPending ||
    unpublishMutation.isPending ||
    deleteMutation.isPending;
  const selectedStatus = selectedForm?.status ?? 'draft';
  const hasSavedForm = draft.id !== null;
  const hasPrivacyNotice =
    draft.controllerEmail.trim() !== '' && draft.privacyPurpose.trim() !== '' && draft.retentionPolicy.trim() !== '';

  function startNewForm() {
    setIsCreating(true);
    setSelectedFormId(null);
    setDraft(createBlankDraft());
    resetMutationErrors();
  }

  function selectForm(formId: string) {
    setIsCreating(false);
    setSelectedFormId(formId);
    resetMutationErrors();
  }

  function saveForm() {
    const input = draftToRequest(draft);
    if (draft.id === null) {
      createMutation.mutate(input);
      return;
    }

    updateMutation.mutate({ formId: draft.id, input });
  }

  function publishSelectedForm() {
    if (draft.id !== null) {
      publishMutation.mutate(draft.id);
    }
  }

  function unpublishSelectedForm() {
    if (draft.id !== null) {
      unpublishMutation.mutate(draft.id);
    }
  }

  function deleteSelectedForm() {
    if (draft.id !== null && window.confirm('Excluir este formulário?')) {
      deleteMutation.mutate(draft.id);
    }
  }

  function updateDraft(patch: Partial<FormDraft>) {
    setDraft((currentDraft) => ({
      ...currentDraft,
      ...patch,
    }));
  }

  function updateField(clientId: string, patch: Partial<FieldDraft>) {
    setDraft((currentDraft) => ({
      ...currentDraft,
      fields: currentDraft.fields.map((field) => (field.clientId === clientId ? { ...field, ...patch } : field)),
    }));
  }

  function addField() {
    setDraft((currentDraft) => ({
      ...currentDraft,
      fields: [...currentDraft.fields, createFieldDraft()],
    }));
  }

  function removeField(clientId: string) {
    setDraft((currentDraft) => ({
      ...currentDraft,
      fields: currentDraft.fields.filter((field) => field.clientId !== clientId),
    }));
  }

  function resetMutationErrors() {
    createMutation.reset();
    updateMutation.reset();
    publishMutation.reset();
    unpublishMutation.reset();
    deleteMutation.reset();
  }

  return (
    <Stack sx={formPagesStyles.pageStack}>
      <Stack sx={formPagesStyles.header}>
        <Stack sx={formPagesStyles.titleBlock}>
          <Typography variant="h4">Formulários</Typography>
          <Typography color="text.secondary">Área administrativa</Typography>
        </Stack>
        <Button variant="contained" onClick={startNewForm}>
          Novo formulário
        </Button>
      </Stack>

      {isBusy ? <LinearProgress sx={formPagesStyles.loadingBar} /> : null}
      {activeError ? <Alert severity="error">{getErrorMessage(activeError)}</Alert> : null}

      <Box sx={formPagesStyles.workspace}>
        <Paper sx={formPagesStyles.sidebarPanel}>
          <Stack sx={formPagesStyles.sidebarHeader}>
            <Typography variant="subtitle1">Formulários salvos</Typography>
            <Chip label={forms.length} size="small" />
          </Stack>
          <Divider />
          <List sx={formPagesStyles.formsList}>
            {forms.length === 0 ? (
              <Stack sx={formPagesStyles.sidebarEmptyState}>
                <Typography variant="body2" color="text.secondary">
                  Nenhum formulário salvo
                </Typography>
              </Stack>
            ) : null}
            {forms.map((form) => (
              <ListItemButton
                key={form.id}
                selected={!isCreating && selectedFormId === form.id}
                onClick={() => selectForm(form.id)}
                sx={formPagesStyles.formListItem}
              >
                <ListItemText primary={form.title} secondary={formatFieldCount(form.fields.length)} />
                <Chip
                  label={statusLabels[form.status]}
                  size="small"
                  color={form.status === 'published' ? 'success' : 'default'}
                />
              </ListItemButton>
            ))}
          </List>
        </Paper>

        <Paper sx={formPagesStyles.editorPanel}>
          <Stack sx={formPagesStyles.editorStack}>
            <Stack sx={formPagesStyles.editorHeader}>
              <Stack sx={formPagesStyles.titleBlock}>
                <Typography variant="h5">{isCreating ? 'Novo formulário' : 'Editor de formulário'}</Typography>
                <Typography color="text.secondary">{hasSavedForm ? draft.id : 'Rascunho não salvo'}</Typography>
              </Stack>
              <Chip
                label={hasSavedForm ? statusLabels[selectedStatus] : statusLabels.draft}
                color={selectedStatus === 'published' ? 'success' : 'default'}
              />
            </Stack>

            <Box sx={formPagesStyles.formGrid}>
              <TextField
                label="Titulo"
                value={draft.title}
                onChange={(event) => updateDraft({ title: event.target.value })}
                required
              />
              <TextField
                label="Descricao"
                value={draft.description}
                onChange={(event) => updateDraft({ description: event.target.value })}
                multiline
                minRows={3}
              />
            </Box>

            <Divider />

            <Stack sx={formPagesStyles.sectionHeader}>
              <Typography variant="h6">Privacidade</Typography>
            </Stack>

            <Box sx={formPagesStyles.formGrid}>
              <TextField
                label="E-mail do controlador"
                type="email"
                value={draft.controllerEmail}
                onChange={(event) => updateDraft({ controllerEmail: event.target.value })}
                required
              />
              <TextField
                label="Finalidade do tratamento"
                value={draft.privacyPurpose}
                onChange={(event) => updateDraft({ privacyPurpose: event.target.value })}
                multiline
                minRows={3}
                required
              />
              <TextField
                label="Retenção das respostas"
                value={draft.retentionPolicy}
                onChange={(event) => updateDraft({ retentionPolicy: event.target.value })}
                multiline
                minRows={3}
                required
              />
            </Box>

            <Divider />

            <Stack sx={formPagesStyles.sectionHeader}>
              <Typography variant="h6">Campos</Typography>
              <Button variant="outlined" onClick={addField}>
                Adicionar campo
              </Button>
            </Stack>

            <Stack sx={formPagesStyles.fieldStack}>
              {draft.fields.map((field, index) => (
                <Stack key={field.clientId} sx={formPagesStyles.fieldEditor}>
                  <Stack sx={formPagesStyles.fieldHeader}>
                    <Typography variant="subtitle2">Campo {index + 1}</Typography>
                    <Button variant="text" color="error" onClick={() => removeField(field.clientId)}>
                      Remover
                    </Button>
                  </Stack>

                  <Box sx={formPagesStyles.fieldGrid}>
                    <TextField
                      select
                      label="Tipo"
                      value={field.type}
                      onChange={(event) => updateField(field.clientId, { type: event.target.value as FieldType })}
                      size="small"
                    >
                      {fieldTypes.map((fieldType) => (
                        <MenuItem key={fieldType} value={fieldType}>
                          {fieldTypeLabels[fieldType]}
                        </MenuItem>
                      ))}
                    </TextField>
                    <TextField
                      label="Rótulo"
                      value={field.label}
                      onChange={(event) => updateField(field.clientId, { label: event.target.value })}
                      size="small"
                      required
                    />
                    <TextField
                      label="Texto de ajuda"
                      value={field.placeholder}
                      onChange={(event) => updateField(field.clientId, { placeholder: event.target.value })}
                      size="small"
                    />
                    <FormControlLabel
                      control={
                        <Switch
                          checked={field.required}
                          onChange={(_, checked) => updateField(field.clientId, { required: checked })}
                        />
                      }
                      label="Obrigatório"
                    />
                  </Box>

                  {field.type === 'select' ? (
                    <TextField
                      label="Opções"
                      value={field.optionsText}
                      onChange={(event) => updateField(field.clientId, { optionsText: event.target.value })}
                      multiline
                      minRows={3}
                    />
                  ) : null}
                </Stack>
              ))}
            </Stack>

            <Divider />

            <Stack sx={formPagesStyles.actionBar}>
              <Button variant="contained" onClick={saveForm} disabled={isBusy || draft.title.trim() === ''}>
                Salvar
              </Button>
              {selectedStatus === 'published' ? (
                <Button variant="outlined" onClick={unpublishSelectedForm} disabled={!hasSavedForm || isBusy}>
                  Despublicar
                </Button>
              ) : (
                <Button variant="outlined" onClick={publishSelectedForm} disabled={!hasSavedForm || isBusy || !hasPrivacyNotice}>
                  Publicar
                </Button>
              )}
              {selectedForm?.publicUrl ? (
                <Button component={RouterLink} to={selectedForm.publicUrl} variant="text">
                  Abrir formulário público
                </Button>
              ) : null}
              {hasSavedForm ? (
                <Button component={RouterLink} to={`/admin/forms/${draft.id}/responses`} variant="text">
                  Respostas
                </Button>
              ) : null}
              <Button color="error" variant="text" onClick={deleteSelectedForm} disabled={!hasSavedForm || isBusy}>
                Excluir
              </Button>
            </Stack>
          </Stack>
        </Paper>
      </Box>
    </Stack>
  );
}

function formatFieldCount(count: number) {
  return `${count} ${count === 1 ? 'campo' : 'campos'}`;
}
