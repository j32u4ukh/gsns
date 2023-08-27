<script setup lang="ts">
import { reactive } from 'vue'
import TinyBlueTickSvg from './Icon/TinyBlueTick.vue'

// 使用 reactive 創建一個響應式數據對象 trendings，包含一組趨勢對象
const trendings = reactive([
    {
        index: 0,
        gender: 'Entertainment',
        label: '#DCcomics',
        stats: '87,9 mil Posts'
    },
    {
        index: 0,
        gender: 'Entertainment',
        label: '#Marvel',
        stats: '75,9 mil Posts'
    },
    {
        index: 0,
        gender: 'Anime',
        label: '#Naruto',
        stats: '25,9 mil Posts'
    },

])


</script>
<template>
    <div class="w-full h-full flex flex-col py-2 px-2 bg-ll-base dark:bg-ld-base overflow-y-auto">
        <!-- 一個包含趨勢列表的樣式塊 -->
        <ul class="w-full bg-ll-neutral dark:bg-ld-neutral p-5 rounded-md mb-5">
            <li class="dark:text-gray-200 text-gray-700 text-lg  border-b pb-3 border-ll-border dark:border-ld-border">
                <p>Worldwide Trending</p>
            </li>
            <!-- 單個趨勢項 -->
            <li class="mb-4">
                <p class="text-lg dark:text-gray-200 text-gray-700 font-bold">#Desencantada</p>
                <p class="text-xs">Original film, Watch now on HBO Max </p>
            </li>

            <!-- 使用 v-for 循環渲染 trendings 數組中的每個趨勢項 -->
            <li class="mb-4" v-for="(trending, index) in trendings" :key="index">
                <p class="text-xs">{{ index + 1 }} - {{ trending.gender }} </p>
                <p class="text-lg dark:text-gray-200 text-gray-700 font-bold">{{ trending.label }}</p>
                <p class="text-xs">{{ trending.stats }}</p>
            </li>
        </ul>

        <!-- 另一個包含“Who to follow”列表的樣式塊 -->
        <ul class="w-full bg-ll-neutral dark:bg-ld-neutral p-5 rounded-md mb-5">
            <li
                class="dark:text-gray-200 text-gray-700 text-lg  border-b pb-3 border-ll-border dark:border-ld-border mb-4">
                <p>Who to follow</p>
            </li>

            <!-- 使用 v-for 循環渲染一些元素 -->
            <div v-for="(item) in [0, 1, 2, 3]" :key="item" class="flex items-center mb-4">
                <!-- 頭像 -->
                <div
                    class="avatar rounded-full bg-ll-base dark:bg-ld-base w-15 h-15 border-2 border-ll-border dark:border-ld-border relative ">
                    <img :src="`https://picsum.photos/id/${21 + item}/200/200`"
                        class="w-full h-full  rounded-full object-cover" alt="">
                </div>

                <!-- 用戶信息 -->
                <div class="flex flex-col ml-2">
                    <p class="text-lg font-bold text-gray-800 dark:text-gray-300">Lukebana</p>
                    <p class="-mt-1 text-xs">@luk{{item}}d</p>
                </div>
                <TinyBlueTickSvg fill="blue"/>
            </div>
        </ul>

    </div>
</template>
<script lang="ts">
export default {

}
</script>
