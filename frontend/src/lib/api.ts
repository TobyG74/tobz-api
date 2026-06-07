import type {
  ApiKey,
  AuthResult,
  CreatedApiKey,
  Envelope,
  ImageResult,
  MediaResult,
  PixivResult,
  PlatformInfo,
  User,
  WhitelistIP,
} from "./types";

const BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080/api/v1";

/** Thrown for any non-success API response; carries a user-safe message. */
export class ApiError extends Error {
  code: string;
  status: number;
  constructor(message: string, code: string, status: number) {
    super(message);
    this.code = code;
    this.status = status;
  }
}

// Access token is held in memory only (never localStorage) to limit XSS blast radius.
let accessToken: string | null = null;
export function setAccessToken(t: string | null) {
  accessToken = t;
}
export function getAccessToken() {
  return accessToken;
}

interface RequestOpts {
  method?: string;
  body?: unknown;
  auth?: boolean; // attach Bearer token
  apiKey?: string; // attach X-API-Key
  headers?: Record<string, string>;
}

async function request<T>(path: string, opts: RequestOpts = {}): Promise<T> {
  const headers: Record<string, string> = { ...(opts.headers ?? {}) };
  if (opts.body !== undefined) headers["Content-Type"] = "application/json";
  if (opts.auth && accessToken) headers["Authorization"] = `Bearer ${accessToken}`;
  if (opts.apiKey) headers["X-API-Key"] = opts.apiKey;

  let res: Response;
  try {
    res = await fetch(`${BASE}${path}`, {
      method: opts.method ?? "GET",
      headers,
      credentials: "include", // send/receive the httpOnly refresh cookie
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    });
  } catch {
    throw new ApiError("Tidak dapat terhubung ke server API", "network", 0);
  }

  let json: Envelope<T> | null = null;
  try {
    json = (await res.json()) as Envelope<T>;
  } catch {
    /* non-JSON response */
  }

  if (!res.ok || !json || json.success === false) {
    const msg = json?.error?.message ?? `Request gagal (${res.status})`;
    const code = json?.error?.code ?? "error";
    throw new ApiError(msg, code, res.status);
  }
  return json.data as T;
}

export const api = {
  base: BASE,

  // ---- Auth ----
  register: (body: { email: string; password: string; display_name?: string; captcha_token: string }) =>
    request<AuthResult>("/auth/register", { method: "POST", body }),

  login: (body: { email: string; password: string; captcha_token: string }) =>
    request<AuthResult>("/auth/login", { method: "POST", body }),

  refresh: () => request<AuthResult>("/auth/refresh", { method: "POST" }),

  logout: () => request<{ message: string }>("/auth/logout", { method: "POST" }),

  me: () => request<User>("/auth/me", { auth: true }),

  changePassword: (body: { current_password: string; new_password: string }) =>
    request<AuthResult>("/auth/change-password", { method: "POST", body, auth: true }),

  oauthURL: (provider: "google" | "github") => `${BASE}/auth/oauth/${provider}/login`,

  // ---- API keys ----
  listKeys: () => request<ApiKey[]>("/keys/", { auth: true }),

  createKey: (body: { name: string; tier?: string }) =>
    request<CreatedApiKey>("/keys/", { method: "POST", body, auth: true }),

  revokeKey: (id: string) => request<{ message: string }>(`/keys/${id}`, { method: "DELETE", auth: true }),

  // ---- IP whitelist ----
  listWhitelist: () => request<{ ips: WhitelistIP[]; max: number }>("/whitelist/", { auth: true }),

  addWhitelist: (body: { ip: string; label?: string }) =>
    request<WhitelistIP>("/whitelist/", { method: "POST", body, auth: true }),

  removeWhitelist: (id: string) =>
    request<{ message: string }>(`/whitelist/${id}`, { method: "DELETE", auth: true }),

  // ---- Downloader ----
  platforms: (apiKey: string) =>
    request<{ platforms: PlatformInfo[] }>("/download/platforms", { apiKey }),

  download: (url: string, apiKey: string) =>
    request<MediaResult>(`/download?url=${encodeURIComponent(url)}`, { apiKey }),

  // ---- Search ----
  searchSources: (apiKey: string) =>
    request<{ sources: string[] }>("/search/sources", { apiKey }),

  searchImages: (q: string, source: string, limit: number, apiKey: string) =>
    request<{ source: string; query: string; count: number; results: ImageResult[] }>(
      `/search/images?q=${encodeURIComponent(q)}&source=${source}&limit=${limit}`,
      { apiKey }
    ),

  searchPixiv: (q: string, type: string, apiKey: string) =>
    request<{ query: string; type: string; count: number; results: PixivResult[] }>(
      `/search/pixiv?q=${encodeURIComponent(q)}&type=${type}`,
      { apiKey }
    ),
};
