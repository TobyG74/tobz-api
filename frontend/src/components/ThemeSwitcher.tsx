import { useEffect, useRef, useState } from "react";
import { AnimatePresence, motion } from "motion/react";
import { Check, Palette } from "lucide-react";
import { useTheme } from "../context/ThemeContext";

export function ThemeSwitcher() {
  const { theme, setTheme, themes } = useTheme();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  // Close on outside click / Escape.
  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => e.key === "Escape" && setOpen(false);
    document.addEventListener("mousedown", onClick);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onClick);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  return (
    <div ref={ref} className="relative">
      <button
        onClick={() => setOpen((o) => !o)}
        aria-label="Ganti tema"
        className="grid h-10 w-10 place-items-center rounded-xl border border-line text-fog transition hover:border-violet/60 hover:text-paper"
      >
        <Palette size={17} />
      </button>

      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: -6, scale: 0.97 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -6, scale: 0.97 }}
            transition={{ duration: 0.16 }}
            className="glass absolute right-0 z-50 mt-2 w-48 rounded-xl p-1.5 shadow-2xl"
          >
            <p className="px-2.5 py-1.5 font-mono text-[11px] uppercase tracking-widest text-mist">
              Tema
            </p>
            {themes.map((t) => (
              <button
                key={t.id}
                onClick={() => {
                  setTheme(t.id);
                  setOpen(false);
                }}
                className="flex w-full items-center gap-3 rounded-lg px-2.5 py-2 text-sm text-fog transition hover:bg-surface-2/70 hover:text-paper"
              >
                <span
                  className="h-5 w-5 shrink-0 rounded-md ring-1 ring-line"
                  style={{ background: `linear-gradient(135deg, ${t.swatch[0]}, ${t.swatch[1]})` }}
                />
                <span className="flex-1 text-left">{t.name}</span>
                {theme === t.id && <Check size={15} className="text-cyan" />}
              </button>
            ))}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
