<template>
    <svg 
        xmlns="http://www.w3.org/2000/svg" 
        :width="width" 
        :height="height" 
        :fill="fill"
        class="w-6 h-6">
        <path 
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M4.5 12a1.5 1.5 0 113 0 1.5 1.5 0 01-3 0zm6 0a1.5 1.5 0 113 0 1.5 1.5 0 01-3 0zm6 0a1.5 1.5 0 113 0 1.5 1.5 0 01-3 0z"/>
    </svg>
</template>
  
<script setup lang="ts">
const props = defineProps({
    width: { type: String, default: '24' },
    height: { type: String, default: '24' },
    fill: { type: String, default: 'currentColor' }
})
</script>
