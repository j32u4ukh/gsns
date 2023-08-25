<template> 
    <!-- Menu2 圖示 -->
    <svg 
        xmlns="http://www.w3.org/2000/svg" 
        :width="width" 
        :height="height" 
        fill="none"
        stroke-width="1.5"
        :stroke="color" 
        class="w-6 h-6">
        <path 
            stroke-linecap="round" 
            stroke-linejoin="round"
            d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12" />
    </svg>
</template>
  
<script setup lang="ts">
const props = defineProps({
    width: { type: String, default: '24' },
    height: { type: String, default: '24' },
    color: { type: String, default: 'currentColor' }
})
</script>
