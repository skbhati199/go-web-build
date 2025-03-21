/**
 * API utility functions for making HTTP requests
 */

const API_BASE_URL = process.env.REACT_APP_API_URL || 'https://api.example.com';

export const fetchData = async (endpoint, options = {}) => {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
};

export const get = (endpoint, options = {}) => {
  return fetchData(endpoint, { method: 'GET', ...options });
};

export const post = (endpoint, data, options = {}) => {
  return fetchData(endpoint, {
    method: 'POST',
    body: JSON.stringify(data),
    ...options,
  });
};