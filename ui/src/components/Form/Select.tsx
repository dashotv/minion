import { Controller } from 'react-hook-form';

import { title } from 'radash';

import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import { InputProps, Option } from './types';

export interface SelectProps {
  options: Option[];
  render?: (option: Option) => React.ReactNode;
}
export const Select = ({ name, label, disabled, options, control, sx, render, onChange }: InputProps & SelectProps) => {
  const renderDefault = (option: Option) => (
    <MenuItem key={option.value} value={option.value}>
      <span>{render ? render(option) : option.label}</span>
    </MenuItem>
  );
  return (
    <Controller
      name={name}
      control={control}
      render={({ field }) => (
        <TextField
          {...field}
          select
          sx={sx}
          id={name}
          disabled={disabled}
          label={title(label || name)}
          variant="standard"
          fullWidth
          onChange={e => {
            onChange && onChange(e);
            field.onChange(e);
          }}
        >
          {options.map(renderDefault)}
        </TextField>
      )}
    />
  );
};
