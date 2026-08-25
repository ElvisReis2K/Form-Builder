import { Button, Paper, Stack, TextField, Typography } from '@mui/material';

import { formPagesStyles } from '../styles/formPages.styles';

export default function AdminHomePage() {
  return (
    <Stack sx={formPagesStyles.pageStack}>
      <Stack sx={formPagesStyles.header}>
        <Typography variant="h4">Forms</Typography>
        <Button variant="contained">New form</Button>
      </Stack>

      <Paper sx={formPagesStyles.editorPanel}>
        <Stack sx={formPagesStyles.fieldStack}>
          <TextField label="Form title" defaultValue="Customer feedback" />
          <TextField label="Description" defaultValue="Simple published form draft" multiline rows={3} />
          <Button variant="outlined" sx={formPagesStyles.publishButton}>
            Publish
          </Button>
        </Stack>
      </Paper>
    </Stack>
  );
}
