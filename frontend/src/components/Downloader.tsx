import { useState } from "react";
import { motion } from "motion/react";
import { toast } from "sonner";
import { ArrowRight, KeyRound, Link2, Loader2 } from "lucide-react";
import { api, ApiError } from "../lib/api";
import type { MediaResult } from "../lib/types";
import { useActiveKey } from "../hooks/useActiveKey";
import { useAuth } from "../context/AuthContext";
import { MediaResultCard } from "./MediaResultCard";
import { Button, inputCls } from "./primitives";

const SAMPLES = [
  "https://www.tiktok.com/@scout2015/video/6718335390845095173",
  "https://youtu.be/dQw4w9WgXcQ",
  "https://www.instagram.com/p/CXY/",
];

export function Downloader({
  onNeedAuth,
  onNeedKey,
}: {
  onNeedAuth: () => void;
  onNeedKey: () => void;
}) {
  const { user } = useAuth();
  const { key } = useActiveKey();
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<MediaResult | null>(null);

  async function go(e: React.FormEvent) {
    e.preventDefault();
    if (!url.trim()) return;
    if (!user) {
      toast.info("Masuk dulu untuk memakai downloader");
      onNeedAuth();
      return;
    }
    if (!key) {
      toast.info("Tambahkan API key dulu");
      onNeedKey();
      return;
    }
    setLoading(true);
    setResult(null);
    try {
      const r = await api.download(url.trim(), key);
      setResult(r);
      toast.success(`Media ${r.platform_name} ditemukan`);
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Gagal mengambil media";
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto w-full max-w-2xl">
      <form onSubmit={go} className="grad-border rounded-2xl p-2">
        <div className="flex flex-col gap-2 sm:flex-row">
          <div className="relative flex-1">
            <Link2
              size={18}
              className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-mist"
            />
            <input
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="Tempel link TikTok, YouTube, Instagram…"
              className={`${inputCls} border-0 bg-transparent pl-11 focus:ring-0`}
              spellCheck={false}
            />
          </div>
          <Button type="submit" loading={loading} className="sm:px-7">
            {!loading && <ArrowRight size={18} />} Ambil
          </Button>
        </div>
      </form>

      <div className="mt-3 flex flex-wrap items-center gap-2 text-xs text-mist">
        <span className="font-mono">coba:</span>
        {SAMPLES.map((s) => (
          <button
            key={s}
            onClick={() => setUrl(s)}
            className="max-w-[12rem] truncate rounded-md border border-line bg-surface/50 px-2 py-1 font-mono transition hover:border-violet/50 hover:text-paper"
          >
            {s.replace("https://www.", "").replace("https://", "")}
          </button>
        ))}
        {user && !key && (
          <button
            onClick={onNeedKey}
            className="ml-auto inline-flex items-center gap-1 rounded-md border border-violet/40 px-2 py-1 text-violet"
          >
            <KeyRound size={12} /> set API key
          </button>
        )}
      </div>

      <div className="mt-6">
        {loading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="flex items-center justify-center gap-3 rounded-2xl border border-line bg-surface/40 py-10 text-mist"
          >
            <Loader2 className="animate-spin text-violet" /> Menganalisis tautan…
          </motion.div>
        )}
        {result && !loading && <MediaResultCard result={result} />}
      </div>
    </div>
  );
}
