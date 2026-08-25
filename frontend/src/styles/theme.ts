import { createTheme } from '@mui/material';

export const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#2155d9',
    },
    secondary: {
      main: '#1f7a5f',
    },
    background: {
      default: '#f7f8fb',
    },
  },
  shape: {
    borderRadius: 8,
  },
});
