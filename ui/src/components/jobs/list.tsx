import { useState } from 'react';

import AccessTimeIcon from '@mui/icons-material/AccessTime';
import BlockIcon from '@mui/icons-material/Block';
import CachedIcon from '@mui/icons-material/Cached';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import PendingIcon from '@mui/icons-material/Pending';
import { Link, Paper, Stack, Typography } from '@mui/material';

import { ButtonMap, ButtonMapButton, Chrono, Row } from 'components/common';

import { Job, JobsDialog } from '.';

const stringToColor = (value: string) => {
  switch (value) {
    case 'flame':
      return '#e57373';
    case 'tower':
      return '#4caf50';
    case 'runic':
      return '#2196f3';
    case 'rift':
      return '#ff9800';
    case 'scry':
      return '#9c27b0';
    case 'minion':
      // return "#ff5722";
      return '#FFFFFF';
    default:
      return '#888888';
  }
};
// const stringToHSL = (value: string) => {
//   let hash = 0;
//   for (let i = 0; i < value.length; i++) {
//     hash = value.charCodeAt(i) + ((hash << 5) - hash);
//   }
//
//   return `hsl(${hash % 360}, 100%, 65%)`;
// };

export interface JobsListProps {
  data?: any;
  handleCancel: (id: string) => void;
}

export function JobsList({ data, handleCancel }: JobsListProps) {
  const [selected, setSelected] = useState<Job | null>(null);
  const open = (job: Job) => {
    setSelected(job);
  };
  const close = () => {
    setSelected(null);
  };

  if (!data?.results || data?.results?.length === 0) {
    return (
      <Paper elevation={0}>
        <Typography color="gray" variant="caption">
          No jobs
        </Typography>
      </Paper>
    );
  }

  return (
    <Paper elevation={0}>
      {data?.results.map(job => {
        return (
          <Link key={job.id} href="#" onClick={() => open(job)}>
            <JobRow {...{ job, handleCancel }} />
          </Link>
        );
      })}

      {selected && <JobsDialog job={selected} close={close} />}
    </Paper>
  );
}

// const JobsListSection = ({
//   jobs,
//   status,
//   open,
//   handleCancel,
// }: {
//   jobs?: Job[];
//   status: string;
//   open: (job: Job) => void;
//   handleCancel: (id: string) => void;
// }) => {
//   return (
//     <Paper elevation={0} sx={{ minHeight: "75px", mb: 2 }}>
//       <Stack direction="row" spacing={1} alignItems="center">
//         <Icon status={status} />
//         <Typography color="primary" fontWeight="bolder">
//           {status.charAt(0).toUpperCase() + status.slice(1)}
//         </Typography>
//       </Stack>
//       {(!jobs || jobs?.length === 0) && (
//         <Typography color="gray" variant="caption">
//           No jobs
//         </Typography>
//       )}
//       {jobs?.map((job) => {
//         return (
//           <Link key={job.id} href="#" onClick={() => open(job)}>
//             <JobRow {...{ job, handleCancel }} />
//           </Link>
//         );
//       })}
//     </Paper>
//   );
// };
const Icon = ({ status }: { status: string }) => {
  switch (status) {
    case 'failed':
      return <ErrorIcon fontSize="small" color="error" />;
    case 'finished':
      return <CheckCircleIcon fontSize="small" color="success" />;
    case 'queued':
      return <AccessTimeIcon fontSize="small" color="secondary" />;
    case 'running':
      return <CachedIcon fontSize="small" color="primary" />;
    case 'cancelled':
      return <BlockIcon fontSize="small" color="warning" />;
    default:
      return <PendingIcon fontSize="small" color="disabled" />;
  }
};
const Error = ({ error }: { error?: string }) => {
  if (!error) return null;
  return (
    <Typography variant="caption" color="error" minWidth="0" /*width={{ xs: '100%', md: 'auto' }}*/ noWrap>
      {error}
    </Typography>
  );
};
const Title = ({ args }: { args: string }) => {
  if (args === '{}') return null;
  const parsed = JSON.parse(args);
  if (!parsed || !parsed.title) return null;
  return (
    <Typography variant="caption" color="gray" minWidth="0" /*width={{ xs: '100%', md: 'auto' }}*/ noWrap>
      {parsed.title}
    </Typography>
  );
};
const ErrorOrTitle = ({ error, args }: { error?: string; args: string }) => {
  if (error) return <Error error={error} />;
  return <Title args={args} />;
};

const Client = ({ client }: { client: string }) => {
  // const names: string[] = ["flame", "tower", "runic", "rift", "scry", "minion"];
  // const i: number = Math.floor(Math.random() * names.length);
  if (!client) client = 'unknown';
  return (
    <Typography sx={{ color: stringToColor(client) }} noWrap variant="button" color="primary.dark">
      {client || 'unknown'}
    </Typography>
  );
};

export function JobRow({
  job: { id, client, kind, queue, status, args, attempts },
  handleCancel,
}: {
  job: Job;
  handleCancel: (id: string) => void;
}) {
  const { started_at, duration, error } = (attempts && attempts.length > 0 && attempts[attempts.length - 1]) || {};

  const buttons: ButtonMapButton[] = [
    {
      Icon: BlockIcon,
      color: 'warning',
      click: e => {
        e.preventDefault();
        e.stopPropagation();
        handleCancel(id);
      },
      title: 'cancel',
    },
  ];
  return (
    <Row key={id}>
      <Stack width="100%" direction={{ xs: 'column', md: 'row' }} alignItems="center" justifyContent="space-between">
        <Stack
          width="100%"
          maxWidth="650px"
          direction={{ xs: 'column', md: 'row' }}
          spacing={1}
          alignItems="center"
          justifyContent="start"
        >
          <Stack
            width={{ xs: '100%', md: 'auto' }}
            direction="row"
            spacing={1}
            alignItems="center"
            justifyContent="start"
          >
            <Icon status={status} />
            <Client client={client} />
            <Typography minWidth="0" color={status === 'failed' ? 'error' : 'primary'} noWrap>
              {kind}
            </Typography>
          </Stack>
          <ErrorOrTitle {...{ error, args }} />
        </Stack>
        <Stack
          minWidth="300px"
          width={{ xs: '100%', md: 'auto' }}
          direction="row"
          spacing={1}
          alignItems="center"
          justifyContent="end"
        >
          <Typography noWrap variant="button" color="primary.dark">
            {queue}
          </Typography>
          <Typography noWrap fontWeight="bolder" color="action">
            {duration ? `${duration.toFixed(1)}s` : ''}
          </Typography>
          <Typography variant="subtitle2" color="gray" noWrap>
            {started_at ? <Chrono fromNow>{started_at.toString()}</Chrono> : ''}
          </Typography>
          <ButtonMap size="small" buttons={buttons} />
        </Stack>
      </Stack>
    </Row>
  );
}
