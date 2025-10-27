import axios from 'axios';
import { User, LoginRequest, RegisterRequest } from '../types';

// Configure backend URLs for each environment
const getApiBaseUrl = () => {
  if (window.location.hostname === 'localhost') {
    return 'http://localhost:8080';
  }

  // For Render production, use the backend service URL
  // Replace with your actual backend URL
  const hostname = window.location.hostname;
  if (hostname.includes('front-prod')) {
    return 'https://ingsw3-back-prod.onrender.com';
  } else if (hostname.includes('front-qa')) {
    return 'https://ingsw3-back-qa.onrender.com';
  }

  // Fallback for other environments
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
