import type { SxProps, Theme } from '@mui/material/styles';

export const appStyles = {
  root: {
    minHeight: '100vh',
    background: 'linear-gradient(180deg, #f9fbf6 0%, #eff5f1 48%, #f7f9f4 100%)',
  },
  appBar: {
    position: 'sticky',
    top: 0,
    zIndex: (theme) => theme.zIndex.appBar,
    bgcolor: 'rgba(255, 255, 255, 0.82)',
    backdropFilter: 'blur(18px)',
    borderBottom: 1,
    borderColor: 'divider',
    boxShadow: '0 8px 24px rgba(20, 33, 31, 0.04)',
  },
  toolbar: {
    minHeight: {
      xs: 64,
      sm: 72,
    },
    gap: 2,
    flexWrap: 'wrap',
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
    flexWrap: 'wrap',
    justifyContent: 'flex-end',
  },
  navButton: {
    color: 'text.secondary',
    px: 1.5,
    minHeight: 34,
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
