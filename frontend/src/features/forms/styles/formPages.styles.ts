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
    alignItems: {
      xs: 'stretch',
      sm: 'center',
    },
  },
  titleBlock: {
    gap: 0.5,
    minWidth: 0,
  },
  workspace: {
    display: 'grid',
    gridTemplateColumns: {
      xs: '1fr',
      md: '320px minmax(0, 1fr)',
    },
    gap: 3,
    alignItems: 'start',
  },
  sidebarPanel: {
    overflow: 'hidden',
  },
  sidebarHeader: {
    p: 2,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 2,
  },
  formsList: {
    maxHeight: {
      xs: 280,
      md: 620,
    },
    overflow: 'auto',
  },
  editorPanel: {
    p: 3,
  },
  editorStack: {
    gap: 3,
  },
  editorHeader: {
    flexDirection: {
      xs: 'column',
      sm: 'row',
    },
    justifyContent: 'space-between',
    gap: 2,
    alignItems: {
      xs: 'stretch',
      sm: 'center',
    },
  },
  formGrid: {
    display: 'grid',
    gridTemplateColumns: '1fr',
    gap: 2,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 2,
  },
  fieldStack: {
    gap: 2,
  },
  fieldEditor: {
    gap: 2,
    border: 1,
    borderColor: 'divider',
    borderRadius: 1,
    p: 2,
    bgcolor: 'background.default',
  },
  fieldHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 2,
  },
  fieldGrid: {
    display: 'grid',
    gridTemplateColumns: {
      xs: '1fr',
      md: '180px minmax(0, 1fr) minmax(0, 1fr) 140px',
    },
    gap: 2,
    alignItems: 'center',
  },
  actionBar: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 1,
  },
  publishButton: {
    alignSelf: 'flex-start',
  },
  publicPanel: {
    maxWidth: 640,
    mx: 'auto',
    p: 3,
  },
  publicHeader: {
    gap: 0.5,
  },
  loadingBar: {
    borderRadius: 1,
  },
} satisfies Record<string, SxProps<Theme>>;
