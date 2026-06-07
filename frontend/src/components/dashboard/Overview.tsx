import { useEffect, useState } from "react";
import { motion } from "motion/react";
import { ArrowUpRight, CalendarDays, CheckCircle2, Download, Gauge, Globe, KeyRound, Search } from "lucide-react";
import type { ReactNode } from "react";
import { useAuth } from "../../context/AuthContext";
import { useActiveKey } from "../../hooks/useActiveKey";
import { api } from "../../lib/api";
import type { GroupId } from "../../lib/endpoints";

function fmtDate(s?: string) {
  if (!s) return "—";
  try {
    return new Date(s).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
  } catch {
    return "—";
  }
}

type Action = { kind: "group"; group: GroupId } | { kind: "keys" };

const FEATURES: { action: Action; icon: ReactNode; title: string; desc: string }[] = [
  {
    action: { kind: "group", group: "downloader" },
    icon: <Download size={20} />,
    title: "Downloader",
    desc: "Unduh media dari TikTok, YouTube, Instagram, Facebook, Twitter/X, Douyin.",
  },
  {
    action: { kind: "group", group: "search" },
    icon: <Search size={20} />,
    title: "Search",
    desc: "Cari gambar (Bing, DuckDuckGo, Pexels, Pinterest) & jelajahi Pixiv.",
  },
  {
    action: { kind: "keys" },
    icon: <KeyRound size={20} />,
    title: "API Keys",
    desc: "Buat & kelola API key dengan kuota harian untuk semua fitur.",
  },
];

const STEPS = [
  "Buat API key di tab API Keys (otomatis diaktifkan).",
  "Pilih fitur — Downloader atau Search.",
  "Masukkan input, jalankan, dapatkan hasil instan.",
];

export function Overview({
  openGroup,
  openKeys,
}: {
  openGroup: (g: GroupId) => void;
  openKeys: () => void;
}) {
  const { user } = useAuth();
  const { key } = useActiveKey();
  const run = (a: Action) => (a.kind === "keys" ? openKeys() : openGroup(a.group));

  const [stats, setStats] = useState({ keys: 0, used: 0, limit: 0, wl: 0 });
  useEffect(() => {
    (async () => {
      try {
        const [keys, wl] = await Promise.all([api.listKeys(), api.listWhitelist()]);
        const active = keys.filter((k) => !k.revoked);
        setStats({
          keys: active.length,
          used: active.reduce((a, k) => a + k.quota_used, 0),
          limit: active.reduce((a, k) => a + k.daily_quota, 0),
          wl: wl.ips.length,
        });
      } catch {
        /* ignore */
      }
    })();
  }, []);

  const pct = stats.limit ? Math.min(100, Math.round((stats.used / stats.limit) * 100)) : 0;
  const STATS: { icon: ReactNode; label: string; value: string; sub?: ReactNode }[] = [
    {
      icon: <Gauge size={16} />,
      label: "Pemakaian hari ini",
      value: `${stats.used.toLocaleString("id-ID")} / ${stats.limit.toLocaleString("id-ID")}`,
      sub: (
        <div className="mt-2 h-1.5 overflow-hidden rounded-full bg-void">
          <div className="h-full rounded-full bg-gradient-to-r from-blue to-violet" style={{ width: `${pct}%` }} />
        </div>
      ),
    },
    { icon: <KeyRound size={16} />, label: "API Keys aktif", value: String(stats.keys) },
    { icon: <Globe size={16} />, label: "IP Whitelist", value: `${stats.wl} / 5` },
    { icon: <CalendarDays size={16} />, label: "Bergabung", value: fmtDate(user?.created_at) },
  ];

  return (
    <div>
      <motion.div initial={{ opacity: 0, y: 14 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.4 }}>
        <h2 className="font-display text-3xl">
          Halo, <span className="text-gradient">{user?.display_name || "developer"}</span> 👋
        </h2>
        <p className="mt-1 text-fog">Selamat datang di dashboard. Pilih fitur untuk mulai.</p>
      </motion.div>

      {/* stats */}
      <div className="mt-6 grid grid-cols-2 gap-3 lg:grid-cols-4">
        {STATS.map((s, i) => (
          <motion.div
            key={s.label}
            initial={{ opacity: 0, y: 14 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.35, delay: i * 0.05 }}
            className="rounded-2xl border border-line bg-surface/40 p-4"
          >
            <div className="flex items-center gap-2 text-mist">
              <span className="text-azure">{s.icon}</span>
              <span className="text-xs uppercase tracking-widest">{s.label}</span>
            </div>
            <p className="mt-2 font-display text-xl text-paper">{s.value}</p>
            {s.sub}
          </motion.div>
        ))}
      </div>

      {/* active key banner */}
      <div className="mt-6 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-line bg-surface/40 p-4">
        <div className="flex items-center gap-3">
          <span className={`grid h-9 w-9 place-items-center rounded-lg ${key ? "bg-cyan/15 text-cyan" : "bg-violet/15 text-violet"}`}>
            {key ? <CheckCircle2 size={18} /> : <KeyRound size={18} />}
          </span>
          <div>
            <p className="text-sm font-500 text-paper">{key ? "API key aktif" : "Belum ada API key aktif"}</p>
            <p className="font-mono text-xs text-mist">
              {key ? key.slice(0, 14) + "••••••••" : "Buat satu untuk memakai fitur"}
            </p>
          </div>
        </div>
        <button
          onClick={openKeys}
          className="inline-flex items-center gap-1 text-sm font-500 text-azure hover:text-cyan"
        >
          Kelola keys <ArrowUpRight size={15} />
        </button>
      </div>

      {/* feature cards */}
      <div className="mt-6 grid gap-4 sm:grid-cols-3">
        {FEATURES.map((f, i) => (
          <motion.button
            key={f.title}
            onClick={() => run(f.action)}
            initial={{ opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: i * 0.06 }}
            className="group rounded-2xl border border-line bg-surface/40 p-6 text-left transition hover:border-violet/50 hover:bg-surface-2/50"
          >
            <div className="mb-4 grid h-11 w-11 place-items-center rounded-xl bg-gradient-to-br from-blue/20 to-violet/20 text-azure ring-1 ring-line transition group-hover:text-cyan">
              {f.icon}
            </div>
            <h3 className="flex items-center gap-1 font-display text-lg">
              {f.title}
              <ArrowUpRight size={16} className="text-mist transition group-hover:text-violet" />
            </h3>
            <p className="mt-1.5 text-sm leading-relaxed text-fog">{f.desc}</p>
          </motion.button>
        ))}
      </div>

      {/* how it works */}
      <div className="mt-6 rounded-2xl border border-line bg-surface/40 p-6">
        <h3 className="font-display text-lg">Cara kerja</h3>
        <div className="mt-4 grid gap-4 sm:grid-cols-3">
          {STEPS.map((s, i) => (
            <div key={i} className="relative rounded-xl bg-void/40 p-4">
              <span className="font-display text-3xl font-700 text-gradient">{i + 1}</span>
              <p className="mt-1 text-sm text-fog">{s}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
