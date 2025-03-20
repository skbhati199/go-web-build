import React from 'react';
import { ButtonProps } from '../../types/components';
import styles from './Button.module.css';

export const Button: React.FC<ButtonProps> = ({
  children,
  className,
  disabled = false,
  variant = 'primary',
  size = 'medium',
  type = 'button',
  onClick,
}) => {
  const buttonClass = `${styles.button} ${styles[variant]} ${styles[size]} ${className || ''}`;

  return (
    <button
      type={type}
      className={buttonClass}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>
  );
};

export default Button;