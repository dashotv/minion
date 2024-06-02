import AccessTimeIcon from '@mui/icons-material/AccessTime';
import ArchiveIcon from '@mui/icons-material/Archive';
import BlockIcon from '@mui/icons-material/Block';
import CachedIcon from '@mui/icons-material/Cached';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import PendingOutlinedIcon from '@mui/icons-material/PendingOutlined';
import { Button, Stack } from '@mui/material';

import { Stats } from './types';

export const JobsStats = ({ stats, setStatus }: { stats: Stats; setStatus: (status: string) => void }) => {
  if (!stats) {
    return null;
  }
  return (
    <Stack direction="row" spacing={0} justifyContent="end">
      <Stat
        setStatus={setStatus}
        name="pending"
        color="primary"
        value={stats.pending || 0}
        icon={<PendingOutlinedIcon color="disabled" fontSize="small" />}
      />
      <Stat
        setStatus={setStatus}
        name="queued"
        color="secondary"
        value={stats.queued || 0}
        icon={<AccessTimeIcon color="secondary" fontSize="small" />}
      />
      <Stat
        setStatus={setStatus}
        name="running"
        color="primary"
        value={stats.running || 0}
        icon={<CachedIcon color="primary" fontSize="small" />}
      />
      <Stat
        setStatus={setStatus}
        name="cancelled"
        color="warning"
        value={stats.cancelled || 0}
        icon={<BlockIcon color="warning" fontSize="small" />}
      />
      <Stat
        setStatus={setStatus}
        name="failed"
        color="error"
        value={stats.failed || 0}
        icon={<ErrorOutlineIcon color="error" fontSize="small" />}
      />
      <Stat
        setStatus={setStatus}
        name="archived"
        color="primary"
        value={0}
        icon={<ArchiveIcon color="disabled" fontSize="small" />}
      />
      {/* <Button title="pending" onClick={() => setStatus('pending')}>
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
      <Button title="failed" onClick={() => setStatus('archived')}>
        <Pill name="A" color="gray" value={stats.archived} icon={<ArchiveIcon fontSize="small" />} />
      </Button> */}
    </Stack>
  );
};

const Stat = ({ name, value, icon, color, setStatus }) => {
  return (
    <Stack direction="row" spacing={0} alignItems="center" title={name}>
      <Button color={color} startIcon={icon} size="small" onClick={() => setStatus(name)}>
        {value}
      </Button>
    </Stack>
  );
};
