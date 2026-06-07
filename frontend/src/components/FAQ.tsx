import { useState } from "react";
import { AnimatePresence, motion } from "motion/react";
import { ChevronDown } from "lucide-react";
import { Kicker } from "./primitives";

const QA = [
  {
    q: "Apa itu Tobz API?",
    a: "REST API untuk mengunduh media dari berbagai platform (TikTok, YouTube, Instagram, Facebook, Twitter/X, Douyin) lewat satu endpoint sederhana dengan respons JSON.",
  },
  {
    q: "Bagaimana cara autentikasinya?",
    a: "Setiap permintaan ke endpoint downloader memakai header X-API-Key. Buat akun, generate API key dari dashboard, lalu sertakan di setiap request.",
  },
  {
    q: "Apakah ada paket gratis?",
    a: "Ya. Paket Free memberi 800 request per hari selamanya, tanpa kartu kredit. Upgrade ke Pro atau VVIP kapan saja saat butuh kuota lebih besar.",
  },
  {
    q: "Bagaimana jika kuota harian habis?",
    a: "Permintaan akan menerima status 429 (rate limited) hingga kuota harian ter-reset otomatis. Tingkatkan paket untuk kuota lebih tinggi.",
  },
  {
    q: "Apakah aman dipakai di produksi?",
    a: "Backend dibangun dengan hardening menyeluruh: hashing Argon2id, JWT + refresh-token rotation, proteksi anti-SSRF, rate limiting, dan security header. Cocok untuk produksi.",
  },
];

function Item({ q, a }: { q: string; a: string }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="border border-line bg-surface/40 first:rounded-t-2xl last:rounded-b-2xl [&:not(:last-child)]:border-b-0">
      <button
        onClick={() => setOpen((o) => !o)}
        className="flex w-full items-center justify-between gap-4 px-5 py-4 text-left"
      >
        <span className="font-500 text-paper">{q}</span>
        <ChevronDown
          size={18}
          className={`shrink-0 text-mist transition-transform ${open ? "rotate-180 text-violet" : ""}`}
        />
      </button>
      <AnimatePresence initial={false}>
        {open && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25 }}
            className="overflow-hidden"
          >
            <p className="px-5 pb-4 text-sm leading-relaxed text-fog">{a}</p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

export function FAQ() {
  return (
    <section id="faq" className="mx-auto max-w-3xl px-4 py-24">
      <div className="mb-10 text-center">
        <div className="flex justify-center">
          <Kicker>Pertanyaan umum</Kicker>
        </div>
        <h2 className="mt-4 font-display text-4xl sm:text-5xl">FAQ</h2>
      </div>
      <div className="rounded-2xl">
        {QA.map((item) => (
          <Item key={item.q} {...item} />
        ))}
      </div>
    </section>
  );
}
