import { alpha, createTheme } from '@mui/material/styles';

const colors = {
  background: '#f7f9f4',
  surface: '#ffffff',
  surfaceSoft: '#fbfcf8',
  border: '#dde6df',
  textPrimary: '#18211f',
  textSecondary: '#5f6d68',
  primary: '#2f6073',
  primaryDark: '#24495a',
  primarySoft: '#e8f1f3',
  secondary: '#6f7f3f',
  success: '#2f7a62',
  error: '#b44f59',
  shadow: '#14211f',
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
          WebkitFontSmoothing: 'antialiased',
          MozOsxFontSmoothing: 'grayscale',
        },
        '::selection': {
          backgroundColor: alpha(colors.primary, 0.18),
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
          boxShadow: `0 16px 38px ${alpha(colors.shadow, 0.055)}, 0 2px 8px ${alpha(colors.shadow, 0.035)}`,
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
          transition:
            'background-color 180ms ease, border-color 180ms ease, box-shadow 180ms ease, color 180ms ease, transform 180ms cubic-bezier(0.2, 0.8, 0.2, 1)',
          willChange: 'transform',
          '&:hover': {
            boxShadow: 'none',
            transform: 'translateY(-2px)',
          },
          '&:active': {
            transform: 'translateY(0) scale(0.985)',
          },
          '&.Mui-focusVisible': {
            boxShadow: `0 0 0 3px ${alpha(colors.primary, 0.18)}`,
          },
          '&.Mui-disabled': {
            transform: 'none',
            boxShadow: 'none',
          },
          '@media (prefers-reduced-motion: reduce)': {
            transition: 'none',
            willChange: 'auto',
            '&:hover': {
              transform: 'none',
            },
            '&:active': {
              transform: 'none',
            },
          },
        },
        contained: {
          boxShadow: `0 10px 22px ${alpha(colors.primary, 0.16)}`,
          '&:hover': {
            boxShadow: `0 14px 26px ${alpha(colors.primary, 0.2)}`,
          },
        },
        outlined: {
          backgroundColor: alpha(colors.surface, 0.86),
          borderColor: colors.border,
          '&:hover': {
            backgroundColor: colors.surface,
            borderColor: alpha(colors.primary, 0.45),
          },
        },
        text: {
          '&:hover': {
            backgroundColor: alpha(colors.primary, 0.07),
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          fontWeight: 700,
          transition: 'background-color 160ms ease, color 160ms ease, border-color 160ms ease',
        },
      },
    },
    MuiAlert: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          border: `1px solid ${colors.border}`,
          boxShadow: 'none',
        },
      },
    },
    MuiCheckbox: {
      styleOverrides: {
        root: {
          transition: 'background-color 160ms ease, transform 160ms ease',
          '&:hover': {
            transform: 'scale(1.04)',
          },
          '@media (prefers-reduced-motion: reduce)': {
            transition: 'none',
            '&:hover': {
              transform: 'none',
            },
          },
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: colors.surface,
          borderRadius: 8,
          transition: 'background-color 170ms ease, box-shadow 170ms ease, transform 170ms ease',
          '&:hover': {
            backgroundColor: colors.surfaceSoft,
          },
          '&.Mui-focused': {
            boxShadow: `0 0 0 3px ${alpha(colors.primary, 0.12)}`,
          },
          '@media (prefers-reduced-motion: reduce)': {
            transition: 'none',
          },
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          transition: 'background-color 160ms ease, color 160ms ease, transform 160ms ease',
          '&:hover': {
            transform: 'translateX(2px)',
          },
          '@media (prefers-reduced-motion: reduce)': {
            transition: 'none',
            '&:hover': {
              transform: 'none',
            },
          },
        },
      },
    },
    MuiLinearProgress: {
      styleOverrides: {
        root: {
          height: 5,
          backgroundColor: colors.primarySoft,
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
    MuiTableRow: {
      styleOverrides: {
        root: {
          transition: 'background-color 150ms ease',
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        fullWidth: true,
      },
    },
  },
});
