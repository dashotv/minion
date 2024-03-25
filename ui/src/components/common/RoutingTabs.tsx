import { Link, Route, Routes, matchPath, useLocation } from 'react-router-dom';

import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';

export interface RoutingTabsRoute {
  label: string;
  element: JSX.Element;
  to: string;
}
export const RoutingTabs = ({ data, route }: { data: RoutingTabsRoute[]; route: string }) => {
  const routes = data.map(({ to }) => (to !== '' ? `${route}/${to}` : route));
  const routeMatch = useRouteMatch(routes);
  const currentTab = routeMatch;

  return (
    <>
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs variant="scrollable" value={currentTab}>
          {data.map(({ label, to }, i) => {
            return <RoutingTab key={i} index={i} label={label} to={to} />;
          })}
        </Tabs>
      </Box>
      <Paper elevation={0} sx={{ pt: 1 }}>
        <Routes>
          {data.map(({ to, element }, i) => {
            return <Route key={i} path={to} element={element} />;
          })}
        </Routes>
      </Paper>
    </>
  );
};

export type RoutingTabProps = {
  key: number;
  index: number;
  to: string;
  label: string;
};
export function RoutingTab({ key, index, to, label }: RoutingTabProps) {
  return (
    <Tab
      key={key}
      to={to}
      label={label}
      id={`routing-tab-${index}`}
      aria-controls={`routing-tabpanel-${index}`}
      component={Link}
    />
  );
}

function useRouteMatch(patterns: readonly string[]) {
  const { pathname } = useLocation();

  for (let i = 0; i < patterns.length; i += 1) {
    if (matchPath(patterns[i], pathname) !== null) {
      return i;
    }
  }

  return null;
}
