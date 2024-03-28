import React from 'react';

import CloseIcon from '@mui/icons-material/Close';
import {
  Alert,
  AlertTitle,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  Paper,
  Stack,
  Typography,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import useMediaQuery from '@mui/material/useMediaQuery';

import { Chrono } from 'components/common';

import { Job } from './types';

const Args = ({ args }: { args: string }) => {
  return (
    <Paper elevation={1} sx={{ p: '8px' }} className="args">
      <pre
        style={{
          fontSize: '12px',
          paddingTop: '1px',
          marginTop: '1px',
          marginBottom: '1px',
        }}
      >
        {JSON.stringify(JSON.parse(args), null, 2)}
      </pre>
    </Paper>
  );
};

const Attempts = ({ attempts }: { attempts: Job['attempts'] }) => {
  return attempts && attempts.map((attempt, i) => <Attempt key={i} attempt={attempt} />);
};

const Attempt = ({ attempt: { status, started_at, duration, error, stacktrace } }: { attempt: Job['attempts'][0] }) => {
  return (
    <Paper elevation={1} className="attempt" sx={{ p: '8px', mb: 1 }}>
      <Stack direction="row" spacing={1} alignItems="center" justifyContent="start" width="100%">
        <Typography variant="subtitle1" color={status === 'failed' ? 'error' : 'action'}>
          {status}
        </Typography>
        <Stack width="100%" direction="row" spacing={1} alignItems="center" justifyContent="end">
          <Typography variant="caption" color="action">
            {duration ? `${duration.toFixed(1)}s` : ''}
          </Typography>
          <Typography minWidth="100px" variant="subtitle2" color="gray" noWrap>
            {started_at ? <Chrono fromNow>{started_at.toString()}</Chrono> : 'Pending'}
          </Typography>
        </Stack>
      </Stack>
      <Stacktrace {...{ error, stacktrace }} />
    </Paper>
  );
};

const Stacktrace = ({ error, stacktrace }: { error: string; stacktrace: string[] }) => {
  if (!error) return null;
  return (
    <Alert severity="error">
      <AlertTitle>{error}</AlertTitle>
      <div role="alert">
        <Stack direction="column">{stacktrace && stacktrace.map((trace, i) => <Trace key={i} trace={trace} />)}</Stack>
      </div>
    </Alert>
  );
};
const Trace = ({ trace }: { trace: Job['attempts'][0]['stacktrace'][0] }) => {
  return (
    <Typography variant="subtitle2" color="action">
      {trace}
    </Typography>
  );
};

export const JobsDialog = ({
  close,
  job: { id, kind, queue, status, args, attempts },
}: {
  close: () => void;
  job: Job;
}) => {
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down('md'));
  const [open, setOpen] = React.useState(true);
  const handleClose = () => {
    setOpen(false);
    close();
  };
  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth fullScreen={fullScreen}>
      <DialogTitle>
        <Stack direction="row" spacing={3} alignItems="center" justifyContent="start" width="100%">
          <Typography fontWeight="bolder" color="primary">
            {kind}
          </Typography>
        </Stack>
        <IconButton
          aria-label="close"
          onClick={handleClose}
          sx={{
            position: 'absolute',
            right: 8,
            top: 8,
            color: theme => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent>
        <Stack width="100%" direction={{ xs: 'column', md: 'row' }} spacing={0} sx={{ mb: 1 }}>
          <Typography minWidth="125px" variant="subtitle1" color="textSecondary">
            ID
          </Typography>
          <Typography variant="button">{id}</Typography>
        </Stack>
        <Stack width="100%" direction={{ xs: 'column', md: 'row' }} spacing={0} sx={{ mb: 1 }}>
          <Typography minWidth="125px" variant="subtitle1" color="textSecondary">
            Status
          </Typography>
          <Typography variant="subtitle1" color={status === 'failed' ? 'error' : 'success'}>
            {status}
          </Typography>
        </Stack>
        <Stack width="100%" direction={{ xs: 'column', md: 'row' }} spacing={0} sx={{ mb: 1 }}>
          <Typography minWidth="125px" variant="subtitle1" color="textSecondary">
            Queue
          </Typography>
          <Typography variant="subtitle1">{queue}</Typography>
        </Stack>

        <Typography variant="subtitle1" color="textSecondary">
          Args {args === '{}' ? '(none)' : ''}
        </Typography>
        <Args {...{ args }} />
        <Typography variant="subtitle1" color="textSecondary">
          Attempts
        </Typography>
        <Attempts {...{ attempts }} />
      </DialogContent>
    </Dialog>
  );
};
