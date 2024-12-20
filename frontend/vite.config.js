import { defineConfig } from "vite";

export default defineConfig({
  server: {
    proxy: {
      "/api": {
        target: "https://localhost:8443", // URL вашего Go сервера
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
