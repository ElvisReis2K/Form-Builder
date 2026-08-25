import type { SxProps, Theme } from '@mui/material/styles';

export const responsesPageStyles = {
  pageStack: {
    gap: {
      xs: 2.5,
      md: 3.5,
    },
  },
  header: {
    gap: 0.75,
    pb: 1,
  },
  panel: {
    overflow: 'hidden',
  },
  tableContainer: {
    overflowX: 'auto',
  },
  table: {
    minWidth: 560,
  },
  tableHead: {
    bgcolor: '#f9fbfd',
  },
  emptyState: {
    minHeight: 260,
    alignItems: 'center',
    justifyContent: 'center',
    gap: 1,
    textAlign: 'center',
    p: 3,
  },
  loadingBar: {
    borderRadius: 1,
  },
} satisfies Record<string, SxProps<Theme>>;
