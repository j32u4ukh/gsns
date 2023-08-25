<template >
    <div class="flex w-full h-full justify-between items-center">
        <!-- 左側菜單按鈕 -->
        <div class="flex items-center">
            <button @click="$emit('onMenuClick')"
                class="w-10 h-10 mr-2 border rounded-md flex justify-center items-center ml-2 border-ll-border dark:border-ld-border bg-ll-base dark:bg-ld-base dark:text-gray-200 active:scale-95 transition-transform transform">
                <Menu2Svg/>
            </button>

            <!-- 網站名稱 -->
            <p class="text-xl font-normal font-bold">Hivy.Space</p>
        </div>

        <!-- 搜索框和搜索按鈕 -->
        <div
            class="bg-ll-base dakr:bg-ld-base rounded-md border  border-ll-border dark:border-ld-border w-2/4 h-3/4 flex items-center">
            <input type="text" placeholder="Search on hivy..."
                class="px-2 bg-transparent text-md outline-none w-full h-full">

            <button class="p-2 active:scale-95 transform transition-transform">
                <!-- 放大鏡圖示 -->
                <SearchSvg/>
            </button>
        </div>

        <!-- 右側菜單按鈕 -->
        <div>
            <button @click="$emit('onRightMenuClick')"
                class="w-10 h-10 mr-2 border rounded-md flex justify-center items-center ml-2 border-ll-border dark:border-ld-border bg-ll-base dark:bg-ld-base dark:text-gray-200 active:scale-95 transition-transform transform">
                <Menu1Svg/>
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import SearchSvg from './Icon/SearchSvg.vue'
import Menu1Svg from './Icon/Menu1Svg.vue'
import Menu2Svg from './Icon/Menu2Svg.vue'
</script>
