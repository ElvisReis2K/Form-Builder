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
  header: {
    gap: 0.5,
  },
  actions: {
    flexDirection: {
      xs: 'column',
      sm: 'row',
    },
    gap: 1,
    alignItems: {
      xs: 'stretch',
      sm: 'center',
    },
  },
} satisfies Record<string, SxProps<Theme>>;
