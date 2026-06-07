import { useEffect } from "react";
import { toast } from "sonner";
import { useAuth } from "./context/AuthContext";
import { Landing } from "./components/Landing";
import { Dashboard } from "./components/dashboard/Dashboard";

export default function App() {
  const { user, loading } = useAuth();

  // Handle the OAuth redirect result, then clean the URL.
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const ok = params.get("login");
    const err = params.get("auth_error");
    if (!ok && !err) return;
    if (ok === "success") toast.success("Login berhasil 🎉");
    if (err) {
      const msg: Record<string, string> = {
        cancelled: "Login dibatalkan",
        invalid_state: "Sesi OAuth tidak valid, coba lagi",
        exchange_failed: "Gagal autentikasi dengan provider",
        provider_unsupported: "Provider tidak didukung",
      };
      toast.error(msg[err] ?? "Login gagal");
    }
    window.history.replaceState({}, "", window.location.pathname);
  }, []);

  if (loading) {
    return (
      <div className="grid min-h-screen place-items-center">
        <div className="flex flex-col items-center gap-4">
          <img src="/logo.png" alt="" className="h-16 w-16 rounded-2xl pulse-ring" />
          <span className="font-mono text-xs uppercase tracking-[0.3em] text-mist">memuat…</span>
        </div>
      </div>
    );
  }

  return user ? <Dashboard /> : <Landing />;
}
