/// <reference types="vite/client" />

// 聲明一個模塊，用於處理 Vue 單文件組件的類型聲明
declare module '*.vue' {
  import type { DefineComponent } from 'vue'

  // 聲明一個常量 component，它是一個 Vue 組件
  const component: DefineComponent<{}, {}, any>

  // 導出這個 Vue 組件，以便在其他地方使用
  export default component
}
