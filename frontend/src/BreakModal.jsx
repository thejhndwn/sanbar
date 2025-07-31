// components/BreakModal.tsx
import { useEffect, useState } from 'react';

export default function BreakModal({ isOpen, onContinue } ) {
  const [timeLeft, setTimeLeft] = useState(300); // 300 seconds = 5 minutes

  useEffect(() => {
    if (!isOpen) return;

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          onContinue(); // Auto-continue when time's up
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [isOpen, onContinue]);

  if (!isOpen) return null;

  const minutes = Math.floor(timeLeft / 60);
  const seconds = timeLeft % 60;

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Take a Break</h2>
        <p>You’ve earned a 5-minute break.</p>
        <div className="countdown">
          {minutes}:{seconds.toString().padStart(2, '0')}
        </div>
        <button
          onClick={onContinue}
          disabled={timeLeft > 0}
        >
          {timeLeft > 0 ? 'Wait...' : 'Continue'}
        </button>
      </div>
    </div>
  );
}
