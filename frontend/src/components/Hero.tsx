import { motion } from "motion/react";
import { ArrowRight, BookOpen, Check } from "lucide-react";
import { Button, Kicker } from "./primitives";
import { CodeWindow } from "./CodeWindow";

const fadeUp = {
  hidden: { opacity: 0, y: 22 },
  show: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { duration: 0.6, delay: i * 0.1, ease: [0.22, 1, 0.36, 1] as const },
  }),
};

const PERKS = ["800 request/hari gratis", "Tanpa kartu kredit", "Integrasi < 5 menit"];

export function Hero({ onStart }: { onStart: () => void }) {
  return (
    <section className="relative mx-auto max-w-6xl px-4 pb-12 pt-16 sm:pt-24">
      <div className="pointer-events-none absolute left-1/4 top-0 -z-10 h-72 w-72 -translate-x-1/2 rounded-full bg-violet/20 blur-[110px]" />
      <div className="pointer-events-none absolute right-0 top-20 -z-10 h-64 w-64 rounded-full bg-blue/20 blur-[110px]" />

      <div className="grid items-center gap-12 lg:grid-cols-[1.05fr_1fr]">
        {/* Copy */}
        <div>
          <motion.div custom={0} variants={fadeUp} initial="hidden" animate="show">
            <Kicker>REST API untuk developer</Kicker>
          </motion.div>

          <motion.h1
            custom={1}
            variants={fadeUp}
            initial="hidden"
            animate="show"
            className="mt-6 font-display text-5xl leading-[1.05] sm:text-6xl"
          >
            Satu API untuk
            <br />
            <span className="text-gradient">unduh media apa saja.</span>
          </motion.h1>

          <motion.p
            custom={2}
            variants={fadeUp}
            initial="hidden"
            animate="show"
            className="mt-6 max-w-xl text-lg text-fog"
          >
            TikTok, YouTube, Instagram, dan lainnya — satu endpoint, respons JSON ringkas,
            uptime tinggi. Daftar, ambil API key, langsung integrasikan ke produkmu.
          </motion.p>

          <motion.div
            custom={3}
            variants={fadeUp}
            initial="hidden"
            animate="show"
            className="mt-8 flex flex-wrap items-center gap-3"
          >
            <Button onClick={onStart} className="px-7 py-3 text-base">
              Mulai gratis <ArrowRight size={18} />
            </Button>
            <a href="#pricing">
              <Button variant="outline" className="px-6 py-3 text-base">
                <BookOpen size={17} /> Lihat harga
              </Button>
            </a>
          </motion.div>

          <motion.ul
            custom={4}
            variants={fadeUp}
            initial="hidden"
            animate="show"
            className="mt-7 flex flex-wrap gap-x-6 gap-y-2 text-sm text-mist"
          >
            {PERKS.map((p) => (
              <li key={p} className="inline-flex items-center gap-1.5">
                <Check size={15} className="text-cyan" /> {p}
              </li>
            ))}
          </motion.ul>
        </div>

        {/* Code mockup */}
        <CodeWindow />
      </div>
    </section>
  );
}
