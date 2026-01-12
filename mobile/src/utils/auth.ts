import AsyncStorage from '@react-native-async-storage/async-storage';
import {API_URL} from '../config';

const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

interface AuthResponse {
  access_token?: string;
  refresh_token?: string;
}

interface FetchOptions extends RequestInit {
  headers?: Record<string, string>;
}

// Store tokens
export const setTokens = async (
  accessToken: string,
  refreshToken: string | null = null,
): Promise<void> => {
  try {
    await AsyncStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    if (refreshToken) {
      await AsyncStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
    }
  } catch (error) {
    console.error('Error storing tokens:', error);
  }
};

// Get access token
export const getAccessToken = async (): Promise<string | null> => {
  try {
    return await AsyncStorage.getItem(ACCESS_TOKEN_KEY);
  } catch (error) {
    console.error('Error getting access token:', error);
    return null;
  }
};

// Get refresh token
export const getRefreshToken = async (): Promise<string | null> => {
  try {
    return await AsyncStorage.getItem(REFRESH_TOKEN_KEY);
  } catch (error) {
    console.error('Error getting refresh token:', error);
    return null;
  }
};

// Clear tokens
export const clearTokens = async (): Promise<void> => {
  try {
    await AsyncStorage.removeItem(ACCESS_TOKEN_KEY);
    await AsyncStorage.removeItem(REFRESH_TOKEN_KEY);
  } catch (error) {
    console.error('Error clearing tokens:', error);
  }
};

// Refresh or initialize access token
export const auth = async (): Promise<boolean> => {
  try {
    // Try to refresh access token first
    const refreshToken = await getRefreshToken();
    if (refreshToken) {
      const response = await fetch(`${API_URL}v1/auth/refreshaccess`, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${refreshToken}`,
        },
      });

      if (response.ok) {
        const data: AuthResponse = await response.json();
        if (data.access_token) {
          await setTokens(data.access_token, data.refresh_token || refreshToken);
          return true;
        }
      }
    }

    // If refresh fails, initialize new access
    const initResponse = await fetch(`${API_URL}v1/auth/initaccess`, {
      method: 'GET',
    });

    if (initResponse.ok) {
      const data: AuthResponse = await initResponse.json();
      if (data.access_token) {
        await setTokens(data.access_token, data.refresh_token || null);
        return true;
      }
    }

    return false;
  } catch (error) {
    console.error('Auth error:', error);
    return false;
  }
};

// Fetch with authentication
export const fetchWithAuth = async (
  endpoint: string,
  options: FetchOptions = {},
): Promise<Response> => {
  const makeRequest = async (): Promise<Response> => {
    const accessToken = await getAccessToken();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    };

    if (accessToken) {
      headers.Authorization = `Bearer ${accessToken}`;
    }

    return fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers,
    });
  };

  let response = await makeRequest();

  // If unauthorized, try to refresh and retry
  if (response.status === 401) {
    const refreshed = await auth();
    if (refreshed) {
      response = await makeRequest();
    }
  }

  return response;
};
