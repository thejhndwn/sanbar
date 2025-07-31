// src/api/gameApi.js
import { api } from './client';

export const createGame = (data) => {
  return api.post('/api/games', data);
};

export const solvePuzzle = (gameId, solution) => {
  return api.post(`/api/game/${gameId}/solve`, { solution });
};

export const takeBreak = (gameId) => {
  return api.post(`/api/game/${gameId}/take-break`, {});
};

export const continueGame = (gameId) => {
  return api.post(`/api/game/${gameId}/continue`, {});
};
