// src/api/gameApi.js
import { api } from './client';

export const createGame = (data) => {
  return api.post('/api/newgame', data);
};

export const startGame = (gameId) => {
    return api.post('/api/game/start', {gameId});
};

export const submitSolve = (gameId) => {
  return api.post(`/api/game/submit`, { gameId});
};

export const endGame = (gameId) => {
    return api.post('/api/game/end', { gameId});
}

export const skipSolve = (gameId) => {
    return api.post(`/api/game/skip`, { gameId});
};

export const takeBreak = (gameId) => {
  return api.post(`/api/game/break`, { gameId });
};

export const continueGame = (gameId) => {
  return api.post(`/api/game/continue`, { gameId });
};
