<template>
    <svg 
        xmlns="http://www.w3.org/2000/svg" 
        :width="width" 
        :height="height"  
        :fill="fill"
        stroke-width="1.5"
        stroke="currentColor" 
        class="w-full h-full">
        <path 
            stroke-linecap="round" 
            stroke-linejoin="round"
            d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
</template>
  
<script setup lang="ts">
const props = defineProps({
    width: { type: String, default: '24' },
    height: { type: String, default: '24' },
    fill: { type: String, default: 'none' }
})
</script>
