import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { Globe, Info, Plus, Shield, Trash2 } from "lucide-react";
import { api, ApiError } from "../lib/api";
import type { WhitelistIP } from "../lib/types";
import { Button, Field, inputCls } from "./primitives";

function fmtDate(s: string) {
  try {
    return new Date(s).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
  } catch {
    return s;
  }
}

export function WhitelistManager() {
  const [ips, setIps] = useState<WhitelistIP[]>([]);
  const [max, setMax] = useState(5);
  const [loading, setLoading] = useState(false);
  const [ip, setIp] = useState("");
  const [label, setLabel] = useState("");
  const [adding, setAdding] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const r = await api.listWhitelist();
      setIps(r.ips);
      setMax(r.max);
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Gagal memuat whitelist");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  async function add(e: React.FormEvent) {
    e.preventDefault();
    setAdding(true);
    try {
      await api.addWhitelist({ ip: ip.trim(), label: label.trim() });
      setIp("");
      setLabel("");
      toast.success("IP ditambahkan");
      load();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Gagal menambahkan IP");
    } finally {
      setAdding(false);
    }
  }

  async function remove(id: string) {
    try {
      await api.removeWhitelist(id);
      toast.success("IP dihapus");
      load();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Gagal menghapus IP");
    }
  }

  const full = ips.length >= max;
  const isPublic = ips.some((x) => x.ip === "0.0.0.0");

  return (
    <div>
      {/* explainer */}
      <div className="mb-5 rounded-2xl border border-line bg-surface/40 p-5">
        <p className="mb-3 flex items-center gap-2 font-mono text-xs uppercase tracking-[0.2em] text-azure">
          <Info size={14} /> Cara kerja whitelist
        </p>
        <ul className="space-y-2 text-sm text-fog">
          <li className="flex gap-2"><Shield size={15} className="mt-0.5 shrink-0 text-cyan" /> Membatasi <b>login</b> & <b>pemakaian API key</b> hanya dari IP terdaftar.</li>
          <li className="flex gap-2"><Globe size={15} className="mt-0.5 shrink-0 text-cyan" /> <code className="font-mono">0.0.0.0</code> = <b>publik</b> (boleh dari IP mana pun).</li>
          <li className="flex gap-2"><Info size={15} className="mt-0.5 shrink-0 text-cyan" /> Daftar kosong = tanpa batasan (semua IP diizinkan).</li>
        </ul>
      </div>

      {/* status */}
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-line bg-void/50 p-4">
        <div>
          <p className="text-xs uppercase tracking-widest text-mist">IP tersimpan</p>
          <p className="font-display text-2xl">
            {ips.length}<span className="text-mist">/{max}</span>
          </p>
        </div>
        <span
          className={`rounded-lg px-3 py-1.5 font-mono text-xs ${
            isPublic ? "bg-cyan/15 text-cyan" : ips.length ? "bg-violet/15 text-violet" : "bg-surface-2 text-mist"
          }`}
        >
          {isPublic ? "MODE PUBLIK" : ips.length ? "TERBATAS" : "TANPA BATASAN"}
        </span>
      </div>

      {/* add form */}
      <form onSubmit={add} className="mb-6 flex flex-col gap-2 sm:flex-row sm:items-end">
        <div className="flex-1">
          <Field label="Alamat IP" hint="mis. 103.10.20.30 atau 0.0.0.0">
            <input
              className={`${inputCls} font-mono`}
              value={ip}
              onChange={(e) => setIp(e.target.value)}
              placeholder="0.0.0.0"
              disabled={full}
            />
          </Field>
        </div>
        <div className="flex-1">
          <Field label="Label (opsional)">
            <input
              className={inputCls}
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              placeholder="mis. Kantor"
              disabled={full}
            />
          </Field>
        </div>
        <Button type="submit" loading={adding} disabled={full}>
          <Plus size={16} /> Tambah
        </Button>
      </form>
      {full && (
        <p className="-mt-4 mb-6 text-xs text-mist">Batas maksimal {max} IP tercapai — hapus salah satu untuk menambah.</p>
      )}

      {/* list */}
      <p className="mb-2 text-xs uppercase tracking-widest text-mist">Daftar IP</p>
      <div className="space-y-2">
        {loading && <p className="text-sm text-mist">Memuat…</p>}
        {!loading && ips.length === 0 && (
          <div className="rounded-xl border border-dashed border-line p-6 text-center text-sm text-mist">
            <Shield className="mx-auto mb-2 opacity-50" /> Belum ada IP — akun bisa diakses dari mana saja.
          </div>
        )}
        {ips.map((w) => (
          <div key={w.id} className="flex items-center justify-between gap-3 rounded-xl border border-line bg-surface/50 p-4">
            <div className="flex items-center gap-3">
              <span className={`grid h-8 w-8 place-items-center rounded-lg ${w.ip === "0.0.0.0" ? "bg-cyan/15 text-cyan" : "bg-surface-2 text-azure"}`}>
                {w.ip === "0.0.0.0" ? <Globe size={16} /> : <Shield size={16} />}
              </span>
              <div>
                <code className="font-mono text-sm text-paper">{w.ip}</code>
                {w.ip === "0.0.0.0" && <span className="ml-2 rounded bg-cyan/15 px-1.5 py-0.5 text-[10px] uppercase text-cyan">publik</span>}
                <p className="text-xs text-mist">
                  {w.label ? `${w.label} · ` : ""}ditambahkan {fmtDate(w.created_at)}
                </p>
              </div>
            </div>
            <button
              onClick={() => remove(w.id)}
              className="rounded-lg p-1.5 text-mist hover:bg-red-500/10 hover:text-red-300"
              aria-label="Hapus"
            >
              <Trash2 size={15} />
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
