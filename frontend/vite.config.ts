import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  // Port 3000 matches the backend's default ALLOWED_ORIGINS.
  server: { port: 3000 },
  preview: { port: 3000 },
});
