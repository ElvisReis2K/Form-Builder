import {
  Alert,
  Button,
  Checkbox,
  FormControlLabel,
  Link,
  LinearProgress,
  MenuItem,
  Paper,
  Stack,
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import { useMutation, useQuery } from '@tanstack/react-query';
import { FormEvent, useEffect, useState } from 'react';
import { Link as RouterLink, useParams } from 'react-router-dom';

import {
  getApiPublicFormsSlug,
  getErrorMessage,
  postApiPublicFormsSlugResponses,
} from '../../../api/generated/client';
import { formPagesStyles } from '../styles/formPages.styles';
import type { FormField } from '../types';

type AnswerState = Record<string, string | boolean>;

export default function PublicFormPage() {
  const { slug } = useParams();
  const [answers, setAnswers] = useState<AnswerState>({});
  const [privacyAcknowledged, setPrivacyAcknowledged] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const formQuery = useQuery({
    queryKey: ['public-form', slug],
    queryFn: () =>
      getApiPublicFormsSlug({
        path: { slug: slug ?? '' },
      }),
    enabled: Boolean(slug),
  });

  const form = formQuery.data;
  const submitMutation = useMutation({
    mutationFn: () =>
      postApiPublicFormsSlugResponses({
        path: { slug: slug ?? '' },
        body: { answers, privacyAcknowledged },
      }),
    onSuccess: () => {
      setSubmitted(true);
    },
  });

  useEffect(() => {
    if (!form) {
      return;
    }

    const initialAnswers: AnswerState = {};
    for (const field of form.fields) {
      initialAnswers[field.id] = field.type === 'checkbox' ? false : '';
    }

    setAnswers(initialAnswers);
    setPrivacyAcknowledged(false);
    setSubmitted(false);
  }, [form]);

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitted(false);
    submitMutation.mutate();
  }

  function updateAnswer(fieldId: string, value: string | boolean) {
    setAnswers((currentAnswers) => ({
      ...currentAnswers,
      [fieldId]: value,
    }));
  }

  return (
    <Paper sx={formPagesStyles.publicPanel}>
      <Stack component="form" onSubmit={handleSubmit} sx={formPagesStyles.fieldStack}>
        {formQuery.isFetching || submitMutation.isPending ? <LinearProgress sx={formPagesStyles.loadingBar} /> : null}
        {formQuery.error ? <Alert severity="error">{getErrorMessage(formQuery.error)}</Alert> : null}
        {submitMutation.error ? <Alert severity="error">{getErrorMessage(submitMutation.error)}</Alert> : null}
        {submitted ? <Alert severity="success">Resposta enviada.</Alert> : null}

        {form ? (
          <>
            <Stack sx={formPagesStyles.publicHeader}>
              <Typography variant="h4">{form.title}</Typography>
              {form.description ? <Typography color="text.secondary">{form.description}</Typography> : null}
            </Stack>

            <Stack sx={formPagesStyles.privacyNotice}>
              <Typography variant="subtitle2">Aviso de privacidade</Typography>
              <Typography variant="body2">Finalidade: {form.privacyPurpose ?? 'Não informada'}</Typography>
              <Typography variant="body2">Retenção: {form.retentionPolicy ?? 'Não informada'}</Typography>
              <Typography variant="body2">Contato do controlador: {form.controllerEmail ?? 'Não informado'}</Typography>
              <Typography variant="caption" color="text.secondary">
                Saiba mais na{' '}
                <Link component={RouterLink} to="/privacidade">
                  Política de Privacidade
                </Link>
                .
              </Typography>
            </Stack>

            <Stack sx={formPagesStyles.fieldStack}>
              {form.fields.map((field) => renderField(field, answers, updateAnswer))}
            </Stack>

            <FormControlLabel
              sx={formPagesStyles.privacyAcknowledgement}
              control={
                <Checkbox
                  checked={privacyAcknowledged}
                  required
                  onChange={(_, checked) => setPrivacyAcknowledged(checked)}
                />
              }
              label="Li e entendi o aviso de privacidade deste formulário."
            />

            <Button
              type="submit"
              variant="contained"
              disabled={submitMutation.isPending || !privacyAcknowledged}
              sx={formPagesStyles.submitButton}
            >
              Enviar resposta
            </Button>
          </>
        ) : null}
      </Stack>
    </Paper>
  );
}

function renderField(
  field: FormField,
  answers: AnswerState,
  updateAnswer: (fieldId: string, value: string | boolean) => void,
) {
  if (field.type === 'checkbox') {
    return (
      <FormControlLabel
        key={field.id}
        control={
          <Switch
            checked={Boolean(answers[field.id])}
            required={field.required}
            onChange={(_, checked) => updateAnswer(field.id, checked)}
          />
        }
        label={`${field.label}${field.required ? ' *' : ''}`}
      />
    );
  }

  if (field.type === 'select') {
    return (
      <TextField
        key={field.id}
        select
        label={field.label}
        required={field.required}
        value={String(answers[field.id] ?? '')}
        onChange={(event) => updateAnswer(field.id, event.target.value)}
      >
        <MenuItem value="">Selecione uma opção</MenuItem>
        {field.options.map((option) => (
          <MenuItem key={option} value={option}>
            {option}
          </MenuItem>
        ))}
      </TextField>
    );
  }

  return (
    <TextField
      key={field.id}
      label={field.label}
      type={inputTypeForField(field)}
      value={String(answers[field.id] ?? '')}
      onChange={(event) => updateAnswer(field.id, event.target.value)}
      placeholder={field.placeholder ?? undefined}
      required={field.required}
      multiline={field.type === 'textarea'}
      minRows={field.type === 'textarea' ? 4 : undefined}
      inputProps={inputPropsForField(field)}
    />
  );
}

function inputTypeForField(field: FormField) {
  if (field.type === 'number' || field.type === 'email') {
    return field.type;
  }

  if (field.type === 'phone') {
    return 'tel';
  }

  return 'text';
}

function inputPropsForField(field: FormField) {
  if (field.type === 'number') {
    return {
      inputMode: 'decimal',
    } as const;
  }

  if (field.type === 'phone') {
    return {
      inputMode: 'numeric',
      maxLength: 12,
      pattern: '[0-9]{12}',
      title: 'Informe 12 dígitos: 3 do DDD e 9 do telefone.',
    } as const;
  }

  return undefined;
}
