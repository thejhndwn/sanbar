// src/api/authApi.js
import { api } from './client';

export const login = (data) => {
  return api.post('/api/login', data);
};

export const register = (data) => {
    return api.post('/api/register', date);
};


/**
export const solvePuzzle = (gameId, solution) => {
  return api.post(`/api/game/${gameId}/solve`, { solution });
};

export const takeBreak = (gameId) => {
  return api.post(`/api/game/${gameId}/take-break`, {});
};

export const continueGame = (gameId) => {
  return api.post(`/api/game/${gameId}/continue`, {});
};
**/ 

