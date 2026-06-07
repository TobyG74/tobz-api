import { useMemo, useState } from "react";
import { Check, Copy } from "lucide-react";
import { toast } from "sonner";
import { generate, LANGS, type ReqDesc } from "../../lib/codegen";
import type { Lang } from "../../lib/endpoints";

export function CodeTabs({ req }: { req: ReqDesc }) {
  const [lang, setLang] = useState<Lang>("curl");
  const [copied, setCopied] = useState(false);

  const code = useMemo(() => generate(lang, req), [lang, req]);

  function copy() {
    navigator.clipboard.writeText(code);
    setCopied(true);
    toast.success("Kode tersalin");
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-line">
      <div className="flex items-center justify-between border-b border-line bg-void/60 px-2">
        <div className="flex overflow-x-auto">
          {LANGS.map((l) => (
            <button
              key={l.id}
              onClick={() => setLang(l.id)}
              className={`whitespace-nowrap px-3 py-2.5 text-xs font-600 font-mono transition ${
                lang === l.id ? "text-cyan" : "text-mist hover:text-paper"
              }`}
            >
              {l.label}
              {lang === l.id && <span className="mt-1 block h-0.5 rounded-full bg-gradient-to-r from-blue to-violet" />}
            </button>
          ))}
        </div>
        <button
          onClick={copy}
          className="mr-1 flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-xs text-fog transition hover:bg-surface-2 hover:text-paper"
        >
          {copied ? <Check size={14} className="text-cyan" /> : <Copy size={14} />}
          {copied ? "Tersalin" : "Copy"}
        </button>
      </div>
      <pre className="overflow-x-auto bg-void/40 p-4 font-mono text-[13px] leading-relaxed text-paper">
        <code>{code}</code>
      </pre>
    </div>
  );
}
