<template> 
    <!-- Menu1 圖示 -->
    <svg 
        xmlns="http://www.w3.org/2000/svg" 
        :width="width" 
        :height="height" 
        :fill="color"
        stroke-width="1.5"
        stroke="currentColor" 
        class="w-6 h-6">
        <path 
            stroke-linecap="round" 
            stroke-linejoin="round"
            d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 010 3.75H5.625a1.875 1.875 0 010-3.75z" />
    </svg>
</template>
  
<script setup lang="ts">
const props = defineProps({
    width: { type: String, default: '24' },
    height: { type: String, default: '24' },
    color: { type: String, default: 'none' }
})
</script>
