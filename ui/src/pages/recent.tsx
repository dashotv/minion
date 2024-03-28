import { useState } from 'react';
import { Helmet } from 'react-helmet-async';

import BlockIcon from '@mui/icons-material/Block';
import ErrorIcon from '@mui/icons-material/Error';
import PendingIcon from '@mui/icons-material/Pending';
import UndoIcon from '@mui/icons-material/Undo';
import { Grid, IconButton, Stack } from '@mui/material';

import { JobsList, JobsStats, deleteJob, useJobsQuery } from 'components/jobs';

const Recent = () => {
  // limit, skip, queries, etc
  const [page] = useState(1);
  const [status, setStatus] = useState('');
  const { data } = useJobsQuery(page, status);

  const handleCancel = (id: string) => {
    console.log('cancel', id);
    deleteJob(id, false);
  };

  const handleDelete = (id: string) => {
    console.log('delete', id);
    deleteJob(id, true);
  };

  return (
    <>
      <Helmet>
        <title>Minion - Jobs</title>
        <meta name="description" content="runic" />
      </Helmet>

      <Grid container spacing={0} sx={{ mb: 2 }}>
        <Grid item xs={12} md={6}>
          <Stack direction="row" spacing={0} alignItems="center">
            <IconButton title="Cancel Pending" onClick={() => handleCancel('pending')}>
              <PendingIcon color="disabled" />
            </IconButton>
            <IconButton title="Delete Cancelled" onClick={() => handleDelete('cancelled')}>
              <BlockIcon color="warning" />
            </IconButton>
            <IconButton title="Delete Failed" onClick={() => handleDelete('failed')}>
              <ErrorIcon color="error" />
            </IconButton>
            <IconButton title="Show All" onClick={() => setStatus('')}>
              <UndoIcon color="primary" />
            </IconButton>
            {data?.stats ? <JobsStats stats={data.stats} setStatus={setStatus} /> : null}
          </Stack>
        </Grid>
        <Grid item xs={12} md={6} justifyContent="end"></Grid>
      </Grid>

      <JobsList {...{ data, status, handleCancel, handleDelete }} />
    </>
  );
};

export default Recent;
