// Game.jsx
import React, { useState, useEffect, useRef } from 'react';
import NewGameModal from '../NewGameModal';
import BreakModal from '../BreakModal';

const GAME_STATES = {
  READY: 'ready',
  PLAYING: 'playing',
  ENDED: 'ended',
  BREAK: 'break',
};

const Game = () => {
  const [gameState, setGameState] = useState(GAME_STATES.READY);

  // Game data from backend
  const [cards, setCards] = useState([]); // Each card: { id, value }
  const [currentScore, setCurrentScore] = useState(0);
  const [problemsSolved, setProblemsSolved] = useState(0);
  const [problemsRemaining, setProblemsRemaining] = useState(0);
  const [timeRemaining, setTimeRemaining] = useState(0);

  // Expression & selected items (for UI)
  const [selectedItems, setSelectedItems] = useState([]); // IDs or values
  const [operation, setOperation] = useState(null);

  // Undo stack: stores full state snapshots
  const [undoStack, setUndoStack] = useState([]);

  // Modals
  const [showNewGameModal, setShowNewGameModal] = useState(false);
  const [showBreakModal, setShowBreakModal] = useState(false);
  const [breakTimeLeft, setBreakTimeLeft] = useState(600);

  const intervalRef = useRef(null);

  // --- API Calls ---
  const startGame = async () => {
    try {
      const res = await fetch('/api/start', { method: 'POST' });
      const data = await res.json();
      loadGameData(data);
      setUndoStack([]); // Reset undo stack
      setGameState(GAME_STATES.PLAYING);
    } catch (err) {
      console.error('Failed to start game:', err);
    }
  };

  const applyOperation = async (card1Id, card2Id, op) => {
    try {
      const res = await fetch('/api/action', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ card1Id, card2Id, operation: op }),
      });
      const data = await res.json();

      if (data.valid) {
        // Save current state to undo stack before updating
        setUndoStack((prev) => [
          ...prev,
          { cards, selectedItems, operation, timeRemaining },
        ]);

        loadGameData(data);
        setSelectedItems([]);
        setOperation(null);
        checkForBreak(data.problemsSolved);
      }
      return data;
    } catch (err) {
      console.error('Operation failed:', err);
    }
  };

  const handleUndo = async () => {
    if (undoStack.length === 0) return;

    const lastState = undoStack[undoStack.length - 1];
    setUndoStack((prev) => prev.slice(0, -1));

    // Optimistically restore
    setCards(lastState.cards);
    setSelectedItems(lastState.selectedItems);
    setOperation(lastState.operation);
    setTimeRemaining(lastState.timeRemaining);

    // Tell backend to undo
    try {
      const res = await fetch('/api/undo', { method: 'POST' });
      const data = await res.json();
      if (data.valid) {
        loadGameData(data); // Ensure sync with backend
      }
    } catch (err) {
      console.error('Undo failed on backend:', err);
      // Optionally: reload game or show warning
    }
  };

  const skipProblem = async () => {
    try {
      const res = await fetch('/api/skip', { method: 'POST' });
      const data = await res.json();
      loadGameData(data);
      setSelectedItems([]);
      setOperation(null);
      checkForBreak(data.problemsSolved);
    } catch (err) {
      console.error('Skip failed:', err);
    }
  };

  const endGame = async () => {
    try {
      const res = await fetch('/api/end', { method: 'POST' });
      const data = await res.json();
      setGameState(GAME_STATES.ENDED);
      setCurrentScore(data.finalScore);
      setProblemsSolved(data.problemsSolved);
    } catch (err) {
      console.error('End game failed:', err);
    }
  };

  const loadGameData = (data) => {
    setCards(data.cards?.map((c) => ({ id: c.id, value: c.value })) || []);
    setCurrentScore(data.score || 0);
    setProblemsSolved(data.problemsSolved || 0);
    setProblemsRemaining(data.problemsRemaining || 0);
    setTimeRemaining(data.timeRemaining || 0);
  };

  const checkForBreak = (currentProblemsSolved) => {
    if (currentProblemsSolved > 0 && currentProblemsSolved % 100 === 0) {
      setBreakTimeLeft(600);
      setGameState(GAME_STATES.BREAK);
      setShowBreakModal(true);
    }
  };

  // --- Break Timer ---
  useEffect(() => {
    if (gameState === GAME_STATES.BREAK) {
      intervalRef.current = setInterval(() => {
        setBreakTimeLeft((prev) => {
          if (prev <= 1) {
            clearInterval(intervalRef.current);
            setShowBreakModal(false);
            setGameState(GAME_STATES.PLAYING);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [gameState]);

  // --- Problem Timer ---
  useEffect(() => {
    if (gameState === GAME_STATES.PLAYING && timeRemaining > 0) {
      intervalRef.current = setInterval(() => {
        setTimeRemaining((prev) => {
          if (prev <= 1) {
            clearInterval(intervalRef.current);
            skipProblem();
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [gameState, timeRemaining]);

  // --- Auto-submit when one card is 24 ---
  useEffect(() => {
    if (gameState === GAME_STATES.PLAYING && cards.length === 1 && cards[0].value === 24) {
      setTimeout(async () => {
        try {
          await fetch('/api/submit', { method: 'POST' });
          // Backend will advance to next problem
          const res = await fetch('/api/state'); // or include in submit response
          const data = await res.json();
          loadGameData(data);
          checkForBreak(data.problemsSolved);
        } catch (err) {
          console.error('Auto-submit failed:', err);
        }
      }, 600);
    }
  }, [cards, gameState]);

  // --- Handlers ---

  const handleCardClick = (card) => {
    if (selectedItems.length >= 2) return; // Already selected two

    const newSelection = [...selectedItems, card];
    setSelectedItems(newSelection);

    if (newSelection.length === 2 && operation) {
      const [c1, c2] = newSelection;
      applyOperation(c1.id, c2.id, operation);
    }
  };

  const handleOperatorClick = (op) => {
    if (selectedItems.length === 0) return;
    if (selectedItems.length === 1) {
      setOperation(op);
    }
    // If two selected, apply immediately
    if (selectedItems.length === 2) {
      const [c1, c2] = selectedItems;
      applyOperation(c1.id, c2.id, op);
    }
  };

  const handleSkip = () => {
    skipProblem();
    setSelectedItems([]);
    setOperation(null);
  };

  const handleEndGame = () => {
    endGame();
    setSelectedItems([]);
    setOperation(null);
  };

  const handlePlayAgain = () => {
    startGame();
    setSelectedItems([]);
    setOperation(null);
  };

  const handleNewGame = () => {
    setShowNewGameModal(true);
  };

  const handleExitBreakEarly = () => {
    setShowBreakModal(false);
    setGameState(GAME_STATES.PLAYING);
    clearInterval(intervalRef.current);
  };

  // --- Render ---

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const renderReadyState = () => (
    <div className="game-ready">
      <h2>🎯 Ready to Play 24?</h2>
      <button onClick={startGame}>Start Game</button>
      <button onClick={handleNewGame}>New Game (Config)</button>
    </div>
  );

  const renderInGameState = () => (
    <div className="game-playing">
      <div className="game-header">
        <div>Score: <strong>{currentScore}</strong></div>
        <div>Problems: {problemsSolved} / {problemsSolved + problemsRemaining}</div>
        <div>Time: <strong>{formatTime(timeRemaining)}</strong></div>
      </div>

      <div className="cards-container">
        {cards.map((card) => (
          <button
            key={card.id}
            className={`card-btn ${selectedItems.find(c => c.id === card.id) ? 'selected' : ''}`}
            onClick={() => handleCardClick(card)}
            disabled={selectedItems.length >= 2 && !selectedItems.find(c => c.id === card.id)}
          >
            {card.value}
          </button>
        ))}
      </div>

      <div className="operation-area">
        {selectedItems.length > 0 && (
          <div className="selection">
            Selected: {selectedItems.map(c => c.value).join(', ')}
            {selectedItems.length === 1 && !operation && (
              <div>
                Choose operator: + − × ÷
              </div>
            )}
          </div>
        )}

        {selectedItems.length === 2 && !operation && (
          <div className="operators">
            {['+', '-', '*', '/'].map((op) => (
              <button key={op} onClick={() => handleOperatorClick(op)}>
                {op === '*' ? '×' : op === '/' ? '÷' : op}
              </button>
            ))}
          </div>
        )}

        {operation && selectedItems.length === 1 && (
          <div>Operator: <strong>{operation}</strong></div>
        )}
      </div>

      <div className="actions">
        <button onClick={handleUndo} disabled={undoStack.length === 0}>
          🔄 Undo
        </button>
        <button onClick={handleSkip}>⏭️ Skip</button>
        <button onClick={handleEndGame}>⏹️ End Game</button>
      </div>
    </div>
  );

  const renderEndGameState = () => (
    <div className="game-ended">
      <h2>🎉 Game Over!</h2>
      <p>Final Score: <strong>{currentScore}</strong></p>
      <p>Problems Solved: {problemsSolved}</p>
      <button onClick={handlePlayAgain}>🔁 Play Again</button>
      <button onClick={handleNewGame}>🆕 New Game</button>
    </div>
  );

  return (
    <div className="game-container">
      {gameState === GAME_STATES.READY && renderReadyState()}
      {gameState === GAME_STATES.PLAYING && renderInGameState()}
      {gameState === GAME_STATES.ENDED && renderEndGameState()}

      {/* Modals */}
      {showNewGameModal && (
        <NewGameModal onClose={() => setShowNewGameModal(false)} onStart={startGame} />
      )}

      {showBreakModal && (
        <BreakModal
          timeLeft={breakTimeLeft}
          onExitEarly={handleExitBreakEarly}
        />
      )}
    </div>
  );
};

export default Game;
