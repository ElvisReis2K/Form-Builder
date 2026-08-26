import type { SxProps, Theme } from '@mui/material/styles';

export const responsesPageStyles = {
  pageStack: {
    gap: {
      xs: 2.25,
      md: 3,
    },
  },
  header: {
    flexDirection: {
      xs: 'column',
      sm: 'row',
    },
    gap: 2,
    justifyContent: 'space-between',
    alignItems: {
      xs: 'stretch',
      sm: 'center',
    },
    pb: 1,
  },
  titleBlock: {
    gap: 0.75,
  },
  panel: {
    overflow: 'hidden',
    bgcolor: 'rgba(255, 255, 255, 0.92)',
  },
  tableContainer: {
    overflowX: 'auto',
  },
  table: {
    minWidth: 760,
  },
  tableHead: {
    bgcolor: 'rgba(232, 241, 243, 0.54)',
  },
  emptyState: {
    minHeight: 240,
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
