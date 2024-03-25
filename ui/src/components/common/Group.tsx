import Chip from '@mui/material/Chip';
import Typography from '@mui/material/Typography';

import { Pill } from '.';

export type GroupProps = {
  group: string | undefined;
  author: string | undefined;
  variant?: 'default' | 'chip' | 'pill';
};
export const Group = ({ group, author, variant }: GroupProps) => {
  const value = group || (author && `[${author}]`) || undefined;
  if (!value) {
    return;
  }

  switch (variant) {
    case 'chip':
      return <Chip label={value} size="small" color="secondary" />;
    case 'pill':
      return <Pill name="G" value={<Typography noWrap>{value}</Typography>} color="secondary.dark" />;
    default:
      return (
        <Typography title={value} maxWidth="95px" noWrap variant="button" color="secondary.dark">
          {value}
        </Typography>
      );
  }
};
