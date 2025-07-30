// src/context/AuthContext.jsx
import React, { createContext, useState } from 'react';

const AuthContext = createContext();

export default AuthContext;

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [guestId] = useState(() => {
    const saved = localStorage.getItem('guestId');
    if (saved) return saved;
    const id = crypto.randomUUID();
    localStorage.setItem('guestId', id);
    return id;
  });

  const [showLoginModal, setShowLoginModal] = useState(false);
  const [showSettingsModal, setShowSettingsModal] = useState(false);
  const [showGameConfigModal, setShowGameConfigModal] = useState(false);

  const login = (userData) => {
    setUser(userData);
    setShowLoginModal(false);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  const openLoginModal = () => setShowLoginModal(true);
  const closeLoginModal = () => setShowLoginModal(false);
  const openSettingsModal = () => setShowSettingsModal(true);
  const closeSettingsModal = () => setShowSettingsModal(false);
  const openGameConfigModal = () => setShowGameConfigModal(true);
  const closeGameConfigModal = () => setShowGameConfigModal(false);

  return (
    <AuthContext.Provider value={{
      user,
      guestId,
      login,
      logout,
      openLoginModal,
      openSettingsModal,
      openGameConfigModal,
      showLoginModal,
      showSettingsModal,
      showGameConfigModal,
      closeLoginModal,
      closeSettingsModal,
      closeGameConfigModal
    }}>
      {children}

      {/* Modals */}
      {showLoginModal && <LoginModal onClose={closeLoginModal} onLogin={login} />}
      {showSettingsModal && <SettingsModal onClose={closeSettingsModal} />}
      {showGameConfigModal && <GameConfigModal onClose={closeGameConfigModal} />}
    </AuthContext.Provider>
  );
}

// Placeholder modals (we'll implement later)
function LoginModal({ onClose, onLogin }) {
  return (
    <div className="modal-overlay">
      <div className="modal">
        <h3>Login</h3>
        <button onClick={() => onLogin({ email: 'test@example.com' })}>Mock Login</button>
        <button onClick={onClose}>Close</button>
      </div>
    </div>
  );
}

function SettingsModal({ onClose }) {
  return (
    <div className="modal-overlay">
      <div className="modal">
        <h3>Settings</h3>
        <p>Game settings go here.</p>
        <button onClick={onClose}>Close</button>
      </div>
    </div>
  );
}

function GameConfigModal({ onClose }) {
  return (
    <div className="modal-overlay">
      <div className="modal" style={{ width: '400px' }}>
        <h3>Game Mode</h3>
        <div>
          <label>
            <input type="radio" name="mode" defaultChecked /> Survival
          </label>
        </div>
        <div>
          <label>Number of Cards: 
            <input type="number" defaultValue="4" min="2" max="6" />
          </label>
        </div>
        <div>
          <label>Target Number: 
            <input type="number" defaultValue="24" />
          </label>
        </div>
        <div style={{ marginTop: '1rem' }}>
          <button onClick={() => {
            onClose();
            window.location.href = '/game'; // Redirect after config
          }}>
            Start Game
          </button>
          <button onClick={onClose}>Cancel</button>
        </div>
      </div>
    </div>
  );
}
