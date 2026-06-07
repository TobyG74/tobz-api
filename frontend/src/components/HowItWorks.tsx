import { motion } from "motion/react";
import { KeyRound, Rocket, UserPlus } from "lucide-react";
import type { ReactNode } from "react";
import { Kicker } from "./primitives";

const STEPS: { icon: ReactNode; title: string; body: string; code?: string }[] = [
  {
    icon: <UserPlus size={20} />,
    title: "1. Daftar akun",
    body: "Buat akun via email atau login Google/GitHub. Verifikasi captcha, selesai dalam hitungan detik.",
  },
  {
    icon: <KeyRound size={20} />,
    title: "2. Ambil API key",
    body: "Generate API key dari dashboard. Tier gratis langsung aktif dengan 800 request/hari.",
    code: "X-API-Key: tobz_live_••••",
  },
  {
    icon: <Rocket size={20} />,
    title: "3. Panggil API",
    body: "Kirim URL ke satu endpoint, terima tautan unduhan dalam JSON. Siap dipakai di produkmu.",
    code: "GET /v1/download?url=…",
  },
];

export function HowItWorks() {
  return (
    <section id="how" className="mx-auto max-w-6xl px-4 py-24">
      <div className="mb-12 text-center">
        <div className="flex justify-center">
          <Kicker>Integrasi kilat</Kicker>
        </div>
        <h2 className="mt-4 font-display text-4xl sm:text-5xl">
          Dari nol ke <span className="text-gradient">produksi</span> dalam 3 langkah
        </h2>
      </div>

      <div className="relative grid gap-5 md:grid-cols-3">
        {/* connecting line */}
        <div className="pointer-events-none absolute left-0 right-0 top-[2.6rem] hidden h-px bg-gradient-to-r from-transparent via-line to-transparent md:block" />
        {STEPS.map((s, i) => (
          <motion.div
            key={s.title}
            initial={{ opacity: 0, y: 18 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.4, delay: i * 0.1 }}
            className="relative rounded-2xl border border-line bg-surface/40 p-6"
          >
            <div className="mb-4 grid h-11 w-11 place-items-center rounded-xl bg-gradient-to-br from-blue to-violet text-white">
              {s.icon}
            </div>
            <h3 className="font-display text-lg">{s.title}</h3>
            <p className="mt-2 text-sm leading-relaxed text-fog">{s.body}</p>
            {s.code && (
              <code className="mt-4 block rounded-lg bg-void/60 px-3 py-2 font-mono text-xs text-cyan ring-1 ring-line">
                {s.code}
              </code>
            )}
          </motion.div>
        ))}
      </div>
    </section>
  );
}
