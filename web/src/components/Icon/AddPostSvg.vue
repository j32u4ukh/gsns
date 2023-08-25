<template>    
    <svg 
        xmlns="http://www.w3.org/2000/svg" 
        :width="width" 
        :height="height" 
        :fill="fill"
        stroke-width="1.5" 
        stroke="currentColor" 
        class="w-6 h-6 mr-2">
        <path 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" />
    </svg>
</template>
  
<script setup lang="ts">
const props = defineProps({
    width: { type: String, default: '24' },
    height: { type: String, default: '24' },
    fill: { type: String, default: 'currentColor' }
})
</script>
