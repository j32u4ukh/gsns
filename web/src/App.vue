<!-- <script setup lang="ts">
// 此起始模板使用了 Vue 3 的 <script setup> 單文件組件 (SFCs)
// 了解更多信息：https://vuejs.org/api/sfc-script-setup.html#script-setup
import Main from './features/main/Main.vue'
</script> -->

<!-- <template> -->
  <!-- 使用 Main 组件 -->
  <!-- <Main></Main>
</template> -->

<!-- <script lang="ts">
// 導出默認配置
export default {
  components: { Main }
}
</script> -->

<template>
  <!-- RouterView 是由 Vue Router 內部定義的一個特殊組件。它在 Vue Router 中扮演著非常重要的角色，
    用於根據當前路由匹配動態渲染相應的組件。 -->
  <RouterView />
</template>


