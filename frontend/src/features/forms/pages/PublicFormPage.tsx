import {
  Alert,
  Button,
  FormControlLabel,
  LinearProgress,
  MenuItem,
  Paper,
  Stack,
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';

import { getErrorMessage } from '../../../lib/api';
import { getPublishedForm } from '../api/formsApi';
import { formPagesStyles } from '../styles/formPages.styles';
import type { FormField } from '../types';
import { fieldTypeLabels } from '../utils/formDraft';

export default function PublicFormPage() {
  const { slug } = useParams();

  const formQuery = useQuery({
    queryKey: ['public-form', slug],
    queryFn: () => getPublishedForm(slug ?? ''),
    enabled: Boolean(slug),
  });

  const form = formQuery.data;

  return (
    <Paper sx={formPagesStyles.publicPanel}>
      <Stack sx={formPagesStyles.fieldStack}>
        {formQuery.isFetching ? <LinearProgress sx={formPagesStyles.loadingBar} /> : null}
        {formQuery.error ? <Alert severity="error">{getErrorMessage(formQuery.error)}</Alert> : null}

        {form ? (
          <>
            <Stack sx={formPagesStyles.publicHeader}>
              <Typography variant="h4">{form.title}</Typography>
              {form.description ? <Typography color="text.secondary">{form.description}</Typography> : null}
            </Stack>

            <Stack sx={formPagesStyles.fieldStack}>{form.fields.map((field) => renderField(field))}</Stack>

            <Button variant="contained" disabled>
              Submit response
            </Button>
          </>
        ) : null}
      </Stack>
    </Paper>
  );
}

function renderField(field: FormField) {
  if (field.type === 'checkbox') {
    return (
      <FormControlLabel
        key={field.id}
        control={<Switch required={field.required} />}
        label={`${field.label}${field.required ? ' *' : ''}`}
      />
    );
  }

  if (field.type === 'select') {
    return (
      <TextField key={field.id} select label={field.label} required={field.required} defaultValue="">
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
      type={field.type === 'number' || field.type === 'email' ? field.type : 'text'}
      placeholder={field.placeholder ?? undefined}
      required={field.required}
      multiline={field.type === 'textarea'}
      minRows={field.type === 'textarea' ? 4 : undefined}
      InputProps={{
        inputProps: {
          'aria-label': fieldTypeLabels[field.type],
        },
      }}
    />
  );
}
