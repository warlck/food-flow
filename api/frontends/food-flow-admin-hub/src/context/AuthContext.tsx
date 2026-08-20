import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';
import { getValidToken, login as authLogin, logout as authLogout } from '@/lib/auth';
import { adminApi } from '@/lib/admin-api';

interface AuthContextValue {
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => getValidToken());

  // Any 401 from the API client clears the session and returns to the
  // login page. 403s surface as error toasts instead — never a logout loop.
  useEffect(() => {
    adminApi.setUnauthorizedHandler(() => setToken(null));
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const newToken = await authLogin(email, password);
    setToken(newToken);
  }, []);

  const logout = useCallback(() => {
    authLogout();
    setToken(null);
  }, []);

  const value = useMemo(() => ({ token, login, logout }), [token, login, logout]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within an AuthProvider');
  return ctx;
}
