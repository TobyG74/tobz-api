import { KeyRound, LogOut, Sparkles } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { Button, Logo } from "./primitives";

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
          <a href="#downloader" className="transition hover:text-paper">Downloader</a>
          <a href="#platforms" className="transition hover:text-paper">Platform</a>
          <a href="#features" className="transition hover:text-paper">Keamanan</a>
        </nav>

        <div className="flex items-center gap-2">
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
            <Button onClick={onAuth}>
              <Sparkles size={16} /> Mulai
            </Button>
          )}
        </div>
      </div>
    </header>
  );
}
