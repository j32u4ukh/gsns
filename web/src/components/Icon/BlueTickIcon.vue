<template> 
    <div 
        :class="`${isExpanded ? 'w-8 h-8 border-4 right-0 bottom-3' : 'w-5 h-5 border-2 -right-2 bottom-0'}  
                bg-ll-primary rounded-full absolute border-ll-base dark:border-ld-base text-white flex justify-center items-center`">
        <svg 
            xmlns="http://www.w3.org/2000/svg" 
            :width="width" 
            :height="height" 
            :stroke="color" 
            fill="none" 
            stroke-width="1.5"
            class="w-5 h-5">
            <path 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                d="M4.5 12.75l6 6 9-13.5" />
        </svg>
    </div>
</template>
  
<script setup lang="ts">
const props = defineProps({
    width: { type: String, default: '24' },
    height: { type: String, default: '24' },
    color: { type: String, default: 'currentColor' },
    isExpanded: Boolean,
})
</script>
