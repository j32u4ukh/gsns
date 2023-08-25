<template>
    <div class="flex flex-col justify-center items-center">
        <p class="text-lg font-bold text-gray-800 dark:text-gray-300">{{formatNumber(number)}}</p>
        <p class="-mt-1 text-xs">{{label}}</p>
    </div>
</template>
  
<script setup lang="ts">
import {formatNumber} from "../../utils"
const props = defineProps({
    label: String,
    number: { type: Number, default: 0 },
})
</script>