import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import  { AuthContext } from './AuthContext';

export default function Navbar() {
  console.log("flag 1", AuthContext)
  const { user, openLoginModal } = useContext(AuthContext);
  console.log("there's an error?")

  return (
    <nav style={{
      padding: '1rem',
      backgroundColor: '#2c3e50',
      color: 'white',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center'
    }}>
      <div>
        <Link to="/" style={{ color: 'white', marginRight: '1rem' }}>Home</Link>
      </div>
      <div>
        {user ? (
          <span>Welcome, {user.email}!</span>
        ) : (
          <button onClick={openLoginModal} style={{ marginRight: '0.5rem' }}>Login</button>
        )}
      </div>
    </nav>
  );
}
