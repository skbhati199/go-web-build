import React from 'react';
import { LayoutProps } from '@/types/components';
import styles from './Layout.module.scss';

export const Layout: React.FC<LayoutProps> = ({
  children,
  header,
  footer,
  sidebar,
  className,
}) => {
  return (
    <div className={`${styles.layout} ${className || ''}`}>
      {header && <header className={styles.header}>{header}</header>}
      <div className={styles.container}>
        {sidebar && <aside className={styles.sidebar}>{sidebar}</aside>}
        <main className={styles.main}>{children}</main>
      </div>
      {footer && <footer className={styles.footer}>{footer}</footer>}
    </div>
  );
};