import { motion } from "motion/react";
import { Fingerprint, Gauge, Lock, ShieldCheck, ToggleRight, Boxes } from "lucide-react";
import type { ReactNode } from "react";
import { Kicker } from "./primitives";

const ITEMS: { icon: ReactNode; title: string; body: string }[] = [
  {
    icon: <Lock size={20} />,
    title: "Argon2id + JWT",
    body: "Password di-hash memory-hard, sesi pakai access token singkat + refresh-token rotation di cookie httpOnly.",
  },
  {
    icon: <ShieldCheck size={20} />,
    title: "Anti-SSRF",
    body: "Setiap fetch ke pihak ketiga memblokir IP privat/loopback & metadata cloud — tidak bisa dipakai menembus jaringan internal.",
  },
  {
    icon: <Fingerprint size={20} />,
    title: "Captcha + lockout",
    body: "Cloudflare Turnstile pada login/daftar, plus penguncian akun otomatis saat brute-force.",
  },
  {
    icon: <Gauge size={20} />,
    title: "Kuota & rate limit",
    body: "API key ber-tier dengan kuota harian atomik dan rate limiting per endpoint.",
  },
  {
    icon: <ToggleRight size={20} />,
    title: "OAuth Google & GitHub",
    body: "Login sosial dengan proteksi state anti-CSRF dan penautan akun via email terverifikasi.",
  },
  {
    icon: <Boxes size={20} />,
    title: "Arsitektur modular",
    body: "Fitur baru (search, lookup, games) cukup tambah satu folder — tidak menumpuk di paket pusat.",
  },
];

export function Features() {
  return (
    <section id="features" className="mx-auto max-w-6xl px-4 py-24">
      <div className="mb-12 max-w-2xl">
        <Kicker>Dibangun untuk produksi</Kicker>
        <h2 className="mt-4 font-display text-4xl leading-tight sm:text-5xl">
          Aman secara <span className="text-gradient">default</span>, cepat secara desain.
        </h2>
        <p className="mt-4 text-fog">
          Backend Go (Fiber + GORM) dengan hardening keamanan menyeluruh — bukan tempelan.
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {ITEMS.map((it, i) => (
          <motion.div
            key={it.title}
            initial={{ opacity: 0, y: 18 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.4, delay: i * 0.05 }}
            className="group rounded-2xl border border-line bg-surface/40 p-6 transition hover:border-violet/50 hover:bg-surface-2/50"
          >
            <div className="mb-4 grid h-11 w-11 place-items-center rounded-xl bg-gradient-to-br from-blue/20 to-violet/20 text-azure ring-1 ring-line transition group-hover:text-cyan">
              {it.icon}
            </div>
            <h3 className="font-display text-lg">{it.title}</h3>
            <p className="mt-2 text-sm leading-relaxed text-fog">{it.body}</p>
          </motion.div>
        ))}
      </div>
    </section>
  );
}
