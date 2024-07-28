import { useState } from 'react';

import { useInterval } from 'usehooks-ts';

import { Box } from '@mui/material';
import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider, createTheme } from '@mui/material/styles';

import { Container, FilterSelect, Option, RoutingTabs, RoutingTabsRoute } from '@dashotv/components';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { Jobs, useJobsStatsQuery } from 'components/jobs';

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
  return (
    <ThemeProvider theme={darkTheme}>
      <QueryClientProvider client={queryClient}>
        <CssBaseline />
        <Container>
          <Routes mount={mount} />
        </Container>
      </QueryClientProvider>
    </ThemeProvider>
  );
};

const clients: Option[] = [
  { label: 'all', value: '' },
  { label: 'tower', value: 'tower' },
  { label: 'flame', value: 'flame' },
  { label: 'runic', value: 'runic' },
  { label: 'rift', value: 'rift' },
  { label: 'scry', value: 'scry' },
  { label: 'arcane', value: 'arcane' },
];

export const Routes = ({ mount }: { mount: string }) => {
  const [client, setClient] = useState('');
  const { data: stats } = useJobsStatsQuery();

  const choose = (choice: string) => {
    setClient(choice);
  };

  useInterval(() => {
    console.log('invalidate');
    queryClient.invalidateQueries({ queryKey: ['jobs'] });
  }, 5000);

  const tabsMap: RoutingTabsRoute[] = [
    {
      label: 'Recent',
      to: '',
      element: <Jobs client={client} status="" />,
    },
    {
      label: `Pending ${stats?.pending || 0}`,
      to: 'pending',
      element: <Jobs client={client} status="pending" />,
    },
    {
      label: `Queued ${stats?.queued || 0}`,
      to: 'queued',
      element: <Jobs client={client} status="queued" />,
    },
    {
      label: `Running ${stats?.running || 0}`,
      to: 'running',
      element: <Jobs client={client} status="running" />,
    },
    {
      label: `Cancelled ${stats?.cancelled || 0}`,
      to: 'cancelled',
      element: <Jobs client={client} status="cancelled" />,
    },
    {
      label: `Failed ${stats?.failed || 0}`,
      to: 'failed',
      element: <Jobs client={client} status="failed" />,
    },
    {
      label: 'Archived',
      to: 'archived',
      element: <Jobs client={client} status="archived" />,
    },
  ];

  return (
    <Box sx={{ pt: 1, pr: 1, pl: 1, mb: 2, position: 'relative' }}>
      <Box sx={{ width: '125px', position: 'absolute', right: '15px', top: '8px', zIndex: 1000 }}>
        <FilterSelect name="client" selected={client} choose={choose} choices={clients} />
      </Box>
      <RoutingTabs data={tabsMap} mount={mount} />
    </Box>
  );
};

export default App;
