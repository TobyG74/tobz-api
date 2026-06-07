const STATS = [
  { value: "6+", label: "Platform didukung" },
  { value: "99.9%", label: "Uptime" },
  { value: "~150ms", label: "Latensi rata-rata" },
  { value: "1 jt+", label: "Request / hari (VVIP)" },
];

export function Stats() {
  return (
    <section className="border-y border-line bg-surface/30">
      <div className="mx-auto grid max-w-6xl grid-cols-2 gap-px sm:grid-cols-4">
        {STATS.map((s) => (
          <div key={s.label} className="px-4 py-8 text-center">
            <div className="font-display text-3xl font-700 text-gradient sm:text-4xl">{s.value}</div>
            <div className="mt-1 text-xs uppercase tracking-widest text-mist">{s.label}</div>
          </div>
        ))}
      </div>
    </section>
  );
}
