import { ThreeDots } from 'react-loader-spinner';

import Box from '@mui/material/Box';

export const LoadingIndicator = () => (
  <Box
    sx={{
      position: 'fixed',
      bottom: 0,
      right: 0,
      left: 0,
      height: '60px',
      width: '100%',
      overflow: 'hidden',
      justifyContent: 'center',
      alignContent: 'center',
      alignItems: 'center',
      display: 'flex',
      zIndex: 1000,
      backgroundImage: 'linear-gradient(to top, rgba(0,0,0,1), rgba(0,0,0,0))',
    }}
  >
    <ThreeDots width="60" height="30" radius="3" color="white" ariaLabel="three-dots-loading" visible={true} />
  </Box>
);
