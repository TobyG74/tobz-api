import { useCallback, useEffect, useState } from "react";

const STORAGE = "tobz_active_key";

/**
 * Holds the raw API key the downloader uses. It is shown only once on creation,
 * so we keep it in localStorage for convenience (it is a service key, not the
 * auth session token, which stays in memory).
 */
export function useActiveKey() {
  const [key, setKey] = useState<string>(() => localStorage.getItem(STORAGE) ?? "");

  useEffect(() => {
    if (key) localStorage.setItem(STORAGE, key);
    else localStorage.removeItem(STORAGE);
  }, [key]);

  const clear = useCallback(() => setKey(""), []);
  return { key, setKey, clear };
}
