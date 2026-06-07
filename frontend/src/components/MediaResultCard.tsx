import { motion } from "motion/react";
import { Download, Film, ImageIcon, Music, User } from "lucide-react";
import type { DownloadItem, MediaResult } from "../lib/types";

function typeIcon(t: string) {
  if (t === "audio") return <Music size={15} />;
  if (t === "image") return <ImageIcon size={15} />;
  return <Film size={15} />;
}

function Item({ item }: { item: DownloadItem }) {
  return (
    <a
      href={item.url}
      target="_blank"
      rel="noreferrer"
      className="group flex items-center justify-between gap-3 rounded-xl border border-line bg-void/50 px-4 py-3 transition hover:border-violet/60 hover:bg-surface-2/60"
    >
      <span className="flex items-center gap-2.5 text-sm">
        <span className="grid h-7 w-7 place-items-center rounded-lg bg-surface-2 text-azure">
          {typeIcon(item.type)}
        </span>
        <span className="font-500 text-paper">{item.label || item.type}</span>
        {item.quality && (
          <span className="rounded-md bg-violet/15 px-1.5 py-0.5 font-mono text-[11px] text-violet">
            {item.quality}
          </span>
        )}
      </span>
      <span className="flex items-center gap-1.5 font-mono text-xs text-mist transition group-hover:text-cyan">
        <Download size={14} /> unduh
      </span>
    </a>
  );
}

export function MediaResultCard({ result }: { result: MediaResult }) {
  const hasImages = result.images && result.images.length > 0;
  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
      className="grad-border overflow-hidden rounded-2xl"
    >
      <div className="grid gap-5 p-5 sm:grid-cols-[auto_1fr]">
        {result.thumbnail ? (
          <img
            src={result.thumbnail}
            alt=""
            className="h-28 w-28 shrink-0 rounded-xl object-cover ring-1 ring-line"
            referrerPolicy="no-referrer"
            onError={(e) => (e.currentTarget.style.display = "none")}
          />
        ) : null}

        <div className="min-w-0">
          <div className="mb-1 flex items-center gap-2">
            <span className="rounded-md bg-gradient-to-r from-blue to-violet px-2 py-0.5 font-mono text-[11px] font-600 uppercase tracking-wider text-white">
              {result.platform_name}
            </span>
            <span className="font-mono text-[11px] text-mist">via {result.downloader}</span>
          </div>
          <h4 className="truncate font-display text-lg text-paper" title={result.title}>
            {result.title || "Media siap diunduh"}
          </h4>
          {result.author_name && (
            <p className="mt-0.5 flex items-center gap-1.5 text-sm text-mist">
              <User size={13} /> {result.author_name}
            </p>
          )}
        </div>
      </div>

      {result.download_items.length > 0 && (
        <div className="space-y-2 px-5 pb-5">
          {result.download_items.map((it) => (
            <Item key={it.key} item={it} />
          ))}
        </div>
      )}

      {hasImages && (
        <div className="grid grid-cols-3 gap-2 px-5 pb-5 sm:grid-cols-4">
          {result.images!.map((src, i) => (
            <a key={i} href={src} target="_blank" rel="noreferrer" className="group relative">
              <img
                src={src}
                alt=""
                referrerPolicy="no-referrer"
                className="aspect-square w-full rounded-lg object-cover ring-1 ring-line transition group-hover:ring-violet/70"
              />
              <span className="absolute inset-0 grid place-items-center rounded-lg bg-void/0 text-cyan opacity-0 transition group-hover:bg-void/50 group-hover:opacity-100">
                <Download size={18} />
              </span>
            </a>
          ))}
        </div>
      )}
    </motion.div>
  );
}
