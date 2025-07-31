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
      const { id } = await createGame({
        game_type: 'solo_survival',
        num_cards: 4,
        target: 24,
      });
      onClose();
      navigate(`/game/${id}`);
    } catch (err) {
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

        <div className="quick-play">
          <h3>Quick Play</h3>
          <button onClick={startQuickGame} disabled={loading}>
            {loading ? 'Starting...' : '4 Cards → 24'}
          </button>
        </div>

        <div className="custom-play">
          <h3>Custom</h3>
          <button onClick={startCustom}>Configure Game</button>
        </div>

        <button className="close-btn" onClick={onClose}>
          Cancel
        </button>
      </div>
    </div>
  );
}
