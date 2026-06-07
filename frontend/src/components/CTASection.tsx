import { motion } from "motion/react";
import { ArrowRight } from "lucide-react";
import { Button } from "./primitives";

export function CTASection({ onStart }: { onStart: () => void }) {
  return (
    <section className="mx-auto max-w-6xl px-4 pb-24">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true, margin: "-80px" }}
        transition={{ duration: 0.5 }}
        className="grad-border relative overflow-hidden rounded-3xl px-8 py-16 text-center"
      >
        <div className="pointer-events-none absolute left-1/2 top-1/2 -z-10 h-64 w-[40rem] -translate-x-1/2 -translate-y-1/2 rounded-full bg-violet/20 blur-[100px]" />
        <h2 className="mx-auto max-w-2xl font-display text-4xl leading-tight sm:text-5xl">
          Siap bangun dengan <span className="text-gradient">Tobz API?</span>
        </h2>
        <p className="mx-auto mt-4 max-w-lg text-fog">
          Daftar gratis, ambil API key-mu, dan kirim request pertama dalam hitungan menit.
        </p>
        <div className="mt-8 flex justify-center">
          <Button onClick={onStart} className="px-8 py-3 text-base">
            Mulai gratis sekarang <ArrowRight size={18} />
          </Button>
        </div>
      </motion.div>
    </section>
  );
}
