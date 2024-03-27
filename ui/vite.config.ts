import { defineConfig } from "vite";
import viteTsconfigPaths from "vite-tsconfig-paths";

import federation from "@originjs/vite-plugin-federation";
import react from "@vitejs/plugin-react-swc";

// @ts-ignore
import pkg from "./package.json";
const { dependencies } = pkg;

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
      shared: {
        ...dependencies,
        react: {
          requiredVersion: dependencies["react"],
        },
        "react-dom": {
          requiredVersion: dependencies["react-dom"],
        },
      },
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
