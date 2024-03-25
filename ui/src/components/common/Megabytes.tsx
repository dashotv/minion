import Typography from '@mui/material/Typography';

export const Megabytes = ({ value, ord }: { value: number; ord: 'bytes' | 'kilobytes' | 'megabytes' }) => {
  const getOutput = (v: number) => {
    switch (ord) {
      case 'bytes':
        return v / (1024 * 1024);
      case 'kilobytes':
        return v / 1024;
      default:
        return v;
    }
  };

  return (
    <Typography variant="subtitle2" color="textSecondary" title={`${value} ${ord}`}>
      {getOutput(value).toFixed(2)}mb
    </Typography>
  );
};
