import { useState } from 'react';
import { Helmet } from 'react-helmet-async';

import { useInterval } from 'usehooks-ts';

import BlockIcon from '@mui/icons-material/Block';
import ErrorIcon from '@mui/icons-material/Error';
import PendingIcon from '@mui/icons-material/Pending';
import UndoIcon from '@mui/icons-material/Undo';
import { Grid, IconButton, Stack } from '@mui/material';

import { Container } from '@dashotv/components';
import { useQueryClient } from '@tanstack/react-query';

import { JobsList, JobsStats, deleteJob, patchJob, useJobsQuery } from 'components/jobs';

const Recent = () => {
  // limit, skip, queries, etc
  const [page] = useState(1);
  const [status, setStatus] = useState('');
  const { data } = useJobsQuery(page, status);
  const queryClient = useQueryClient();

  const handleCancel = (id: string) => {
    console.log('cancel', id);
    deleteJob(id, false);
  };

  const handleDelete = (id: string) => {
    console.log('delete', id);
    deleteJob(id, true);
  };

  const handleRequeue = (id: string) => {
    console.log('delete', id);
    patchJob(id);
  };

  useInterval(() => {
    queryClient.invalidateQueries({ queryKey: ['jobs'] });
  }, 5000);

  return (
    <>
      <Helmet>
        <title>Minion - Jobs</title>
        <meta name="description" content="runic" />
      </Helmet>

      <Container>
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
            </Stack>
          </Grid>
          <Grid item xs={12} md={6} justifyContent="end">
            {data?.stats ? <JobsStats stats={data.stats} setStatus={setStatus} /> : null}
          </Grid>
        </Grid>
      </Container>
      <Container>
        <JobsList {...{ data, handleCancel, handleDelete, handleRequeue }} />
      </Container>
    </>
  );
};

export default Recent;
