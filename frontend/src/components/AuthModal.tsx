import { useState } from "react";
import { toast } from "sonner";
import { Github, Mail } from "lucide-react";
import { api, ApiError } from "../lib/api";
import { useAuth } from "../context/AuthContext";
import { Button, Field, inputCls, Modal } from "./primitives";
import { Turnstile } from "./Turnstile";

type Mode = "login" | "register";

export function AuthModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { setSession } = useAuth();
  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [captcha, setCaptcha] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!captcha) {
      toast.error("Selesaikan captcha terlebih dahulu");
      return;
    }
    setLoading(true);
    try {
      const res =
        mode === "login"
          ? await api.login({ email, password, captcha_token: captcha })
          : await api.register({ email, password, display_name: name, captcha_token: captcha });
      setSession(res);
      toast.success(mode === "login" ? "Selamat datang kembali 👋" : "Akun berhasil dibuat 🎉");
      onClose();
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Terjadi kesalahan";
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Modal open={open} onClose={onClose} title={mode === "login" ? "Masuk" : "Buat akun"}>
      <div className="mb-5 grid grid-cols-2 gap-3">
        <a href={api.oauthURL("google")} className="contents">
          <Button variant="outline" type="button" className="w-full">
            <Mail size={16} /> Google
          </Button>
        </a>
        <a href={api.oauthURL("github")} className="contents">
          <Button variant="outline" type="button" className="w-full">
            <Github size={16} /> GitHub
          </Button>
        </a>
      </div>

      <div className="mb-5 flex items-center gap-3 text-xs uppercase tracking-widest text-mist">
        <span className="h-px flex-1 bg-line" /> atau email <span className="h-px flex-1 bg-line" />
      </div>

      <form onSubmit={submit} className="space-y-4">
        {mode === "register" && (
          <Field label="Nama">
            <input
              className={inputCls}
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Nama tampilan"
              autoComplete="name"
            />
          </Field>
        )}
        <Field label="Email">
          <input
            className={inputCls}
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="kamu@email.com"
            autoComplete="email"
          />
        </Field>
        <Field label="Password" hint={mode === "register" ? "min 8 char, huruf+angka" : undefined}>
          <input
            className={inputCls}
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            autoComplete={mode === "login" ? "current-password" : "new-password"}
          />
        </Field>

        <Turnstile onToken={setCaptcha} />

        <Button type="submit" loading={loading} className="w-full">
          {mode === "login" ? "Masuk" : "Daftar sekarang"}
        </Button>
      </form>

      <p className="mt-5 text-center text-sm text-mist">
        {mode === "login" ? "Belum punya akun?" : "Sudah punya akun?"}{" "}
        <button
          onClick={() => setMode(mode === "login" ? "register" : "login")}
          className="font-600 text-azure hover:text-cyan"
        >
          {mode === "login" ? "Daftar" : "Masuk"}
        </button>
      </p>
    </Modal>
  );
}
