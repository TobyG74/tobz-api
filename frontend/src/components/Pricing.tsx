import { motion } from "motion/react";
import { Check } from "lucide-react";
import clsx from "clsx";
import { Button, Kicker } from "./primitives";

interface Plan {
  name: string;
  price: string;
  period: string;
  tagline: string;
  quota: string;
  features: string[];
  cta: string;
  featured?: boolean;
}

// Prices are placeholders (IDR) — adjust to your real pricing.
// Quotas map to the backend tiers: free=800, pro=50k, vvip=1jt requests/day.
const PLANS: Plan[] = [
  {
    name: "Free",
    price: "Rp0",
    period: "selamanya",
    tagline: "Untuk coba-coba & proyek kecil.",
    quota: "800 request / hari",
    features: ["Semua platform", "Rate limit standar", "Dukungan komunitas", "1 API key"],
    cta: "Mulai gratis",
  },
  {
    name: "Pro",
    price: "Rp49rb",
    period: "/ bulan",
    tagline: "Untuk produk yang sedang bertumbuh.",
    quota: "50.000 request / hari",
    features: ["Semua di Free", "Prioritas antrian", "Dukungan email", "5 API key", "Statistik penggunaan"],
    cta: "Pilih Pro",
    featured: true,
  },
  {
    name: "VVIP",
    price: "Rp199rb",
    period: "/ bulan",
    tagline: "Untuk skala besar & bisnis.",
    quota: "1.000.000 request / hari",
    features: ["Semua di Pro", "Rate limit tinggi", "Dukungan prioritas", "API key tak terbatas", "SLA 99.9%"],
    cta: "Pilih VVIP",
  },
];

export function Pricing({ onChoose }: { onChoose: () => void }) {
  return (
    <section id="pricing" className="mx-auto max-w-6xl px-4 py-24">
      <div className="mb-12 text-center">
        <div className="flex justify-center">
          <Kicker>Harga transparan</Kicker>
        </div>
        <h2 className="mt-4 font-display text-4xl sm:text-5xl">
          Bayar sesuai <span className="text-gradient">skala</span> kamu
        </h2>
        <p className="mx-auto mt-4 max-w-xl text-fog">
          Mulai gratis, upgrade kapan saja. Tanpa biaya tersembunyi.
        </p>
      </div>

      <div className="grid items-stretch gap-5 lg:grid-cols-3">
        {PLANS.map((plan, i) => (
          <motion.div
            key={plan.name}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.4, delay: i * 0.08 }}
            className={clsx(
              "relative flex flex-col rounded-2xl p-7",
              plan.featured
                ? "grad-border glow-violet"
                : "border border-line bg-surface/40"
            )}
          >
            {plan.featured && (
              <span className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-gradient-to-r from-blue to-violet px-3 py-1 font-mono text-[11px] font-600 uppercase tracking-wider text-white">
                Paling populer
              </span>
            )}

            <h3 className="font-display text-xl">{plan.name}</h3>
            <p className="mt-1 text-sm text-mist">{plan.tagline}</p>

            <div className="mt-5 flex items-end gap-1">
              <span className="font-display text-4xl font-700 text-paper">{plan.price}</span>
              <span className="mb-1 text-sm text-mist">{plan.period}</span>
            </div>

            <div className="mt-4 rounded-lg bg-void/50 px-3 py-2 text-center font-mono text-sm text-cyan ring-1 ring-line">
              {plan.quota}
            </div>

            <ul className="mt-6 flex-1 space-y-3 text-sm">
              {plan.features.map((f) => (
                <li key={f} className="flex items-center gap-2.5 text-fog">
                  <Check size={16} className="shrink-0 text-cyan" /> {f}
                </li>
              ))}
            </ul>

            <Button
              onClick={onChoose}
              variant={plan.featured ? "primary" : "outline"}
              className="mt-7 w-full py-3"
            >
              {plan.cta}
            </Button>
          </motion.div>
        ))}
      </div>

      <p className="mt-6 text-center font-mono text-xs text-mist">
        * Harga contoh dalam Rupiah — sesuaikan dengan harga aktualmu.
      </p>
    </section>
  );
}
