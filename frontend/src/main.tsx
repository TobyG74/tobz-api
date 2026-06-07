import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Toaster } from "sonner";
import { ThemeProvider } from "./context/ThemeContext";
import { AuthProvider } from "./context/AuthContext";
import App from "./App";
import "./index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider>
      <AuthProvider>
        <App />
        <div className="grain" aria-hidden />
        <Toaster
          position="bottom-right"
          toastOptions={{
            style: {
              background: "var(--color-surface)",
              border: "1px solid var(--color-line)",
              color: "var(--color-paper)",
              fontFamily: "Sora, sans-serif",
            },
          }}
        />
      </AuthProvider>
    </ThemeProvider>
  </StrictMode>
);
