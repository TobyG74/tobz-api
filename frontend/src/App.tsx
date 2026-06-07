import { useEffect, useState } from "react";
import { motion } from "motion/react";
import { toast } from "sonner";
import { Github, ShieldCheck, Zap } from "lucide-react";
import { useAuth } from "./context/AuthContext";
import { Nav } from "./components/Nav";
import { Downloader } from "./components/Downloader";
import { PlatformTicker } from "./components/PlatformTicker";
import { Features } from "./components/Features";
import { Footer } from "./components/Footer";
import { AuthModal } from "./components/AuthModal";
import { KeysPanel } from "./components/KeysPanel";
import { Kicker } from "./components/primitives";

const fadeUp = {
  hidden: { opacity: 0, y: 22 },
  show: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { duration: 0.6, delay: i * 0.1, ease: [0.22, 1, 0.36, 1] as const },
  }),
};

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

  return (
    <div id="top">
      <Nav onAuth={() => setAuthOpen(true)} onKeys={() => setKeysOpen(true)} />

      {/* ---------- Hero ---------- */}
      <section id="downloader" className="relative mx-auto max-w-6xl px-4 pb-10 pt-16 sm:pt-24">
        {/* decorative glow orb */}
        <div className="pointer-events-none absolute left-1/2 top-0 -z-10 h-72 w-72 -translate-x-1/2 rounded-full bg-violet/20 blur-[100px]" />

        <div className="mx-auto max-w-3xl text-center">
          <motion.div custom={0} variants={fadeUp} initial="hidden" animate="show" className="flex justify-center">
            <Kicker>REST API · v2 · Go + Fiber</Kicker>
          </motion.div>

          <motion.h1
            custom={1}
            variants={fadeUp}
            initial="hidden"
            animate="show"
            className="mt-6 font-display text-5xl leading-[1.05] sm:text-7xl"
          >
            Unduh apa saja.
            <br />
            <span className="text-gradient">Bangun apa saja.</span>
          </motion.h1>

          <motion.p
            custom={2}
            variants={fadeUp}
            initial="hidden"
            animate="show"
            className="mx-auto mt-6 max-w-xl text-lg text-fog"
          >
            Satu API untuk downloader media multi-platform + autentikasi modern, captcha,
            dan API key ber-kuota. Cepat, aman, fleksibel.
          </motion.p>

          <motion.div
            custom={3}
            variants={fadeUp}
            initial="hidden"
            animate="show"
            className="mt-5 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-sm text-mist"
          >
            <span className="inline-flex items-center gap-1.5"><Zap size={15} className="text-cyan" /> respon ringkas</span>
            <span className="inline-flex items-center gap-1.5"><ShieldCheck size={15} className="text-cyan" /> anti-SSRF & captcha</span>
            <span className="inline-flex items-center gap-1.5"><Github size={15} className="text-cyan" /> OAuth Google/GitHub</span>
          </motion.div>
        </div>

        <motion.div custom={4} variants={fadeUp} initial="hidden" animate="show" className="mt-12">
          <Downloader onNeedAuth={() => setAuthOpen(true)} onNeedKey={() => setKeysOpen(true)} />
        </motion.div>
      </section>

      <PlatformTicker />
      <Features />
      <Footer />

      <AuthModal open={authOpen} onClose={() => setAuthOpen(false)} />
      <KeysPanel open={keysOpen} onClose={() => setKeysOpen(false)} />
    </div>
  );
}
