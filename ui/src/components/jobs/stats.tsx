import { Button, Stack } from '@mui/material';

import { Pill } from 'components/common';

import { Stats } from './types';

export const JobsStats = ({ stats, setStatus }: { stats: Stats; setStatus: (status: string) => void }) => {
  if (!stats) {
    return null;
  }
  return (
    <Stack direction="row" spacing={0} justifyContent="end">
      <Button onClick={() => setStatus('pending')}>
        <Pill name="Pending" color="gray" value={stats.pending || 0} />
      </Button>
      <Button onClick={() => setStatus('queued')}>
        <Pill name="Queued" color="secondary" value={stats.queued || 0} />
      </Button>
      <Button onClick={() => setStatus('running')}>
        <Pill name="Running" color="primary" value={stats.running || 0} />
      </Button>
      <Button onClick={() => setStatus('cancelled')}>
        <Pill name="Cancelled" color="warning" value={stats.cancelled} />
      </Button>
      <Button onClick={() => setStatus('failed')}>
        <Pill name="Failed" color="error" value={stats.failed} />
      </Button>
    </Stack>
  );
};
