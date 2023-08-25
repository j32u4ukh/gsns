// 引入 defineConfig 函數從 vite 模塊
import { defineConfig } from 'vite'

// 引入 vite-plugin-vue 插件，用於處理 Vue 組件
import vue from '@vitejs/plugin-vue'

// 引入 vite-plugin-windicss 插件，用於處理 WindiCSS 樣式
import WindiCSS from 'vite-plugin-windicss'

// 定義 Vite 配置
// https://vitejs.dev/config/
export default defineConfig({
  // 配置插件，包括 vue 和 WindiCSS
  plugins: [
    // 使用 vue 插件處理 Vue 組件
    vue(), 
     // 使用 WindiCSS 插件處理 WindiCSS 樣式
    WindiCSS()
  ],
  server: {
    host: true,
    port: 8080,
    // 熱重載
    hmr: true, 
  },
})
