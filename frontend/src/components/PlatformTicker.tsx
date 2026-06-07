const PLATFORMS = [
  "TikTok",
  "YouTube",
  "Instagram",
  "Facebook",
  "Twitter / X",
  "Douyin",
  "Threads",
  "Spotify",
  "SoundCloud",
];

export function PlatformTicker() {
  const row = [...PLATFORMS, ...PLATFORMS];
  return (
    <section id="platforms" className="relative overflow-hidden border-y border-line py-6">
      <div className="pointer-events-none absolute inset-y-0 left-0 z-10 w-24 bg-gradient-to-r from-canvas to-transparent" />
      <div className="pointer-events-none absolute inset-y-0 right-0 z-10 w-24 bg-gradient-to-l from-canvas to-transparent" />
      <div className="flex w-max marquee gap-10 whitespace-nowrap">
        {row.map((p, i) => (
          <span
            key={i}
            className="flex items-center gap-10 font-display text-2xl font-600 text-fog/40"
          >
            {p}
            <span className="text-violet">✦</span>
          </span>
        ))}
      </div>
    </section>
  );
}
