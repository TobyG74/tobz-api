import type { Lang } from "./endpoints";

export interface ReqDesc {
  method: string;
  url: string; // base URL without query string
  query: Record<string, string>;
  body?: Record<string, string>;
  headers: Record<string, string>;
}

export const LANGS: { id: Lang; label: string }[] = [
  { id: "curl", label: "cURL" },
  { id: "python", label: "Python" },
  { id: "node", label: "Node.js" },
  { id: "go", label: "Go" },
  { id: "ruby", label: "Ruby" },
  { id: "php", label: "PHP" },
];

function qs(query: Record<string, string>): string {
  const entries = Object.entries(query).filter(([, v]) => v !== "");
  if (!entries.length) return "";
  return "?" + entries.map(([k, v]) => `${k}=${encodeURIComponent(v)}`).join("&");
}

const fullURL = (r: ReqDesc) => r.url + qs(r.query);

export function generate(lang: Lang, r: ReqDesc): string {
  switch (lang) {
    case "curl":
      return curl(r);
    case "python":
      return python(r);
    case "node":
      return node(r);
    case "go":
      return go(r);
    case "ruby":
      return ruby(r);
    case "php":
      return php(r);
  }
}

function curl(r: ReqDesc): string {
  const lines = [`curl -X ${r.method} "${fullURL(r)}"`];
  for (const [k, v] of Object.entries(r.headers)) lines.push(`  -H "${k}: ${v}"`);
  if (r.body) {
    lines.push(`  -H "Content-Type: application/json"`);
    lines.push(`  -d '${JSON.stringify(r.body)}'`);
  }
  return lines.join(" \\\n");
}

function python(r: ReqDesc): string {
  const lines = ["import requests", "", `url = "${r.url}"`];
  if (Object.keys(r.query).length) lines.push(`params = ${pyDict(r.query)}`);
  if (Object.keys(r.headers).length) lines.push(`headers = ${pyDict(r.headers)}`);
  if (r.body) lines.push(`payload = ${pyDict(r.body)}`);

  const args = ["url"];
  if (Object.keys(r.query).length) args.push("params=params");
  if (r.body) args.push("json=payload");
  if (Object.keys(r.headers).length) args.push("headers=headers");

  lines.push("", `resp = requests.${r.method.toLowerCase()}(${args.join(", ")})`, "print(resp.json())");
  return lines.join("\n");
}

function node(r: ReqDesc): string {
  const opts: string[] = [`method: "${r.method}"`];
  if (Object.keys(r.headers).length || r.body) {
    const h = { ...r.headers };
    if (r.body) h["Content-Type"] = "application/json";
    opts.push(`headers: ${jsObj(h)}`);
  }
  if (r.body) opts.push(`body: JSON.stringify(${jsObj(r.body)})`);
  return [
    `const res = await fetch("${fullURL(r)}", {`,
    `  ${opts.join(",\n  ")}`,
    `});`,
    `const data = await res.json();`,
    `console.log(data);`,
  ].join("\n");
}

function go(r: ReqDesc): string {
  const imports = r.body
    ? `import (\n\t"fmt"\n\t"io"\n\t"net/http"\n\t"strings"\n)`
    : `import (\n\t"fmt"\n\t"io"\n\t"net/http"\n)`;
  const reqLine = r.body
    ? `req, _ := http.NewRequest("${r.method}", "${fullURL(r)}", strings.NewReader(\`${JSON.stringify(r.body)}\`))`
    : `req, _ := http.NewRequest("${r.method}", "${fullURL(r)}", nil)`;
  const headerLines = Object.entries(r.headers).map(([k, v]) => `\treq.Header.Set("${k}", "${v}")`);
  if (r.body) headerLines.push(`\treq.Header.Set("Content-Type", "application/json")`);
  return [
    "package main",
    "",
    imports,
    "",
    "func main() {",
    `\t${reqLine}`,
    ...headerLines,
    "\tres, _ := http.DefaultClient.Do(req)",
    "\tdefer res.Body.Close()",
    "\tb, _ := io.ReadAll(res.Body)",
    "\tfmt.Println(string(b))",
    "}",
  ].join("\n");
}

function ruby(r: ReqDesc): string {
  const lines = [`require "net/http"`, `require "json"`, "", `uri = URI("${fullURL(r)}")`];
  const reqClass =
    r.method === "GET" ? "Get" : r.method === "POST" ? "Post" : r.method.charAt(0) + r.method.slice(1).toLowerCase();
  lines.push(`req = Net::HTTP::${reqClass}.new(uri)`);
  for (const [k, v] of Object.entries(r.headers)) lines.push(`req["${k}"] = "${v}"`);
  if (r.body) {
    lines.push(`req["Content-Type"] = "application/json"`);
    lines.push(`req.body = ${JSON.stringify(JSON.stringify(r.body))}`);
  }
  lines.push(
    `res = Net::HTTP.start(uri.hostname, uri.port, use_ssl: uri.scheme == "https") { |http| http.request(req) }`,
    "puts res.body"
  );
  return lines.join("\n");
}

function php(r: ReqDesc): string {
  const headers = Object.entries(r.headers).map(([k, v]) => `"${k}: ${v}"`);
  if (r.body) headers.push(`"Content-Type: application/json"`);
  const lines = [
    "<?php",
    `$ch = curl_init("${fullURL(r)}");`,
    "curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);",
    `curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "${r.method}");`,
  ];
  if (headers.length) lines.push(`curl_setopt($ch, CURLOPT_HTTPHEADER, [${headers.join(", ")}]);`);
  if (r.body) lines.push(`curl_setopt($ch, CURLOPT_POSTFIELDS, '${JSON.stringify(r.body)}');`);
  lines.push("$response = curl_exec($ch);", "curl_close($ch);", "echo $response;");
  return lines.join("\n");
}

// --- formatting helpers ---
function jsObj(o: Record<string, string>): string {
  const entries = Object.entries(o);
  if (!entries.length) return "{}";
  return "{ " + entries.map(([k, v]) => `"${k}": "${v}"`).join(", ") + " }";
}
function pyDict(o: Record<string, string>): string {
  const entries = Object.entries(o);
  if (!entries.length) return "{}";
  return "{" + entries.map(([k, v]) => `"${k}": "${v}"`).join(", ") + "}";
}
