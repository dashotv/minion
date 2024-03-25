import Typography from '@mui/material/Typography';

import { Chrono } from 'components/common';

export const Published = ({ date }: { date: string }) => {
  const never = (Date.now() - Date.parse(date)) / (1000 * 60 * 60 * 24 * 365) > 100;
  return (
    <Typography variant="subtitle2" color="gray" pl="3px" noWrap>
      {never ? 'never' : <Chrono fromNow>{date}</Chrono>}
    </Typography>
  );
};
