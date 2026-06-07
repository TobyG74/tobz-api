// API endpoint catalog that powers the Swagger-style docs in the dashboard.

export type Lang = "curl" | "python" | "node" | "go" | "ruby" | "php";
export type GroupId = "downloader" | "search";
export type AuthKind = "apikey" | "bearer" | "none";

export interface Param {
  name: string;
  loc: "query" | "body";
  required: boolean;
  example: string;
  desc: string;
}

export interface Endpoint {
  id: string;
  group: GroupId;
  label: string;
  method: "GET" | "POST";
  path: string;
  summary: string;
  auth: AuthKind;
  params: Param[];
  response: string; // pretty JSON example
}

export interface Group {
  id: GroupId;
  label: string;
  endpoints: Endpoint[];
}

const mediaResponse = `{
  "success": true,
  "data": {
    "platform": "tiktok",
    "platform_name": "TikTok",
    "downloader": "MusicalDown",
    "title": "...",
    "author_name": "...",
    "thumbnail": "https://...",
    "download_items": [
      { "key": "video_hd", "label": "Video HD", "type": "video",
        "url": "https://.../hd.mp4", "mime_type": "video/mp4", "quality": "HD" }
    ],
    "images": []
  }
}`;

const imageResponse = `{
  "success": true,
  "data": {
    "source": "bing",
    "query": "sunset",
    "count": 30,
    "results": [
      { "url": "https://.../full.jpg", "thumbnail": "https://.../thumb.jpg",
        "title": "...", "width": 1920, "height": 1080, "source": "bing", "type": "image" }
    ]
  }
}`;

const pixivResponse = `{
  "success": true,
  "data": {
    "query": "miku", "type": "artworks", "count": 60,
    "results": [
      { "id": "12345678", "title": "...", "user_id": "111",
        "user_name": "artist", "type": "artworks" }
    ]
  }
}`;

function dl(id: string, label: string, platform: string, exampleURL: string): Endpoint {
  return {
    id,
    group: "downloader",
    label,
    method: "GET",
    path: "/download",
    summary: `Ambil tautan unduhan media dari ${platform}. Platform dideteksi otomatis dari URL.`,
    auth: "apikey",
    params: [{ name: "url", loc: "query", required: true, example: exampleURL, desc: `URL ${platform}` }],
    response: mediaResponse,
  };
}

function img(id: string, label: string, source: string): Endpoint {
  return {
    id,
    group: "search",
    label,
    method: "GET",
    path: "/search/images",
    summary: `Cari gambar via ${label}.`,
    auth: "apikey",
    params: [
      { name: "q", loc: "query", required: true, example: "sunset", desc: "Kata kunci pencarian" },
      { name: "source", loc: "query", required: false, example: source, desc: "Sumber gambar" },
      { name: "limit", loc: "query", required: false, example: "30", desc: "Jumlah hasil (maks 150)" },
    ],
    response: imageResponse,
  };
}

export const GROUPS: Group[] = [
  {
    id: "downloader",
    label: "Downloader",
    endpoints: [
      dl("dl_tiktok", "TikTok Downloader", "TikTok", "https://www.tiktok.com/@user/video/123"),
      dl("dl_youtube", "YouTube Downloader", "YouTube", "https://youtu.be/dQw4w9WgXcQ"),
      dl("dl_instagram", "Instagram Downloader", "Instagram", "https://www.instagram.com/p/Cabc/"),
      dl("dl_facebook", "Facebook Downloader", "Facebook", "https://www.facebook.com/reel/123"),
      dl("dl_twitter", "Twitter / X Downloader", "Twitter/X", "https://x.com/user/status/123"),
      dl("dl_douyin", "Douyin Downloader", "Douyin", "https://v.douyin.com/abcdef/"),
      {
        id: "dl_platforms",
        group: "downloader",
        label: "List Platforms",
        method: "GET",
        path: "/download/platforms",
        summary: "Daftar semua platform downloader yang didukung.",
        auth: "apikey",
        params: [],
        response: `{ "success": true, "data": { "platforms": [
  { "id": "musicaldown", "platform": "tiktok", "platform_name": "TikTok", "downloader": "musicaldown" }
] } }`,
      },
    ],
  },
  {
    id: "search",
    label: "Search",
    endpoints: [
      img("se_bing", "Bing Images", "bing"),
      img("se_ddg", "DuckDuckGo Images", "duckduckgo"),
      img("se_pexels", "Pexels", "pexels"),
      img("se_pinterest", "Pinterest", "pinterest"),
      {
        id: "se_pixiv",
        group: "search",
        label: "Pixiv Search",
        method: "GET",
        path: "/search/pixiv",
        summary: "Cari artworks/manga/novel di Pixiv.",
        auth: "apikey",
        params: [
          { name: "q", loc: "query", required: true, example: "miku", desc: "Kata kunci" },
          { name: "type", loc: "query", required: false, example: "artworks", desc: "artworks | manga | novel" },
        ],
        response: pixivResponse,
      },
    ],
  },
];

export function findEndpoint(id: string): Endpoint | undefined {
  for (const g of GROUPS) {
    const e = g.endpoints.find((x) => x.id === id);
    if (e) return e;
  }
  return undefined;
}
