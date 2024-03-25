import Chip from '@mui/material/Chip';
import Typography from '@mui/material/Typography';

import { Pill } from '.';

export type ResolutionProps = {
  resolution: string | number | undefined;
  variant?: 'default' | 'chip' | 'pill';
};
export const Resolution = ({ resolution, variant }: ResolutionProps) => {
  if (!resolution) {
    return;
  }

  switch (variant) {
    case 'chip':
      return <Chip label={resolution} size="small" color="primary" />;
    case 'pill':
      return <Pill name="R" value={resolution} color="primary.dark" />;
    default:
      return (
        <Typography variant="button" color="primary.dark">
          {resolution}p
        </Typography>
      );
  }
};
export const ResolutionTitle = ({ title, variant }: { title: string; variant?: 'default' | 'chip' | 'pill' }) => {
  if (!title) {
    return;
  }
  if (title.includes('2160') || title.includes('2160p') || title.includes('4k')) {
    return <Resolution resolution={2160} variant={variant} />;
  }
  if (title.includes('1080') || title.includes('1080p') || title.includes('2k')) {
    return <Resolution resolution={1080} variant={variant} />;
  }
  if (title.includes('720') || title.includes('720p')) {
    return <Resolution resolution={720} variant={variant} />;
  }
  return '';
};
