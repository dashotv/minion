import * as React from 'react';
import { IconType } from 'react-icons/lib';

import type { SvgIconTypeMap } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import type { OverridableComponent } from '@mui/material/OverridableComponent';
import Stack from '@mui/material/Stack';
import SvgIcon from '@mui/material/SvgIcon';

export interface ButtonMapButton {
  title: string;
  Icon: OverridableComponent<SvgIconTypeMap> & { muiName: string };
  color?: 'primary' | 'secondary' | 'error' | 'warning' | 'info' | 'success' | 'disabled' | 'action' | 'inherit';
  click?: (ev: React.MouseEvent<HTMLElement>) => void;
  Component?: IconType;
  link?: string;
}

export interface ButtonMapProps {
  size?: 'small' | 'medium' | 'large';
  buttons: ButtonMapButton[];
}

export function ButtonMap({ buttons, size = 'medium' }: ButtonMapProps) {
  return (
    <Stack spacing="5px" direction="row" sx={{ justifyContent: 'end' }}>
      {buttons.map(({ click, title, Icon, Component, color, link }, index) => (
        <IconButton
          key={index}
          sx={{ p: '0px' }}
          size={size}
          onClick={ev => {
            click && click(ev);
          }}
          href={link || ''}
          target={link ? '_blank' : undefined}
          title={title}
        >
          {Component ? (
            <SvgIcon inheritViewBox component={Component} fontSize={size} color={color} />
          ) : (
            <Icon fontSize={size} color={color} />
          )}
        </IconButton>
      ))}
    </Stack>
  );
}
