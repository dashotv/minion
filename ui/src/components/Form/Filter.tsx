import { Checkbox, MenuItem, TextField } from '@mui/material';

export interface Choice {
  name: string;
  type: string;
}

export const FilterSelect = ({
  name,
  choices,
  choose,
  selected,
}: {
  name: string;
  selected: string;
  choices: Choice[];
  choose: (choice: string) => void;
}) => {
  const handleChange = event => {
    choose(event.target.value);
  };
  return (
    <TextField
      select
      id={name}
      label={name}
      variant="standard"
      fullWidth
      onChange={handleChange}
      value={selected}
      defaultValue={choices[0].type}
    >
      {choices.map(({ name, type }) => (
        <MenuItem key={type} value={type}>
          {name}
        </MenuItem>
      ))}
    </TextField>
  );
};

export const FilterCheckbox = ({
  checked,
  icon,
  checkedIcon,
  change,
}: {
  checked: boolean;
  icon: React.ReactNode;
  checkedIcon: React.ReactNode;
  change: (checked: boolean) => void;
}) => {
  const handleChange = event => {
    change(event.target.checked);
  };
  return <Checkbox checked={checked} icon={icon} checkedIcon={checkedIcon} onChange={handleChange} />;
};
