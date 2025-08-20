// src/context/AuthContext.jsx
import React, { createContext, useState, useEffect } from 'react';
import { newuser} from './api/authApi';

export const AuthContext = createContext();

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null);
    useEffect(() => {
        const initUser = async () => {
            console.log("sanity check that we're grabbing the token")
            const savedToken = localStorage.getItem('sanbarToken')
            let token = savedToken;
            if (!savedToken) {
                token = crypto.randomUUID();
                localStorage.setItem('sanbarToken', token);

                try {
                    await newuser({
                        "token": token
                    })
                    console.log("making the new user was great")
                } catch (error) {
                    console.error("Failed to create user:", error);
                }
            }

        };

        initUser();

//        initUser().finally(() => {
//           setLoading(false)
//        })
    }, [])


    const [showLoginModal, setShowLoginModal] = useState(false);
    const [showSettingsModal, setShowSettingsModal] = useState(false);

    // TODO: add request to login route for api client
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
