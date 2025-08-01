// src/context/AuthContext.jsx
import React, { createContext, useState } from 'react';

export const AuthContext = createContext();

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

  return (
    <AuthContext.Provider value={{
      user,
      guestId,
      login,
      logout,
      openLoginModal,
      openSettingsModal,
      showLoginModal,
      showSettingsModal,
      closeLoginModal,
      closeSettingsModal,
    }}>
      {children}

      {/* Modals */}
      {showLoginModal && <LoginModal onClose={closeLoginModal} onLogin={login} />}
      {showSettingsModal && <SettingsModal onClose={closeSettingsModal} />}
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
