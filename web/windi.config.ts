import { defineConfig } from 'windicss/helpers'

// 定義 WindiCSS 的配置
export default defineConfig({
    // 設置暗黑模式，可選值有 'class', 'media', 'false'
    darkMode: 'class',
    // safelist: 'p-3 p-4 p-5',
    // 主題配置
    theme: {
        // 擴展主題
        extend: {
            // 自定義顏色
            colors: {
                // 淺色主題顏色
                ll: {
                    base: "#00000014",
                    neutral: "#FFFFFF",
                    primary: "#2697FE",
                    secondary: "#1EE9AC",
                    accent: "#FFD025",
                    info: "#1CB5E7",
                    success: "#2BD3A2",
                    warning: "#FAB34A",
                    error: "#FE4443",
                    border: "#EAEBEF"
                },
                // 深色主題顏色
                ld: {
                    base: "#212332",
                    neutral: "#2A2E3F",
                    primary: "#2697FE",
                    secondary: "#1EE9AC",
                    accent: "#FFD025",
                    info: "#1CB5E7",
                    success: "#2BD3A2",
                    warning: "#FAB34A",
                    error: "#FE4443",
                    border: "#111219"
                },

            },
        },
    },
    // 插件配置
    plugins: [],
})