import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// 正确配置结构
export default defineConfig({
  assetsInclude: ["**/*.md", "**/*.png"], // 新增配置
  plugins: [react()],
  server: {
    // 开发服务器配置
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path,
      },
    },
  },
  base: "/", // 生产环境基础路径
  build: {
    // 构建配置（需与server同级）
    outDir: "dist",
    assetsDir: "assets",
    emptyOutDir: true,
  },
});
