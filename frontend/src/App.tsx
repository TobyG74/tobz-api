import { useEffect, useState } from "react";
import { toast } from "sonner";
import { useAuth } from "./context/AuthContext";
import { Nav } from "./components/Nav";
import { Hero } from "./components/Hero";
import { Stats } from "./components/Stats";
import { PlatformTicker } from "./components/PlatformTicker";
import { Features } from "./components/Features";
import { HowItWorks } from "./components/HowItWorks";
import { Pricing } from "./components/Pricing";
import { FAQ } from "./components/FAQ";
import { CTASection } from "./components/CTASection";
import { Footer } from "./components/Footer";
import { AuthModal } from "./components/AuthModal";
import { KeysPanel } from "./components/KeysPanel";

export default function App() {
  const { loading } = useAuth();
  const [authOpen, setAuthOpen] = useState(false);
  const [keysOpen, setKeysOpen] = useState(false);

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

  const openAuth = () => setAuthOpen(true);

  return (
    <div id="top">
      <Nav onAuth={openAuth} onKeys={() => setKeysOpen(true)} />

      <Hero onStart={openAuth} />
      <Stats />
      <PlatformTicker />
      <Features />
      <HowItWorks />
      <Pricing onChoose={openAuth} />
      <FAQ />
      <CTASection onStart={openAuth} />
      <Footer />

      <AuthModal open={authOpen} onClose={() => setAuthOpen(false)} />
      <KeysPanel open={keysOpen} onClose={() => setKeysOpen(false)} />
    </div>
  );
}
