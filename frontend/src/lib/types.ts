// API contract types mirroring the Go backend's JSON responses.

export interface Envelope<T> {
  success: boolean;
  data?: T;
  error?: { code: string; message: string };
}

export interface User {
  id: string;
  email: string;
  email_verified: boolean;
  display_name: string;
  avatar_url?: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResult {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface ApiKey {
  id: string;
  name: string;
  prefix: string;
  tier: string;
  daily_quota: number;
  quota_used: number;
  quota_reset_at: string;
  revoked: boolean;
  expires_at: string | null;
  last_used_at: string | null;
  created_at: string;
}

export interface CreatedApiKey {
  id: string;
  name: string;
  tier: string;
  daily_quota: number;
  prefix: string;
  api_key: string;
  warning: string;
}

export interface DownloadItem {
  key: string;
  label: string;
  type: "video" | "audio" | "image" | "playlist_item" | string;
  url: string;
  mime_type?: string;
  quality?: string;
}

export interface MediaResult {
  platform: string;
  platform_name: string;
  downloader: string;
  title?: string;
  author_name?: string;
  duration?: number;
  thumbnail?: string;
  download_items: DownloadItem[];
  images?: string[];
}

export interface PlatformInfo {
  id: string;
  platform: string;
  platform_name: string;
  downloader: string;
}

export interface ImageResult {
  url: string;
  thumbnail: string;
  title?: string;
  width?: number;
  height?: number;
  source?: string;
  author?: string;
  type?: string;
}

export interface PixivResult {
  id: string;
  title: string;
  user_id: string;
  user_name: string;
  type: string;
}

export interface WhitelistIP {
  id: string;
  user_id: string;
  ip: string;
  label: string;
  created_at: string;
}
