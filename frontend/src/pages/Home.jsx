// src/pages/Home.jsx
import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthContext from '../context/AuthContext';

export default function Home() {
  const { openGameConfigModal } = useContext(AuthContext);
  const navigate = useNavigate();

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      height: '80vh',
      textAlign: 'center'
    }}>
      <h1>Welcome to 24 Game</h1>
      <p>Solve arithmetic puzzles to make 24!</p>
      <button
        onClick={openGameConfigModal}
        style={{
          padding: '12px 24px',
          fontSize: '1.2rem',
          backgroundColor: '#3498db',
          color: 'white',
          border: 'none',
          borderRadius: '6px',
          cursor: 'pointer'
        }}
      >
        Start Game
      </button>
    </div>
  );
}
