import type { SxProps, Theme } from '@mui/material/styles';

export const appStyles = {
  root: {
    minHeight: '100vh',
  },
  toolbar: {
    gap: 2,
  },
  brand: {
    flexGrow: 1,
  },
  nav: {
    flexDirection: 'row',
    gap: 1,
  },
  content: {
    py: 4,
  },
} satisfies Record<string, SxProps<Theme>>;
