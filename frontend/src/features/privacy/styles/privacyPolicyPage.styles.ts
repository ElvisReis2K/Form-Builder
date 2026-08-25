import type { SxProps, Theme } from '@mui/material/styles';

export const privacyPolicyPageStyles = {
  panel: {
    maxWidth: 860,
    mx: 'auto',
    p: {
      xs: 2.5,
      sm: 4,
    },
  },
  stack: {
    gap: {
      xs: 2.5,
      md: 3,
    },
  },
  header: {
    gap: 0.75,
    pb: 1,
  },
  section: {
    gap: 0.75,
  },
} satisfies Record<string, SxProps<Theme>>;
