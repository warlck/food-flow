const TOKEN_STORAGE_KEY = 'foodflow.admin-token';
const AUTH_API_BASE_URL = import.meta.env.VITE_AUTH_API_URL || '';

// NotAuthenticatedError is thrown by the API client when no valid token is
// available; the AuthProvider treats it as a signal to render the login page.
export class NotAuthenticatedError extends Error {
  constructor() {
    super('Not authenticated');
    this.name = 'NotAuthenticatedError';
  }
}

// login exchanges credentials for a JWT at the auth service and stores it.
// Error messages intentionally mirror the backend's generic responses so the
// UI never reveals which part of the credentials failed.
export async function login(email: string, password: string): Promise<string> {
  const res = await fetch(`${AUTH_API_BASE_URL}/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  if (res.status === 429) throw new Error('Too many attempts. Try again later.');
  if (res.status === 403) throw new Error('This account is not authorized for Restaurant Studio.');
  if (!res.ok) throw new Error('Invalid email or password');

  const { token } = (await res.json()) as { token: string };
  window.localStorage.setItem(TOKEN_STORAGE_KEY, token);
  return token;
}

export function logout(): void {
  window.localStorage.removeItem(TOKEN_STORAGE_KEY);
}

// getValidToken returns the stored token unless it is missing, malformed, or
// expired, in which case the store is cleared and null is returned.
export function getValidToken(): string | null {
  const token = window.localStorage.getItem(TOKEN_STORAGE_KEY);
  if (!token) return null;

  try {
    const { exp } = JSON.parse(window.atob(token.split('.')[1])) as { exp?: number };
    if (!exp || exp * 1000 <= Date.now()) {
      logout();
      return null;
    }
    return token;
  } catch {
    logout();
    return null;
  }
}
