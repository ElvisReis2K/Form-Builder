import type { SxProps, Theme } from '@mui/material/styles';

export const loginPageStyles = {
  panel: {
    maxWidth: 420,
    mx: 'auto',
    p: 3,
  },
  form: {
    gap: 2,
  },
} satisfies Record<string, SxProps<Theme>>;
