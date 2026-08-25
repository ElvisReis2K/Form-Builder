import type { SxProps, Theme } from '@mui/material/styles';

export const formPagesStyles = {
  pageStack: {
    gap: {
      xs: 2.5,
      md: 3.5,
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
    minWidth: 0,
  },
  workspace: {
    display: 'grid',
    gridTemplateColumns: {
      xs: '1fr',
      md: '300px minmax(0, 1fr)',
      lg: '320px minmax(0, 1fr)',
    },
    gap: {
      xs: 2,
      md: 3,
    },
    alignItems: 'start',
  },
  sidebarPanel: {
    overflow: 'hidden',
    position: {
      md: 'sticky',
    },
    top: {
      md: 96,
    },
  },
  sidebarHeader: {
    p: 2.25,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 2,
  },
  formsList: {
    p: 1,
    maxHeight: {
      xs: 280,
      md: 650,
    },
    overflow: 'auto',
  },
  sidebarEmptyState: {
    minHeight: 132,
    alignItems: 'center',
    justifyContent: 'center',
    textAlign: 'center',
    px: 2,
  },
  formListItem: {
    mb: 0.75,
    gap: 1.5,
    '&:last-of-type': {
      mb: 0,
    },
    '&.Mui-selected': {
      bgcolor: 'primary.main',
      color: 'primary.contrastText',
      '& .MuiListItemText-secondary': {
        color: 'rgba(255,255,255,0.72)',
      },
      '&:hover': {
        bgcolor: 'primary.dark',
      },
    },
  },
  editorPanel: {
    p: {
      xs: 2.25,
      md: 3,
    },
  },
  editorStack: {
    gap: {
      xs: 2.5,
      md: 3,
    },
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
    gap: 2.25,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 2,
  },
  fieldStack: {
    gap: 1.5,
  },
  fieldEditor: {
    gap: 2,
    border: 1,
    borderColor: 'divider',
    borderRadius: 1,
    p: {
      xs: 1.75,
      md: 2,
    },
    bgcolor: '#f9fbfd',
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
    pt: 0.5,
  },
  publishButton: {
    alignSelf: 'flex-start',
  },
  publicPanel: {
    maxWidth: 700,
    mx: 'auto',
    mt: {
      xs: 2,
      md: 4,
    },
    p: {
      xs: 2.5,
      sm: 4,
    },
  },
  publicHeader: {
    gap: 0.75,
    pb: 1,
  },
  submitButton: {
    alignSelf: {
      xs: 'stretch',
      sm: 'flex-start',
    },
  },
  loadingBar: {
    borderRadius: 1,
  },
} satisfies Record<string, SxProps<Theme>>;
