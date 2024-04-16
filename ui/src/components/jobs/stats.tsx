import AccessTimeIcon from '@mui/icons-material/AccessTime';
import BlockIcon from '@mui/icons-material/Block';
import CachedIcon from '@mui/icons-material/Cached';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import PendingOutlinedIcon from '@mui/icons-material/PendingOutlined';
import { Button, Stack } from '@mui/material';

import { Pill } from '@dashotv/components';

import { Stats } from './types';

export const JobsStats = ({ stats, setStatus }: { stats: Stats; setStatus: (status: string) => void }) => {
  if (!stats) {
    return null;
  }
  return (
    <Stack direction="row" spacing={0} justifyContent="end">
      <Button title="pending" onClick={() => setStatus('pending')}>
        <Pill name="P" color="gray" value={stats.pending || 0} icon={<PendingOutlinedIcon fontSize="small" />} />
      </Button>
      <Button title="queued" onClick={() => setStatus('queued')}>
        <Pill name="Q" color="secondary" value={stats.queued || 0} icon={<AccessTimeIcon fontSize="small" />} />
      </Button>
      <Button title="running" onClick={() => setStatus('running')}>
        <Pill name="R" color="primary" value={stats.running || 0} icon={<CachedIcon fontSize="small" />} />
      </Button>
      <Button title="cancelled" onClick={() => setStatus('cancelled')}>
        <Pill name="C" color="warning" value={stats.cancelled} icon={<BlockIcon fontSize="small" />} />
      </Button>
      <Button title="failed" onClick={() => setStatus('failed')}>
        <Pill name="F" color="error" value={stats.failed} icon={<ErrorOutlineIcon fontSize="small" />} />
      </Button>
    </Stack>
  );
};
