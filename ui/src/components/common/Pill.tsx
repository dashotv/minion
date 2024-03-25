import { useTheme } from '@mui/material';
import Box from '@mui/material/Box';

// TODO: support size, move styling to scss
export const Pill = ({
  name,
  value,
  color = 'primary',
  variant = 'row',
  icon,
  width,
}: {
  name: string;
  value?: React.ReactElement | string | number;
  color?: string;
  variant?: 'row' | 'column';
  icon?: React.ReactElement;
  width?: string | number;
  size?: 'small' | 'medium' | 'large';
}) => {
  const display = variant == 'row' ? 'flex' : 'block';
  const theme = useTheme();

  const colorValue = () => {
    switch (color) {
      case 'primary':
        return theme.palette.primary.main;
      case 'secondary':
        return theme.palette.secondary.main;
      case 'error':
        return theme.palette.error.main;
      case 'warning':
        return theme.palette.warning.main;
      case 'info':
        return theme.palette.info.main;
      case 'success':
        return theme.palette.success.main;
      default:
        return color;
    }
  };
  return (
    <Box
      sx={{
        display: display,
        border: '1px solid',
        borderRadius: '8px',
        borderColor: colorValue(),
        overflow: 'hidden',
        width: 'fit-content',
        height: variant == 'row' ? '22px' : '42px',
      }}
    >
      <Box
        sx={{
          color: 'black',
          backgroundColor: colorValue(),
          pr: '4px',
          pl: '4px',
          fontSize: icon ? 'auto' : 'small',
        }}
      >
        {icon ? icon : name}
      </Box>
      <Box
        sx={{
          color: colorValue(),
          backgroundColor: '#363636',
          pr: '4px',
          pl: '4px',
          overflow: 'hidden',
          maxWidth: width,
          whiteSpace: 'nowrap',
          textOverflow: 'ellipsis',
          fontSize: 'small',
        }}
      >
        {value}
      </Box>
    </Box>
  );
};
