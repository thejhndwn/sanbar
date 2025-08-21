// src/api/authApi.js
import { api } from './client';

export const login = (data) => {
  return api.post('/api/login', data);
};

export const register = (data) => {
    return api.post('/api/register', data);
};

export const newuser = (data) => {
    return api.post('/api/newuser', data);
};
