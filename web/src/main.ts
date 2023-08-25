// 引入 Vue 的 createApp 函數，用於創建 Vue 應用實例
import { createApp } from 'vue'
import router from "./router";
// 引入 'virtual:windi.css'，這是一個用於動態生成 Tailwind CSS 類的插件
import 'virtual:windi.css'

// 引入自定義的樣式文件
import './style.css'
import { createPinia } from "pinia";
// 引入 App 組件，這是根組件，將被掛載到頁面上
import App from './App.vue'

// 使用 createApp 函數創建 Vue 應用實例，並將 App 組件掛載到 #app 元素上
createApp(App).use(router).use(createPinia()).mount('#app')
