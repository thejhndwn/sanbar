// components/NewGameModal.tsx
import { useState } from 'react';
import { createGame } from './api/gameApi';
import { useNavigate } from 'react-router-dom';


export default function NewGameModal({ isOpen, onClose }) {
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    if (!isOpen) return null;

    const startQuickGame = async () => {
        setLoading(true);
        try {
            console.log("welcome to the thunderdome")
            const { id }= await createGame({
                game_type: 'solo_survival',
                num_cards: 4,
                target: 24,
            });
            console.log("I cannot belive the game was created")
            onClose();
            navigate(`/game/${id}`);
        } catch (err) {
            console.log(err)
            alert('Failed to start game');
            setLoading(false);
        }
    };

    const startCustom = () => {
        // TODO: Implement custom config
        alert('Custom mode coming soon!');
    };

    return (
        <div className="modal-overlay">
        <div className="modal">
        <h2>Start New Game</h2>

        <div className="modal-options">
        <button onClick={startQuickGame} disabled={loading}>
        {loading ? 'Starting...' : '24'}
        </button>

        <button onClick={startCustom}>Custom</button>
        </div>

        <button className="close-btn" onClick={onClose}>
        Cancel
        </button>
        </div>
        </div>
    );
}
