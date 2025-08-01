// src/pages/Home.jsx
import React, { useState } from 'react';
import NewGameModal from '../NewGameModal';

export default function Home() {

  const [showNewGameModal, setShowNewGameModal] = useState(false);


  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      height: '80vh',
      textAlign: 'center'
    }}>
      <h1>Sanbar</h1>
      <button
        onClick={() => setShowNewGameModal(true)}
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
       Play 
      </button>
      {showNewGameModal &&
      <NewGameModal isOpen={showNewGameModal} onClose={() => setShowNewGameModal(false)}/>}
    </div>
  );
}
