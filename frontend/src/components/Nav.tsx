import { KeyRound, LogOut, Sparkles } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { Button, Logo } from "./primitives";
import { ThemeSwitcher } from "./ThemeSwitcher";

export function Nav({
  onAuth,
  onKeys,
}: {
  onAuth: () => void;
  onKeys: () => void;
}) {
  const { user, logout } = useAuth();

  return (
    <header className="sticky top-0 z-40">
      <div className="mx-auto mt-4 flex max-w-6xl items-center justify-between gap-4 rounded-2xl glass px-4 py-3 sm:px-6">
        <Logo />

        <nav className="hidden items-center gap-7 text-sm text-fog md:flex">
          <a href="#features" className="transition hover:text-paper">Fitur</a>
          <a href="#how" className="transition hover:text-paper">Cara Kerja</a>
          <a href="#pricing" className="transition hover:text-paper">Harga</a>
          <a href="#faq" className="transition hover:text-paper">FAQ</a>
        </nav>

        <div className="flex items-center gap-2">
          <ThemeSwitcher />
          {user ? (
            <>
              <Button variant="outline" onClick={onKeys} className="hidden sm:inline-flex">
                <KeyRound size={16} /> API Keys
              </Button>
              <div className="flex items-center gap-2 rounded-xl border border-line bg-surface/60 py-1.5 pl-1.5 pr-3">
                <span className="grid h-7 w-7 place-items-center rounded-lg bg-gradient-to-br from-blue to-violet text-xs font-700 text-white">
                  {user.display_name.slice(0, 1).toUpperCase()}
                </span>
                <span className="max-w-[9rem] truncate text-sm text-paper">{user.display_name}</span>
              </div>
              <Button variant="ghost" onClick={logout} aria-label="Keluar">
                <LogOut size={16} />
              </Button>
            </>
          ) : (
            <>
              <button
                onClick={onAuth}
                className="hidden text-sm text-fog transition hover:text-paper sm:block"
              >
                Masuk
              </button>
              <Button onClick={onAuth}>
                <Sparkles size={16} /> Mulai gratis
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
