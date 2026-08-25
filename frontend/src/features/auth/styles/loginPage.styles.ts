import type { SxProps, Theme } from '@mui/material/styles';

export const loginPageStyles = {
  panel: {
    maxWidth: 440,
    mx: 'auto',
    mt: {
      xs: 3,
      md: 7,
    },
    p: {
      xs: 3,
      sm: 4,
    },
  },
  form: {
    gap: 2.25,
  },
  header: {
    gap: 0.75,
    pb: 1,
    textAlign: 'center',
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
  actionButton: {
    minWidth: {
      sm: 132,
    },
  },
  googleButton: {
    mt: 0.5,
  },
} satisfies Record<string, SxProps<Theme>>;
