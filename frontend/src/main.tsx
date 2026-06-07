import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Toaster } from "sonner";
import { AuthProvider } from "./context/AuthContext";
import App from "./App";
import "./index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AuthProvider>
      <App />
      <div className="grain" aria-hidden />
      <Toaster
        theme="dark"
        position="bottom-right"
        toastOptions={{
          style: {
            background: "#0d1224",
            border: "1px solid #1e2547",
            color: "#eaecff",
            fontFamily: "Sora, sans-serif",
          },
        }}
      />
    </AuthProvider>
  </StrictMode>
);
