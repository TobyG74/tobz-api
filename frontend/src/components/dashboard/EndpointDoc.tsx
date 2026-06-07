import { useMemo, useState } from "react";
import { Loader2, Play } from "lucide-react";
import clsx from "clsx";
import { api, getAccessToken } from "../../lib/api";
import type { Endpoint } from "../../lib/endpoints";
import type { ReqDesc } from "../../lib/codegen";
import { useActiveKey } from "../../hooks/useActiveKey";
import { Button, inputCls } from "../primitives";
import { CodeTabs } from "./CodeTabs";
import { NoKeyNotice } from "./shared";

const methodColor: Record<string, string> = {
  GET: "from-blue to-azure",
  POST: "from-violet to-magenta",
  DELETE: "from-red-500 to-orange-500",
};

export function EndpointDoc({ endpoint, goKeys }: { endpoint: Endpoint; goKeys: () => void }) {
  const { key } = useActiveKey();
  // Editable param values, prefilled with examples.
  const [values, setValues] = useState<Record<string, string>>(() =>
    Object.fromEntries(endpoint.params.map((p) => [p.name, p.example]))
  );
  const [running, setRunning] = useState(false);
  const [resp, setResp] = useState<{ status: number; body: string } | null>(null);

  const authHeaders = useMemo<Record<string, string>>(() => {
    const h: Record<string, string> = {};
    if (endpoint.auth === "apikey") h["X-API-Key"] = key || "YOUR_API_KEY";
    else if (endpoint.auth === "bearer") h["Authorization"] = `Bearer ${getAccessToken() || "YOUR_ACCESS_TOKEN"}`;
    return h;
  }, [endpoint.auth, key]);

  const req: ReqDesc = useMemo(() => {
    const query: Record<string, string> = {};
    const body: Record<string, string> = {};
    for (const p of endpoint.params) {
      const v = values[p.name] ?? "";
      if (p.loc === "query") query[p.name] = v;
      else body[p.name] = v;
    }
    return {
      method: endpoint.method,
      url: api.base + endpoint.path,
      query,
      headers: authHeaders,
      body: Object.keys(body).length ? body : undefined,
    };
  }, [endpoint, values, authHeaders]);

  const needsKey = endpoint.auth === "apikey" && !key;

  async function run() {
    setRunning(true);
    setResp(null);
    try {
      const q = Object.entries(req.query)
        .filter(([, v]) => v !== "")
        .map(([k, v]) => `${k}=${encodeURIComponent(v)}`)
        .join("&");
      const headers: Record<string, string> = { ...authHeaders };
      if (req.body) headers["Content-Type"] = "application/json";
      const res = await fetch(req.url + (q ? `?${q}` : ""), {
        method: req.method,
        headers,
        credentials: "include",
        body: req.body ? JSON.stringify(req.body) : undefined,
      });
      const json = await res.json().catch(() => ({}));
      setResp({ status: res.status, body: JSON.stringify(json, null, 2) });
    } catch {
      setResp({ status: 0, body: "Gagal terhubung ke server" });
    } finally {
      setRunning(false);
    }
  }

  return (
    <div>
      {/* header */}
      <div className="mb-2 flex flex-wrap items-center gap-x-3 gap-y-2">
        <span
          className={clsx(
            "rounded-lg bg-gradient-to-r px-2.5 py-1 font-mono text-xs font-700 text-white",
            methodColor[endpoint.method]
          )}
        >
          {endpoint.method}
        </span>
        <code className="break-all font-mono text-sm text-paper">/api/v1{endpoint.path}</code>
        <span
          className={clsx(
            "rounded-md px-2 py-0.5 font-mono text-[11px] uppercase sm:ml-auto",
            endpoint.auth === "none" ? "bg-surface-2 text-mist" : "bg-violet/15 text-violet"
          )}
        >
          {endpoint.auth === "apikey" ? "API key" : endpoint.auth === "bearer" ? "Bearer" : "public"}
        </span>
      </div>
      <h2 className="font-display text-2xl">{endpoint.label}</h2>
      <p className="mt-1 text-fog">{endpoint.summary}</p>

      {/* params */}
      {endpoint.params.length > 0 && (
        <div className="mt-6">
          <h3 className="mb-3 font-mono text-xs uppercase tracking-[0.2em] text-azure">Parameter</h3>
          <div className="space-y-3">
            {endpoint.params.map((p) => (
              <div key={p.name} className="grid gap-2 sm:grid-cols-[180px_1fr] sm:items-center">
                <div>
                  <code className="font-mono text-sm text-paper">{p.name}</code>
                  {p.required && <span className="ml-1 text-red-400">*</span>}
                  <span className="ml-2 rounded bg-surface-2 px-1.5 py-0.5 font-mono text-[10px] text-mist">
                    {p.loc}
                  </span>
                  <p className="mt-0.5 text-xs text-mist">{p.desc}</p>
                </div>
                <input
                  value={values[p.name] ?? ""}
                  onChange={(e) => setValues((v) => ({ ...v, [p.name]: e.target.value }))}
                  placeholder={p.example}
                  className={`${inputCls} py-2 font-mono text-sm`}
                />
              </div>
            ))}
          </div>
        </div>
      )}

      {/* code */}
      <div className="mt-6">
        <h3 className="mb-3 font-mono text-xs uppercase tracking-[0.2em] text-azure">Contoh kode</h3>
        <CodeTabs req={req} />
      </div>

      {/* try it */}
      <div className="mt-6">
        <div className="mb-3 flex items-center justify-between">
          <h3 className="font-mono text-xs uppercase tracking-[0.2em] text-azure">Coba langsung</h3>
          <Button onClick={run} loading={running} disabled={needsKey}>
            {!running && <Play size={15} />} Jalankan
          </Button>
        </div>

        {needsKey ? (
          <NoKeyNotice goKeys={goKeys} />
        ) : resp ? (
          <div className="overflow-hidden rounded-2xl border border-line">
            <div className="flex items-center gap-2 border-b border-line bg-void/60 px-4 py-2 font-mono text-xs">
              <span
                className={clsx(
                  "h-2 w-2 rounded-full",
                  resp.status >= 200 && resp.status < 300 ? "bg-[#28c840]" : "bg-red-500"
                )}
              />
              <span className="text-mist">status {resp.status || "—"}</span>
            </div>
            <pre className="max-h-96 overflow-auto bg-void/40 p-4 font-mono text-[13px] leading-relaxed text-paper">
              <code>{resp.body}</code>
            </pre>
          </div>
        ) : (
          <div className="rounded-2xl border border-dashed border-line p-6 text-center text-sm text-mist">
            {running ? (
              <span className="inline-flex items-center gap-2">
                <Loader2 className="animate-spin" /> Mengirim…
              </span>
            ) : (
              "Klik Jalankan untuk mengirim request nyata & melihat respons."
            )}
          </div>
        )}
      </div>

      {/* example response */}
      <div className="mt-6">
        <h3 className="mb-3 font-mono text-xs uppercase tracking-[0.2em] text-azure">Contoh respons</h3>
        <pre className="overflow-x-auto rounded-2xl border border-line bg-void/40 p-4 font-mono text-[13px] leading-relaxed text-fog">
          <code>{endpoint.response}</code>
        </pre>
      </div>
    </div>
  );
}
