import axios from 'axios';
import { User, LoginRequest, RegisterRequest } from '../types';

// Relative URLs work in deployment (same domain), localhost for development
// In Render, frontend and backend are on the same domain when using proximity
const API_BASE_URL = window.location.hostname === 'localhost' ? 'http://localhost:8080' : '';
const API_URL = `${API_BASE_URL}/api/auth`;

export const authService = {
  // Login de usuario
  async login(credentials: LoginRequest): Promise<User> {
    const response = await axios.post<User>(`${API_URL}/login`, credentials);
    return response.data;
  },

  // Registro de usuario
  async register(data: RegisterRequest): Promise<User> {
    const response = await axios.post<User>(`${API_URL}/register`, data);
    return response.data;
  }
};
