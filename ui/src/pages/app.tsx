import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider, createTheme } from '@mui/material/styles';

import { Container } from '@dashotv/components';
import { RoutingTabs, RoutingTabsRoute } from '@dashotv/components';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import Recent from 'pages/recent';
import Search from 'pages/search';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
  components: {
    MuiLink: {
      styleOverrides: {
        root: {
          textDecoration: 'none',
        },
      },
    },
  },
});

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 5,
      staleTime: 5 * 1000,
      throwOnError: true,
    },
  },
});

const App = ({ mount }: { mount: string }) => {
  const tabsMap: RoutingTabsRoute[] = [
    {
      label: 'Recent',
      to: '',
      element: <Recent />,
    },
    {
      label: 'Search',
      to: 'search',
      element: <Search />,
    },
  ];
  return (
    <ThemeProvider theme={darkTheme}>
      <QueryClientProvider client={queryClient}>
        <CssBaseline />
        <Container>
          <RoutingTabs data={tabsMap} mount={mount} />
        </Container>
      </QueryClientProvider>
    </ThemeProvider>
  );
};

export default App;
