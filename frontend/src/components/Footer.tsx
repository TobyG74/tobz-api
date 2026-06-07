import { Logo } from "./primitives";

export function Footer() {
  return (
    <footer className="border-t border-line">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-4 py-10 sm:flex-row">
        <Logo />
        <p className="font-mono text-xs text-mist">
          © {new Date().getFullYear()} Tobz API · Go + Fiber · dibuat dengan ⚡
        </p>
        <div className="flex gap-5 text-sm text-fog">
          <a href="#downloader" className="hover:text-paper">Downloader</a>
          <a href="#features" className="hover:text-paper">Keamanan</a>
        </div>
      </div>
    </footer>
  );
}
