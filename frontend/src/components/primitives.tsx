import { AnimatePresence, motion } from "motion/react";
import { X } from "lucide-react";
import clsx from "clsx";
import type { ButtonHTMLAttributes, ReactNode } from "react";

/* ---- Logo ---- */
export function Logo({ className }: { className?: string }) {
  return (
    <a href="#top" className={clsx("flex items-center gap-2.5 group", className)}>
      <img
        src="/logo.png"
        alt="Tobz API"
        className="h-9 w-9 rounded-lg object-cover ring-1 ring-line transition-transform group-hover:scale-105"
      />
      <span className="font-display text-lg font-700 tracking-tight">
        Tobz<span className="text-gradient">API</span>
      </span>
    </a>
  );
}

/* ---- Button ---- */
type Variant = "primary" | "ghost" | "outline" | "danger";
interface BtnProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  loading?: boolean;
}
export function Button({ variant = "primary", loading, className, children, disabled, ...rest }: BtnProps) {
  const base =
    "inline-flex items-center justify-center gap-2 rounded-xl px-5 py-2.5 text-sm font-600 font-display tracking-wide transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed select-none";
  const variants: Record<Variant, string> = {
    primary:
      "text-white bg-gradient-to-r from-blue to-violet hover:shadow-[0_0_28px_-4px_rgba(139,92,255,0.7)] active:scale-[0.98]",
    outline: "border border-line text-paper hover:border-violet/60 hover:bg-surface-2/50",
    ghost: "text-fog hover:text-paper hover:bg-surface-2/60",
    danger: "border border-red-500/40 text-red-300 hover:bg-red-500/10",
  };
  return (
    <button className={clsx(base, variants[variant], className)} disabled={disabled || loading} {...rest}>
      {loading && (
        <span className="h-3.5 w-3.5 animate-spin rounded-full border-2 border-white/30 border-t-white" />
      )}
      {children}
    </button>
  );
}

/* ---- Field ---- */
export function Field({
  label,
  hint,
  children,
}: {
  label: string;
  hint?: string;
  children: ReactNode;
}) {
  return (
    <label className="block">
      <span className="mb-1.5 flex items-center justify-between text-xs font-500 uppercase tracking-widest text-mist">
        {label}
        {hint && <span className="font-mono normal-case tracking-normal text-mist/70">{hint}</span>}
      </span>
      {children}
    </label>
  );
}

export const inputCls =
  "w-full rounded-xl border border-line bg-void/60 px-4 py-3 text-paper placeholder:text-mist/50 outline-none transition focus:border-violet/70 focus:ring-2 focus:ring-violet/20 font-body";

/* ---- Section label (mono kicker) ---- */
export function Kicker({ children }: { children: ReactNode }) {
  return (
    <span className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-[0.25em] text-azure">
      <span className="h-px w-6 bg-gradient-to-r from-transparent to-azure" />
      {children}
    </span>
  );
}

/* ---- Modal ---- */
export function Modal({
  open,
  onClose,
  children,
  title,
}: {
  open: boolean;
  onClose: () => void;
  children: ReactNode;
  title?: string;
}) {
  return (
    <AnimatePresence>
      {open && (
        <motion.div
          className="fixed inset-0 z-[60] flex items-center justify-center p-4"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
        >
          <div className="absolute inset-0 bg-void/80 backdrop-blur-sm" onClick={onClose} />
          <motion.div
            className="grad-border relative w-full max-w-md rounded-2xl p-6 shadow-2xl"
            initial={{ y: 24, scale: 0.96, opacity: 0 }}
            animate={{ y: 0, scale: 1, opacity: 1 }}
            exit={{ y: 16, scale: 0.97, opacity: 0 }}
            transition={{ type: "spring", stiffness: 320, damping: 28 }}
          >
            <div className="mb-5 flex items-center justify-between">
              {title && <h3 className="font-display text-xl">{title}</h3>}
              <button
                onClick={onClose}
                className="ml-auto rounded-lg p-1.5 text-mist hover:bg-surface-2 hover:text-paper"
              >
                <X size={18} />
              </button>
            </div>
            {children}
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}

/* ---- Sheet (right slide-over) ---- */
export function Sheet({
  open,
  onClose,
  children,
  title,
}: {
  open: boolean;
  onClose: () => void;
  children: ReactNode;
  title?: string;
}) {
  return (
    <AnimatePresence>
      {open && (
        <motion.div
          className="fixed inset-0 z-[60]"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
        >
          <div className="absolute inset-0 bg-void/70 backdrop-blur-sm" onClick={onClose} />
          <motion.aside
            className="glass absolute right-0 top-0 flex h-full w-full max-w-lg flex-col p-6"
            initial={{ x: "100%" }}
            animate={{ x: 0 }}
            exit={{ x: "100%" }}
            transition={{ type: "spring", stiffness: 280, damping: 32 }}
          >
            <div className="mb-6 flex items-center justify-between">
              <h3 className="font-display text-2xl">{title}</h3>
              <button
                onClick={onClose}
                className="rounded-lg p-2 text-mist hover:bg-surface-2 hover:text-paper"
              >
                <X size={20} />
              </button>
            </div>
            <div className="flex-1 overflow-y-auto pr-1">{children}</div>
          </motion.aside>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
