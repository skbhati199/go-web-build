import React from 'react';

function Home() {
  return (
    <div className="home-page">
      <h2>Welcome to My React App</h2>
      <p>This is a starter template for your React application.</p>
      <button onClick={() => alert('Hello from React!')}>Click Me</button>
    </div>
  );
}

export default Home;