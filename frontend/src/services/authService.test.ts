import axios from 'axios';
import { authService, getApiBaseUrl } from './authService';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('getApiBaseUrl', () => {
  test('returns localhost URL for localhost hostname', () => {
    expect(getApiBaseUrl('localhost')).toBe('http://localhost:8080');
  });

  test('returns production URL for front-prod hostname', () => {
    expect(getApiBaseUrl('myapp-front-prod.onrender.com')).toBe('https://ingsw3-back-prod.onrender.com');
  });

  test('returns QA URL for front-qa hostname', () => {
    expect(getApiBaseUrl('myapp-front-qa.onrender.com')).toBe('https://ingsw3-back-qa.onrender.com');
  });

  test('returns environment variable URL when set', () => {
    expect(getApiBaseUrl('unknown.com', 'https://custom-backend.com')).toBe('https://custom-backend.com');
  });

  test('returns empty string when no conditions match', () => {
    expect(getApiBaseUrl('unknown.com')).toBe('');
  });
});

describe('authService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('login', () => {
    test('calls API correctly with credentials', async () => {
      const mockUser = {
        id: 1,
        email: 'test@example.com',
        username: 'testuser',
        created_at: '2025-01-01'
      };

      mockedAxios.post.mockResolvedValueOnce({
        status: 200,
        statusText: 'OK',
        headers: {},
        config: { url: '' },
        data: mockUser
      });

      const result = await authService.login({
        email: 'test@example.com',
        password: '123456'
      });

      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8080/api/auth/login',
        {
          email: 'test@example.com',
          password: '123456'
        }
      );
      expect(result).toEqual(mockUser);
    });

    test('rejects when credentials are invalid', async () => {
      const error = new Error('Credenciales inválidas');
      mockedAxios.post.mockRejectedValueOnce(error);

      await expect(
        authService.login({
          email: 'wrong@example.com',
          password: 'wrong'
        })
      ).rejects.toEqual(error);
    });

    test('rejects when there is a network error', async () => {
      const error = new Error('Network Error');
      mockedAxios.post.mockRejectedValueOnce(error);

      await expect(
        authService.login({
          email: 'test@example.com',
          password: 'password'
        })
      ).rejects.toEqual(error);
    });
  });

  describe('register', () => {
    test('calls API correctly with registration data', async () => {
      const mockUser = {
        id: 1,
        email: 'newuser@example.com',
        username: 'newuser',
        created_at: '2025-01-01'
      };

      mockedAxios.post.mockResolvedValueOnce({
        status: 200,
        statusText: 'OK',
        headers: {},
        config: { url: '' },
        data: mockUser
      });

      const result = await authService.register({
        email: 'newuser@example.com',
        password: '123456',
        username: 'newuser'
      });

      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8080/api/auth/register',
        {
          email: 'newuser@example.com',
          password: '123456',
          username: 'newuser'
        }
      );
      expect(result).toEqual(mockUser);
    });

    test('rejects when email already exists', async () => {
      const error = new Error('el email ya está registrado');
      mockedAxios.post.mockRejectedValueOnce(error);

      await expect(
        authService.register({
          email: 'existing@example.com',
          password: '123456',
          username: 'testuser'
        })
      ).rejects.toEqual(error);
    });

    test('rejects when validation fails', async () => {
      const error = new Error('la contraseña debe tener al menos 6 caracteres');
      mockedAxios.post.mockRejectedValueOnce(error);

      await expect(
        authService.register({
          email: 'test@example.com',
          password: '123',
          username: 'testuser'
        })
      ).rejects.toEqual(error);
    });
  });
});
