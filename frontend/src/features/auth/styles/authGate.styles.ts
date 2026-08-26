import type { SxProps, Theme } from '@mui/material/styles';

export const authGateStyles = {
  panel: {
    maxWidth: 420,
    mx: 'auto',
    p: {
      xs: 3,
      sm: 4,
    },
  },
  content: {
    alignItems: 'center',
    gap: 2,
    textAlign: 'center',
  },
} satisfies Record<string, SxProps<Theme>>;
