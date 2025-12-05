const PROD_BACKEND_URL = 'https://api.hubplanner-ai.click';
const DEV_BACKEND_URL = 'http://localhost:8080';

// In a production build (import.meta.env.PROD is true), use the production URL.
// Otherwise, use the VITE_API_BASE_URL from the environment file, or fall back to the dev URL.
export const BACKEND_URL = import.meta.env.PROD
  ? PROD_BACKEND_URL
  : import.meta.env.VITE_BACKEND_URL || DEV_BACKEND_URL;

export const API_BASE_URL = `${BACKEND_URL}/api/v1`;
