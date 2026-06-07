import { KeyRound, Sparkles } from "lucide-react";
import type { ReactNode } from "react";
import { Button } from "../primitives";

export function PanelHeader({ icon, title, desc }: { icon: ReactNode; title: string; desc: string }) {
  return (
    <div className="mb-6 flex items-start gap-3">
      <div className="grid h-11 w-11 place-items-center rounded-xl bg-gradient-to-br from-blue to-violet text-white">
        {icon}
      </div>
      <div>
        <h2 className="font-display text-2xl">{title}</h2>
        <p className="text-sm text-fog">{desc}</p>
      </div>
    </div>
  );
}

/** Collapsible-free "how it works" explainer used on each feature panel. */
export function HowBox({ steps }: { steps: string[] }) {
  return (
    <div className="mb-6 rounded-2xl border border-line bg-surface/40 p-5">
      <p className="mb-3 flex items-center gap-2 font-mono text-xs uppercase tracking-[0.2em] text-azure">
        <Sparkles size={14} /> Cara kerja
      </p>
      <ol className="space-y-2.5">
        {steps.map((s, i) => (
          <li key={i} className="flex gap-3 text-sm text-fog">
            <span className="grid h-5 w-5 shrink-0 place-items-center rounded-full bg-violet/15 font-mono text-[11px] text-violet">
              {i + 1}
            </span>
            {s}
          </li>
        ))}
      </ol>
    </div>
  );
}

/** Shown when no active API key is set; routes to the keys section. */
export function NoKeyNotice({ goKeys }: { goKeys: () => void }) {
  return (
    <div className="rounded-2xl border border-dashed border-violet/40 bg-violet/5 p-8 text-center">
      <div className="mx-auto mb-3 grid h-12 w-12 place-items-center rounded-xl bg-violet/15 text-violet">
        <KeyRound size={22} />
      </div>
      <h3 className="font-display text-lg">Butuh API key dulu</h3>
      <p className="mx-auto mt-1 max-w-sm text-sm text-fog">
        Fitur ini memakai API key. Buat satu di tab API Keys — key baru otomatis diaktifkan.
      </p>
      <Button onClick={goKeys} className="mt-4">
        <KeyRound size={16} /> Ke API Keys
      </Button>
    </div>
  );
}
