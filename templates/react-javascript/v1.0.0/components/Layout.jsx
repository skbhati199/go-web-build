import React from 'react';

function Layout({ children }) {
  return (
    <div className="layout">
      <header>
        <nav>
          <h1>My React App</h1>
          <ul>
            <li><a href="/">Home</a></li>
            <li><a href="/about">About</a></li>
          </ul>
        </nav>
      </header>
      <main>{children}</main>
      <footer>
        <p>&copy; {new Date().getFullYear()} My React App</p>
      </footer>
    </div>
  );
}

export default Layout;