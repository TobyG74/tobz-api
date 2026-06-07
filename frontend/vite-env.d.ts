/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE: string;
  readonly VITE_TURNSTILE_SITEKEY: string;
}
interface ImportMeta {
  readonly env: ImportMetaEnv;
}

// Cloudflare Turnstile global (loaded via script tag in index.html).
interface Window {
  turnstile?: {
    render: (
      el: HTMLElement,
      opts: {
        sitekey: string;
        theme?: "light" | "dark" | "auto";
        callback?: (token: string) => void;
        "expired-callback"?: () => void;
        "error-callback"?: () => void;
      }
    ) => string;
    reset: (id?: string) => void;
    remove: (id?: string) => void;
  };
}
