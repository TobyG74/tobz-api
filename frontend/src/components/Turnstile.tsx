import { useEffect, useRef } from "react";

const SITEKEY = import.meta.env.VITE_TURNSTILE_SITEKEY ?? "1x00000000000000000000AA";

/** Renders the Cloudflare Turnstile widget and reports the token upward. */
export function Turnstile({ onToken }: { onToken: (token: string) => void }) {
  const ref = useRef<HTMLDivElement>(null);
  const widgetId = useRef<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    const tryRender = () => {
      if (cancelled || !ref.current || !window.turnstile) return false;
      if (widgetId.current) return true;
      widgetId.current = window.turnstile.render(ref.current, {
        sitekey: SITEKEY,
        theme: "dark",
        callback: (token) => onToken(token),
        "expired-callback": () => onToken(""),
        "error-callback": () => onToken(""),
      });
      return true;
    };

    if (!tryRender()) {
      const iv = setInterval(() => {
        if (tryRender()) clearInterval(iv);
      }, 200);
      return () => {
        cancelled = true;
        clearInterval(iv);
      };
    }
    return () => {
      cancelled = true;
    };
  }, [onToken]);

  return <div ref={ref} className="flex min-h-[65px] justify-center" />;
}
