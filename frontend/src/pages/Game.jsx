import NewGameModal from '../NewGameModal';
import BreakModal from '../BreakModal';
import { useParams } from 'react-router-dom'
import { startGame, submitSolve, skipSolve } from '../api/gameApi';
import React, { useState, useEffect, useRef } from 'react';
import { evaluate } from 'mathjs';

const GAME_STATES = {
    READY: 'ready',
    PLAYING: 'playing',
    ENDED: 'ended',
    BREAK: 'break',
};

// Specifically this is for re
const Game = () => {
    const [gameState, setGameState] = useState(GAME_STATES.READY);
    const { gameID } = useParams();

    // Game data from backend
    //
    const cardIdRef = useRef(0)
    const createCard = (value, position) => ({
        id: cardIdRef.current++,
        value: Number(value),
        position: position,
        active: true,
    });
    // TODO: Put this in their own component? 
    const [cards, setCards] = useState([]); // Each card: { id, value }
    const [currentScore, setCurrentScore] = useState(0);
    const [problemsSolved, setProblemsSolved] = useState(0);
    const [problemsRemaining, setProblemsRemaining] = useState(0);
    const [timeRemaining, setTimeRemaining] = useState(0);

    // Expression & selected items (for UI)
    const [ selectedItems, setSelectedItems] = useState([]); // IDs or values
    const [ operations, setOperations] = useState([]);
    const [ cardSelector, setCardSelector ] = useState('');
    const [ operatorSelector, setOperatorSelector ] = useState('');

    const [showNewGameModal, setShowNewGameModal] = useState(false);
    const [showBreakModal, setShowBreakModal] = useState(false);
    const [breakTimeLeft, setBreakTimeLeft] = useState(600);

    // TODO: what is this? 
    const intervalRef = useRef(null);

    // --- API Calls ---
    const beginGame = async () => {
        try {
            const {combo} = await startGame(gameID)
            loadGameData(combo, {});
            setGameState(GAME_STATES.PLAYING);
        } catch (err) {
            console.error('Failed to start game:', err);
        }
    };

    const skipProblem = async () => {
        try {
            const {combo} = await skipSolve(gameID)
            loadGameData(combo);
            setSelectedItems([])
            setOperations(null);
            //checkForBreak(data.problemsSolved);
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

    const loadGameData = (combo, data) => {
        const newCards = combo.split("-").map(createCard);
        setCards(newCards || []);
        //TODO: change the below to use the data.values OR use the frontend saved values
        setCurrentScore(data.score || 0);
        setProblemsSolved(data.problemsSolved || 0);
        setProblemsRemaining(data.problemsRemaining || 0);
        setTimeRemaining(data.timeRemaining || 0);
    };

    const checkForBreak = (currentProblemsSolved) => {
        //TODO: make constant
        if (currentProblemsSolved > 0 && currentProblemsSolved % 100 === 0) {
            setBreakTimeLeft(600);
            setGameState(GAME_STATES.BREAK);
            setShowBreakModal(true);
        }
    };

    const addExpression = (selectedCard) => {
        // add card state to expressions
        console.log("we are tyring to make an expression")
        console.log(selectedItems)
        const selected = selectedItems[0];
        console.log("selected:", selected)
        setOperations([...operations, {
            cards,
            selected,
        }])
        console.log("operation set")
        // clear selectedItems, add the new expressionSol and update the cardSelector
        const operatorMap = { 'x': '*', '%': '/' };
        const operator = operatorMap[selectedItems[1]] || selectedItems[1];
        const expr = selected.value + operator + selectedCard.value
        console.log("we made an expr:", expr)
        const result = evaluate(expr); //TODO: turn this into a card
        console.log("evaluated to:", result)
        const resultCard = createCard(result, selectedCard.position)
        setSelectedItems([resultCard])
        setCards(prevCards => {
            const updatedCards = prevCards.map(card =>
                card.id === selected.id || selectedCard.id === card.id
                ? { ...card, active: false}
                : card
            );
            console.log(updatedCards)
            return [...updatedCards, resultCard]
        })
        setCardSelector(resultCard)
        setOperatorSelector(null)
    }

    const popOperation = () => {
        console.log("we are trying to pop an operation")
        console.log(operations)
        if (operations.length === 0) {
            return
        }
        const lastState = operations[operations.length - 1];
        setCards(lastState.cards)
        setCardSelector(lastState.selected)
        setSelectedItems([lastState.selected])
        setOperations(prev => prev.slice(0, -1));
    }

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
            }, 10);
        }

        return () => {
            if (intervalRef.current) clearInterval(intervalRef.current);
        };
    }, [gameState, timeRemaining]);

    // --- Auto-submit when one card is 24 ---
    useEffect(() => {
        const activeCards = cards.filter(card => card.active);
        if (gameState === GAME_STATES.PLAYING && activeCards.length === 1 && activeCards[0].value === 24) {
            setTimeout(async () => {
                try {
                    console.log("we are trying to make the submission")
                    const { combo } = await submitSolve(gameID)
                    loadGameData(combo);
                    //checkForBreak(data.problemsSolved);
                } catch (err) {
                    console.error('Auto-submit failed:', err);
                }
            }, 600);
        }
    }, [cards, gameState]);

    // --- Handlers ---

    const handleCardClick = (card) => {
        if (selectedItems.length == 0) {
            setSelectedItems([card])
            setCardSelector(card)
        }
        else if (selectedItems.length == 1) {
            setSelectedItems([card])
            setCardSelector(card)
        }
        else if (selectedItems.length == 2) {
            addExpression(card)
        }
    };

    const handleOperatorClick = (operator) => {
        if (selectedItems.length == 0) {
        }
        else if (selectedItems.length == 1) {
            setSelectedItems([...selectedItems, operator])
            setOperatorSelector(operator)
        }
        else if (selectedItems.length ==2) {
            setSelectedItems([selectedItems[0], operator])
            setOperatorSelector(operator)
        }
    };

    const handleSkip = () => {
        console.log("handleskip pressed")
        skipProblem();
    };

    const handleEndGame = () => {
        endGame();
        setSelectedItems([]);
        setOperations(null);
    };

    const handlePlayAgain = () => {
        beginGame();
        setSelectedItems([]);
        setOperations(null);
    };

    const handleUndo = () => {
        console.log("undo stuff")
        console.log(selectedItems)
        if (selectedItems.length ==0) {
        }
        else if (selectedItems.length ==1) {
            popOperation()
        }
        else if (selectedItems.length ==2) {
            setSelectedItems(selectedItems.slice(0,1))
            setOperatorSelector(null)
        }
        console.log('this is the undo')
        console.log(selectedItems)
    }   


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
        <button onClick={beginGame}>Start Game</button>
        </div>
    );

    const renderInGameState = () => (
        <div className="game-playing">
        <div className="game-header">
        <div>Score: <strong>{currentScore}</strong></div>
        <div>Problems: {problemsSolved} / {problemsSolved + problemsRemaining}</div>
        <div>Time: <strong>{formatTime(timeRemaining)}</strong></div>
        </div>

        <div className="cards-container" style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '12px', padding: '10px' }}>
        {Array.from({ length: 4 }, (_, pos) => {
            const card = cards.find(c => c.position === pos && c.active);
            return (
                <div key={pos} className="card-slot">
                {card ? (
                    <button
                    className={`card-btn ${cardSelector.id === card.id ? 'selected' : ''}`}
                    onClick={() => handleCardClick(card)}
                    disabled={cardSelector.id === card.id}
                    >
                    {card.value}
                    </button>
                ) : (
                    <div className="card-placeholder"></div>
                )}
                </div>
            );
        })}
        </div>
        <div className="operations-container">
        {['+', '-', 'x', '%'].map((operator) => (
            <button
            key={operator}
            className={`operation-btn`}
            onClick={() => handleOperatorClick(operator)}
            disabled={operatorSelector === operator}
            >
            {operator}
            </button>
        ))}

        </div>


        <div className="actions">
        <button onClick={handleUndo}>🔄 Undo </button>
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
        </div>
    );

    return (
        <div className="game-container">
        {gameState === GAME_STATES.READY && renderReadyState()}
        {gameState === GAME_STATES.PLAYING && renderInGameState()}
        {gameState === GAME_STATES.ENDED && renderEndGameState()}

        {/* Modals */}
        {showNewGameModal && (
            <NewGameModal onClose={() => setShowNewGameModal(false)} onStart={beginGame} />
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
