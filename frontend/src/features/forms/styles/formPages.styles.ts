import type { SxProps, Theme } from '@mui/material/styles';

export const formPagesStyles = {
  pageStack: {
    gap: 3,
  },
  header: {
    flexDirection: {
      xs: 'column',
      sm: 'row',
    },
    gap: 2,
    justifyContent: 'space-between',
  },
  editorPanel: {
    p: 3,
  },
  fieldStack: {
    gap: 2,
  },
  publishButton: {
    alignSelf: 'flex-start',
  },
  publicPanel: {
    maxWidth: 640,
    mx: 'auto',
    p: 3,
  },
} satisfies Record<string, SxProps<Theme>>;
