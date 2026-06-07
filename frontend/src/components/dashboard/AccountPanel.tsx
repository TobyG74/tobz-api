import { useState } from "react";
import { toast } from "sonner";
import { Lock, Mail, ShieldCheck, User } from "lucide-react";
import { api, ApiError } from "../../lib/api";
import { useAuth } from "../../context/AuthContext";
import { Button, Field, inputCls } from "../primitives";

export function AccountPanel() {
  const { user, setSession } = useAuth();
  const [current, setCurrent] = useState("");
  const [next, setNext] = useState("");
  const [confirm, setConfirm] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (next !== confirm) {
      toast.error("Konfirmasi password tidak cocok");
      return;
    }
    setLoading(true);
    try {
      const res = await api.changePassword({ current_password: current, new_password: next });
      setSession(res); // refresh the session (others were revoked)
      setCurrent("");
      setNext("");
      setConfirm("");
      toast.success("Password berhasil diganti. Sesi lain telah keluar.");
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Gagal mengganti password");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="max-w-xl">
      {/* profile */}
      <div className="mb-6 rounded-2xl border border-line bg-surface/40 p-5">
        <p className="mb-4 font-mono text-xs uppercase tracking-[0.2em] text-azure">Profil</p>
        <div className="flex items-center gap-4">
          <span className="grid h-14 w-14 place-items-center rounded-2xl bg-gradient-to-br from-blue to-violet text-xl font-700 text-white">
            {(user?.display_name || "U").slice(0, 1).toUpperCase()}
          </span>
          <div className="min-w-0">
            <p className="flex items-center gap-1.5 font-display text-lg text-paper">
              <User size={15} className="text-mist" /> {user?.display_name}
            </p>
            <p className="flex items-center gap-1.5 truncate text-sm text-mist">
              <Mail size={13} /> {user?.email}
              {user?.email_verified && <ShieldCheck size={13} className="text-cyan" />}
            </p>
          </div>
        </div>
      </div>

      {/* change password */}
      <div className="rounded-2xl border border-line bg-surface/40 p-5">
        <p className="mb-1 flex items-center gap-2 font-display text-lg">
          <Lock size={17} className="text-violet" /> Ganti Password
        </p>
        <p className="mb-5 text-sm text-mist">
          Setelah diganti, semua sesi lain otomatis keluar demi keamanan.
        </p>

        <form onSubmit={submit} className="space-y-4">
          <Field label="Password saat ini" hint="kosongkan jika akun sosial & belum punya">
            <input
              className={inputCls}
              type="password"
              value={current}
              onChange={(e) => setCurrent(e.target.value)}
              placeholder="••••••••"
              autoComplete="current-password"
            />
          </Field>
          <Field label="Password baru" hint="min 8 char, huruf + angka">
            <input
              className={inputCls}
              type="password"
              required
              value={next}
              onChange={(e) => setNext(e.target.value)}
              placeholder="••••••••"
              autoComplete="new-password"
            />
          </Field>
          <Field label="Konfirmasi password baru">
            <input
              className={inputCls}
              type="password"
              required
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              placeholder="••••••••"
              autoComplete="new-password"
            />
          </Field>
          <Button type="submit" loading={loading} className="w-full">
            Simpan password baru
          </Button>
        </form>
      </div>
    </div>
  );
}
