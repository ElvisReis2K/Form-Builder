import { alpha, createTheme } from '@mui/material/styles';

const colors = {
  background: '#f5f7fb',
  surface: '#ffffff',
  border: '#dfe6ef',
  textPrimary: '#172033',
  textSecondary: '#627187',
  primary: '#345d96',
  primaryDark: '#284878',
  secondary: '#2d846f',
  success: '#278264',
  error: '#b9474a',
  shadow: '#172033',
};

export const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: colors.primary,
      dark: colors.primaryDark,
      contrastText: colors.surface,
    },
    secondary: {
      main: colors.secondary,
    },
    success: {
      main: colors.success,
    },
    error: {
      main: colors.error,
    },
    background: {
      default: colors.background,
      paper: colors.surface,
    },
    divider: colors.border,
    text: {
      primary: colors.textPrimary,
      secondary: colors.textSecondary,
    },
  },
  typography: {
    fontFamily: '"Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
    allVariants: {
      letterSpacing: 0,
    },
    h4: {
      fontSize: '2rem',
      fontWeight: 750,
      lineHeight: 1.18,
    },
    h5: {
      fontSize: '1.45rem',
      fontWeight: 750,
      lineHeight: 1.22,
    },
    h6: {
      fontWeight: 700,
    },
    subtitle1: {
      fontWeight: 700,
    },
    button: {
      fontWeight: 700,
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          backgroundColor: colors.background,
          color: colors.textPrimary,
          textRendering: 'optimizeLegibility',
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
          border: `1px solid ${colors.border}`,
          boxShadow: `0 18px 45px ${alpha(colors.shadow, 0.07)}`,
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          boxShadow: 'none',
          minHeight: 38,
          textTransform: 'none',
          transition: 'background-color 160ms ease, border-color 160ms ease, color 160ms ease, transform 160ms ease',
          '&:hover': {
            boxShadow: 'none',
            transform: 'translateY(-1px)',
          },
          '&.Mui-disabled': {
            transform: 'none',
          },
        },
        contained: {
          boxShadow: `0 10px 22px ${alpha(colors.primary, 0.18)}`,
          '&:hover': {
            boxShadow: `0 12px 24px ${alpha(colors.primary, 0.22)}`,
          },
        },
        outlined: {
          backgroundColor: colors.surface,
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          fontWeight: 700,
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: colors.surface,
          borderRadius: 8,
          transition: 'background-color 160ms ease, box-shadow 160ms ease',
          '&.Mui-focused': {
            boxShadow: `0 0 0 3px ${alpha(colors.primary, 0.12)}`,
          },
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        root: {
          borderBottomColor: colors.border,
        },
        head: {
          color: colors.textSecondary,
          fontWeight: 700,
        },
      },
    },
  },
});
