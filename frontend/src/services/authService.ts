import axios from 'axios';
import { User, LoginRequest, RegisterRequest } from '../types';

// Configure backend URLs for each environment
export const getApiBaseUrl = (hostname = window.location.hostname, envVar?: string) => {
  if (hostname === 'localhost') {
    return 'http://localhost:8080';
  }

  // For Render production, use the backend service URL
  // Replace with your actual backend URL
  if (hostname.includes('front-prod')) {
    return 'https://ingsw3-back-prod.onrender.com';
  } else if (hostname.includes('front-qa')) {
    return 'https://ingsw3-back-qa.onrender.com';
  }

  // Fallback: use environment variable if set
  const envUrl = envVar || (process && process.env && process.env.REACT_APP_BACKEND_URL) || undefined;
  if (envUrl) {
    return envUrl;
  }

  // Final fallback
  return '';
};

const API_BASE_URL = getApiBaseUrl();
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
