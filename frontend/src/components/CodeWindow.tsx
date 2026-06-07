import { motion } from "motion/react";

/** A faux terminal/editor card showing a sample API request → response. */
export function CodeWindow() {
  return (
    <motion.div
      initial={{ opacity: 0, y: 24, rotateX: 8 }}
      animate={{ opacity: 1, y: 0, rotateX: 0 }}
      transition={{ duration: 0.7, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
      className="grad-border overflow-hidden rounded-2xl text-left shadow-2xl"
    >
      {/* title bar */}
      <div className="flex items-center gap-2 border-b border-line bg-void/60 px-4 py-3">
        <span className="h-3 w-3 rounded-full bg-[#ff5f57]" />
        <span className="h-3 w-3 rounded-full bg-[#febc2e]" />
        <span className="h-3 w-3 rounded-full bg-[#28c840]" />
        <span className="ml-3 font-mono text-xs text-mist">~ tobz-api · request</span>
      </div>

      {/* request */}
      <pre className="overflow-x-auto px-5 py-4 font-mono text-[13px] leading-relaxed">
        <code>
          <span className="text-cyan">curl</span>{" "}
          <span className="text-fog">https://api.tobz.dev/v1/download</span> <span className="text-mist">\</span>
          {"\n  "}
          <span className="text-violet">-H</span> <span className="text-paper">"X-API-Key: tobz_live_••••"</span>{" "}
          <span className="text-mist">\</span>
          {"\n  "}
          <span className="text-violet">--data-urlencode</span>{" "}
          <span className="text-paper">"url=https://tiktok.com/@user/video/123"</span>
        </code>
      </pre>

      {/* response */}
      <div className="border-t border-line bg-void/40 px-5 py-4">
        <div className="mb-2 flex items-center gap-2 font-mono text-[11px] uppercase tracking-widest text-mist">
          <span className="h-1.5 w-1.5 rounded-full bg-[#28c840]" /> 200 OK · 142ms
        </div>
        <pre className="overflow-x-auto font-mono text-[13px] leading-relaxed">
          <code>
            <span className="text-mist">{"{"}</span>
            {"\n  "}
            <span className="text-azure">"platform"</span>
            <span className="text-mist">: </span>
            <span className="text-paper">"tiktok"</span>
            <span className="text-mist">,</span>
            {"\n  "}
            <span className="text-azure">"download_items"</span>
            <span className="text-mist">: [{"{"}</span>
            {"\n    "}
            <span className="text-azure">"label"</span>
            <span className="text-mist">: </span>
            <span className="text-paper">"Video HD"</span>
            <span className="text-mist">, </span>
            <span className="text-azure">"url"</span>
            <span className="text-mist">: </span>
            <span className="text-violet">"https://…/hd.mp4"</span>
            {"\n  "}
            <span className="text-mist">{"}]"}</span>
            {"\n"}
            <span className="text-mist">{"}"}</span>
          </code>
        </pre>
      </div>
    </motion.div>
  );
}
