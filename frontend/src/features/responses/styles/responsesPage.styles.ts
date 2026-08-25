import type { SxProps, Theme } from '@mui/material/styles';

export const responsesPageStyles = {
  pageStack: {
    gap: 3,
  },
  header: {
    gap: 0.5,
  },
  panel: {
    p: 3,
  },
  tableContainer: {
    overflowX: 'auto',
  },
  emptyState: {
    minHeight: 220,
    alignItems: 'center',
    justifyContent: 'center',
    gap: 1,
    textAlign: 'center',
  },
  loadingBar: {
    borderRadius: 1,
  },
} satisfies Record<string, SxProps<Theme>>;
