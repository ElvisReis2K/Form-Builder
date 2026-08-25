import type { SxProps, Theme } from '@mui/material/styles';

export const appStyles = {
  root: {
    minHeight: '100vh',
    bgcolor: 'background.default',
  },
  appBar: {
    position: 'sticky',
    top: 0,
    zIndex: (theme) => theme.zIndex.appBar,
    bgcolor: 'rgba(255, 255, 255, 0.88)',
    backdropFilter: 'blur(14px)',
    borderBottom: 1,
    borderColor: 'divider',
  },
  toolbar: {
    minHeight: {
      xs: 64,
      sm: 72,
    },
    gap: 2,
  },
  brand: {
    flexGrow: 1,
    color: 'text.primary',
    fontWeight: 800,
    lineHeight: 1.15,
  },
  nav: {
    flexDirection: 'row',
    gap: 0.75,
  },
  navButton: {
    color: 'text.secondary',
    px: 1.5,
    '&:hover': {
      bgcolor: 'action.hover',
      color: 'primary.main',
    },
  },
  content: {
    py: {
      xs: 3,
      md: 5,
    },
  },
} satisfies Record<string, SxProps<Theme>>;
