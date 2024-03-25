import React from 'react';

import { Paper } from '@mui/material';

export const Row = ({ variant, children }: { variant?: string; children: React.ReactNode }) => {
  let color: string | undefined = undefined;
  let hover: string = '#242424';

  switch (variant) {
    case 'success':
      color = '#226622';
      hover = '#228822';
      break;
    case 'selected':
      color = '#222266';
      hover = '#222288';
      break;
    case 'transparent':
      color = 'transparent';
      break;
  }

  return (
    <Paper elevation={0} sx={{ p: 0.5, '&:hover': { backgroundColor: hover }, zIndex: 1000, backgroundColor: color }}>
      {children}
    </Paper>
  );
};
