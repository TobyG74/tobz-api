import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { api, setAccessToken } from "../lib/api";
import type { AuthResult, User } from "../lib/types";

interface AuthCtx {
  user: User | null;
  loading: boolean;
  setSession: (r: AuthResult) => void;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const Ctx = createContext<AuthCtx | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const setSession = useCallback((r: AuthResult) => {
    setAccessToken(r.access_token);
    setUser(r.user);
  }, []);

  const logout = useCallback(async () => {
    try {
      await api.logout();
    } catch {
      /* ignore */
    }
    setAccessToken(null);
    setUser(null);
  }, []);

  const refreshUser = useCallback(async () => {
    const u = await api.me();
    setUser(u);
  }, []);

  // On load, try to silently restore a session via the refresh cookie.
  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const r = await api.refresh();
        if (active) setSession(r);
      } catch {
        /* no active session */
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [setSession]);

  const value = useMemo(
    () => ({ user, loading, setSession, logout, refreshUser }),
    [user, loading, setSession, logout, refreshUser]
  );

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useAuth() {
  const v = useContext(Ctx);
  if (!v) throw new Error("useAuth must be used within AuthProvider");
  return v;
}
