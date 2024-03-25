import { defineConfig } from "vite";
import viteTsconfigPaths from "vite-tsconfig-paths";

import federation from "@originjs/vite-plugin-federation";
import react from "@vitejs/plugin-react-swc";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    viteTsconfigPaths(),
    federation({
      name: "minion",
      filename: "remote.js",
      exposes: {
        "./App": "./src/pages/app.tsx",
      },
      shared: [
        "react",
        "react-dom",
        "react-router-dom",
        "axios",
        "react-truncate-inside",
        "dayjs",
        "radash",
        "react-helmet-async",
        "react-hook-form",
        "@mui/material",
        "@mui/icons-material",
        "@tanstack/react-query",
      ],
    }),
  ],
  build: {
    target: "esnext", //browsers can handle the latest ES features
    outDir: "../static",
  },
  server: {
    proxy: {
      "/api/minion": {
        target: "http://localhost:59010",
        changeOrigin: true,
        secure: false,
        ws: true,
        rewrite: (path) => path.replace(/^\/api\/minion/, ""),
      },
    },
  },
});
