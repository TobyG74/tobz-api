import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { Copy, KeyRound, Plus, Star, Trash2 } from "lucide-react";
import { api, ApiError } from "../lib/api";
import type { ApiKey, CreatedApiKey } from "../lib/types";
import { useActiveKey } from "../hooks/useActiveKey";
import { Button, Field, inputCls, Sheet } from "./primitives";

export function KeysPanel({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { key: activeKey, setKey: setActiveKey } = useActiveKey();
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [loading, setLoading] = useState(false);
  const [name, setName] = useState("");
  const [creating, setCreating] = useState(false);
  const [justCreated, setJustCreated] = useState<CreatedApiKey | null>(null);
  const [manual, setManual] = useState("");

  const load = useCallback(async () => {
    setLoading(true);
    try {
      setKeys(await api.listKeys());
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Gagal memuat keys");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (open) load();
  }, [open, load]);

  async function create(e: React.FormEvent) {
    e.preventDefault();
    setCreating(true);
    try {
      const created = await api.createKey({ name: name || "Default" });
      setJustCreated(created);
      setActiveKey(created.api_key); // auto-activate for the downloader
      setName("");
      toast.success("API key dibuat & diaktifkan");
      load();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Gagal membuat key");
    } finally {
      setCreating(false);
    }
  }

  async function revoke(id: string) {
    try {
      await api.revokeKey(id);
      toast.success("Key dicabut");
      load();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Gagal mencabut key");
    }
  }

  function copy(text: string) {
    navigator.clipboard.writeText(text);
    toast.success("Tersalin");
  }

  return (
    <Sheet open={open} onClose={onClose} title="API Keys">
      {/* Active key status */}
      <div className="mb-5 rounded-xl border border-line bg-void/50 p-4">
        <p className="mb-1 text-xs uppercase tracking-widest text-mist">Key aktif (dipakai downloader)</p>
        {activeKey ? (
          <div className="flex items-center justify-between gap-2">
            <code className="truncate font-mono text-sm text-cyan">
              {activeKey.slice(0, 12)}••••••••
            </code>
            <button onClick={() => setActiveKey("")} className="text-xs text-mist hover:text-red-300">
              hapus
            </button>
          </div>
        ) : (
          <p className="text-sm text-mist">Belum ada — buat baru atau tempel di bawah.</p>
        )}
        <div className="mt-3 flex gap-2">
          <input
            value={manual}
            onChange={(e) => setManual(e.target.value)}
            placeholder="tempel API key (tobz_…)"
            className={`${inputCls} py-2 font-mono text-xs`}
          />
          <Button
            variant="outline"
            type="button"
            onClick={() => {
              if (manual.trim()) {
                setActiveKey(manual.trim());
                setManual("");
                toast.success("Key aktif diperbarui");
              }
            }}
          >
            Pakai
          </Button>
        </div>
      </div>

      {/* Freshly created key (shown once) */}
      {justCreated && (
        <div className="mb-5 rounded-xl border border-violet/50 bg-violet/10 p-4">
          <p className="mb-2 text-sm font-600 text-violet">Simpan key ini — hanya tampil sekali!</p>
          <div className="flex items-center gap-2">
            <code className="flex-1 truncate rounded-lg bg-void/60 px-3 py-2 font-mono text-sm text-paper">
              {justCreated.api_key}
            </code>
            <Button variant="outline" onClick={() => copy(justCreated.api_key)}>
              <Copy size={15} />
            </Button>
          </div>
        </div>
      )}

      {/* Create */}
      <form onSubmit={create} className="mb-6 flex items-end gap-2">
        <div className="flex-1">
          <Field label="Buat key baru">
            <input
              className={inputCls}
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Nama (mis. produksi)"
            />
          </Field>
        </div>
        <Button type="submit" loading={creating}>
          <Plus size={16} /> Buat
        </Button>
      </form>

      {/* List */}
      <p className="mb-2 text-xs uppercase tracking-widest text-mist">Keys kamu</p>
      <div className="space-y-2">
        {loading && <p className="text-sm text-mist">Memuat…</p>}
        {!loading && keys.length === 0 && (
          <div className="rounded-xl border border-dashed border-line p-6 text-center text-sm text-mist">
            <KeyRound className="mx-auto mb-2 opacity-50" /> Belum ada API key.
          </div>
        )}
        {keys.map((k) => {
          const pct = Math.min(100, Math.round((k.quota_used / k.daily_quota) * 100));
          return (
            <div key={k.id} className="rounded-xl border border-line bg-surface/50 p-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="font-500 text-paper">{k.name || "(tanpa nama)"}</span>
                  <span className="rounded-md bg-surface-2 px-1.5 py-0.5 font-mono text-[11px] uppercase text-azure">
                    {k.tier}
                  </span>
                  {k.revoked && (
                    <span className="rounded-md bg-red-500/15 px-1.5 py-0.5 text-[11px] text-red-300">
                      dicabut
                    </span>
                  )}
                </div>
                {!k.revoked && (
                  <button
                    onClick={() => revoke(k.id)}
                    className="rounded-lg p-1.5 text-mist hover:bg-red-500/10 hover:text-red-300"
                    aria-label="Cabut"
                  >
                    <Trash2 size={15} />
                  </button>
                )}
              </div>
              <code className="mt-1 block font-mono text-xs text-mist">{k.prefix}••••••••</code>

              <div className="mt-3">
                <div className="mb-1 flex justify-between font-mono text-[11px] text-mist">
                  <span>kuota harian</span>
                  <span>
                    {k.quota_used}/{k.daily_quota}
                  </span>
                </div>
                <div className="h-1.5 overflow-hidden rounded-full bg-void">
                  <div
                    className="h-full rounded-full bg-gradient-to-r from-blue to-violet"
                    style={{ width: `${pct}%` }}
                  />
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="mt-6 flex items-center gap-2 rounded-xl border border-line bg-void/40 p-3 text-xs text-mist">
        <Star size={14} className="shrink-0 text-cyan" />
        Key yang baru dibuat otomatis diaktifkan untuk downloader.
      </div>
    </Sheet>
  );
}
