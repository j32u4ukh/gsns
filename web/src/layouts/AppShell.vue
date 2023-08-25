<template >
    <div :class="`${darkmode ? 'dark' : 'light'}`">
        <!-- 頁面主要的布局結構 -->
        <div class="w-screen h-screen bg-ll-base dark:bg-ld-base flex flex-col text-gray-500 ">
            <!-- 頁面頂部區域 -->
            <div v-if="!props.fullSidebar"
                class="w-full h-14 bg-ll-neutral dark:bg-ld-neutral border-ll-border dark:border-ld-border border-b flex justify-between items-center px-5">
                <slot name="header"></slot>
                <button 
                    @click="changeMode"
                    class="w-10 h-10 border rounded-md flex justify-center items-center ml-2 border-ll-border dark:border-ld-border bg-ll-base dark:bg-ld-base dark:text-gray-200 active:scale-95 transition-transform transform">
                    <!-- 用於切換暗色模式的圖示 -->
                    <DarkModeSvg :darkmode="darkmode"/>
                </button>
            </div>

            <!-- 頁面主要區域 -->
            <div class="w-full flex h-full relative overflow-hidden relative ">
                <!-- 左側導航欄 -->
                <div :class="`absolute left-0 top-0 z-10 w-full md:relative origin-left overflow-x-hidden ${props.navbarExpanded ? 'md:w-110' : 'w-0 md:w-20'} transition-all  border-r h-full bg-ll-neutral dark:bg-ld-neutral border-ll-border dark:border-ld-border flex flex-col`">
                    <slot name="navbar"></slot>
                </div>

                <!-- 主內容區域 -->
                <div class="w-full h-full flex flex-col">
                    <!-- 若使用全尺寸側邊欄，則在頁面上方顯示頂部區域 -->
                    <div v-if="props.fullSidebar"
                        class="w-full h-14 bg-ll-neutral dark:bg-ld-neutral border-ll-border dark:border-ld-border border-b flex justify-between items-center px-5">
                        <slot name="header"></slot>
                        <button @click="changeMode"
                            class="w-10 h-10 border rounded-md flex justify-center items-center ml-2 border-ll-border dark:border-ld-border bg-ll-base dark:bg-ld-base dark:text-gray-200 active:scale-95 transition-transform transform">
                            <!-- 用於切換暗色模式的圖示 -->
                            <DarkModeSvg :darkmode="darkmode"/>
                        </button>
                    </div>

                    <!-- 主要內容區域 -->
                    <div class="w-full h-full flex flex-col overflow-auto">
                        <slot name="body"></slot>
                    </div>
                </div>

                <!-- 右側導航欄 -->
                <div v-if="$slots.rightNavbar"
                    :class="`origin-left overflow-x-hidden ${props.rightNavbarExpanded ? 'w-130' : 'w-0'} transition-all border-l h-full bg-ll-neutral dark:bg-ld-neutral border-ll-border dark:border-ld-border flex flex-col`">
                    <slot name="rightNavbar"></slot>
                </div>

            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, PropType } from 'vue';
import DarkModeSvg from '../components/Icon/DarkModeSvg.vue'

// 創建一個 reactive 變數來控制暗色模式
let darkmode = ref(false);

// 使用 defineProps 定義父組件傳遞的 props，並指定其類型
const props = defineProps({
    // 是否擴展側邊欄
    fullSidebar: Boolean,
    // 是否擴展導航欄
    navbarExpanded: Boolean as PropType<boolean>,
    // 是否擴展右側導航欄
    rightNavbarExpanded: Boolean as PropType<boolean>
})

const changeMode = () => {
    darkmode.value = !darkmode.value;
    // $emit('onChangeTheme', darkmode);
    console.log('darkmode:' + darkmode.value);
}
</script>

<script lang="ts">
export default {

}
</script>

<style lang="">
    
</style>