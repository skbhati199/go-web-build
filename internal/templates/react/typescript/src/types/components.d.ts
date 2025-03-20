import { ReactNode } from 'react';

export interface BaseProps {
  className?: string;
  children?: ReactNode;
}

export interface ButtonProps extends BaseProps {
  onClick?: () => void;
  disabled?: boolean;
  variant?: 'primary' | 'secondary' | 'outline';
  size?: 'small' | 'medium' | 'large';
  type?: 'button' | 'submit' | 'reset';
}

export interface InputProps extends BaseProps {
  value?: string;
  onChange?: (value: string) => void;
  placeholder?: string;
  type?: 'text' | 'password' | 'email' | 'number';
  error?: string;
  label?: string;
}

export interface CardProps extends BaseProps {
  title?: string;
  footer?: ReactNode;
  header?: ReactNode;
}

export interface LayoutProps extends BaseProps {
  header?: ReactNode;
  footer?: ReactNode;
  sidebar?: ReactNode;
}