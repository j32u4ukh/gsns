<template>
    <svg 
        xmlns="http://www.w3.org/2000/svg" 
        :width="width" 
        :height="height" 
        :stroke="stroke" 
        fill="none" 
        stroke-width="1.5" 
        class="w-6 h-6">
        <path 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
    </svg>
</template>
  
<script setup lang="ts">
const props = defineProps({
    width: { type: String, default: '24' },
    height: { type: String, default: '24' },
    stroke: { type: String, default: 'currentColor' }
})
</script>
