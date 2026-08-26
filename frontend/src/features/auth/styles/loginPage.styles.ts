import type { SxProps, Theme } from '@mui/material/styles';

export const loginPageStyles = {
  panel: {
    maxWidth: 460,
    mx: 'auto',
    mt: {
      xs: 3,
      md: 6,
    },
    p: {
      xs: 3,
      sm: 4,
    },
    borderTop: 3,
    borderTopColor: 'secondary.main',
  },
  form: {
    gap: 2,
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
    alignItems: 'center',
    justifyContent: 'center',
  },
  actionButton: {
    minWidth: {
      sm: 132,
    },
    width: {
      xs: '100%',
      sm: 132,
    },
    maxWidth: {
      xs: 220,
      sm: 'none',
    },
  },
  googleButton: {
    mt: 0.25,
  },
  privacyText: {
    textAlign: 'center',
  },
} satisfies Record<string, SxProps<Theme>>;
