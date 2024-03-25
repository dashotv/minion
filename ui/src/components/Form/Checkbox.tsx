import { Controller } from 'react-hook-form';

import { title } from 'radash';

import CheckBoxIcon from '@mui/icons-material/CheckBox';
import CheckBoxOutlineBlankIcon from '@mui/icons-material/CheckBoxOutlineBlank';
import MUICheckbox from '@mui/material/Checkbox';
import FormControl from '@mui/material/FormControl';
import FormControlLabel from '@mui/material/FormControlLabel';

import { InputProps } from './types';

export interface IconCheckboxProps {
  icon: React.ReactNode;
  checkedIcon: React.ReactNode;
}
export const IconCheckbox = ({
  name,
  label,
  icon,
  disabled,
  checkedIcon,
  control,
  sx,
  onChange,
}: InputProps & IconCheckboxProps) => {
  //   error={!!errors[name]}
  //   helperText={errors[name]?.message}
  return (
    <Controller
      name={name}
      control={control}
      render={({ field }) => (
        <FormControl variant="standard">
          <FormControlLabel
            sx={sx}
            control={
              <MUICheckbox
                {...field}
                disabled={disabled}
                icon={icon}
                checkedIcon={checkedIcon}
                checked={field.value}
                id={name}
                onChange={e => {
                  field.onChange(e);
                  onChange && onChange(e);
                }}
              />
            }
            label={label}
            title={title(label || name)}
          />
        </FormControl>
      )}
    />
  );
};

export const Checkbox = ({ name, label, disabled, control, sx, onChange }: InputProps) => {
  //   error={!!errors[name]}
  //   helperText={errors[name]?.message}
  return (
    <IconCheckbox
      {...{ name, label, disabled, control, sx, onChange }}
      icon={<CheckBoxOutlineBlankIcon />}
      checkedIcon={<CheckBoxIcon />}
    />
  );
};
