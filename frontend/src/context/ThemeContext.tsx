import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

export interface ThemeDef {
  id: string;
  name: string;
  /** Two swatch colors used in the theme picker. */
  swatch: [string, string];
  light?: boolean;
}

export const THEMES: ThemeDef[] = [
  { id: "cyber", name: "Cyber", swatch: ["#3b6bff", "#8b5cff"] },
  { id: "aurora", name: "Aurora", swatch: ["#2f5fff", "#7b3ff2"], light: true },
  { id: "sunset", name: "Sunset", swatch: ["#ff7a18", "#ff2e63"] },
  { id: "emerald", name: "Emerald", swatch: ["#10b981", "#a3e635"] },
];

const STORAGE = "tobz_theme";
const META = { cyber: "#070b18", aurora: "#f6f8fe", sunset: "#160a08", emerald: "#061410" } as Record<string, string>;

interface ThemeCtx {
  theme: string;
  setTheme: (id: string) => void;
  themes: ThemeDef[];
}

const Ctx = createContext<ThemeCtx | null>(null);

function initialTheme(): string {
  const saved = localStorage.getItem(STORAGE);
  if (saved && THEMES.some((t) => t.id === saved)) return saved;
  return "cyber";
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<string>(initialTheme);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem(STORAGE, theme);
    const meta = document.querySelector('meta[name="theme-color"]');
    if (meta && META[theme]) meta.setAttribute("content", META[theme]);
  }, [theme]);

  const setTheme = useCallback((id: string) => setThemeState(id), []);

  const value = useMemo(() => ({ theme, setTheme, themes: THEMES }), [theme, setTheme]);
  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useTheme() {
  const v = useContext(Ctx);
  if (!v) throw new Error("useTheme must be used within ThemeProvider");
  return v;
}
