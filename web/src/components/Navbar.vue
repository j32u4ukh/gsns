<script setup lang="ts">
import { ref, markRaw } from 'vue'
import BlueTickIcon from './Icon/BlueTickIcon.vue'
import PencilSvg from './Icon/PencilSvg.vue'
import HomeSvg from './Icon/HomeSvg.vue';
import ExploreSvg from './Icon/ExploreSvg.vue';
import NotificationsSvg from './Icon/NotificationsSvg.vue';
import FavoriteSvg from './Icon/FavoriteSvg.vue';
import DirectSvg from './Icon/DirectSvg.vue';
import StarsSvg from './Icon/StarsSvg.vue';
import SettingsSvg from './Icon/SettingsSvg.vue';
import BriefInfo from "./Label/BriefInfo.vue";

const nPost = ref(3100);
const nFollower = ref(298451);
const nFollowing = ref(20500000);

// 菜單列表
const menus = ref({
    // 當前激活的菜單索引
    active: ref(0), // 使用 ref 包裹以保持響應性
    menusList: [
        // 菜單項的信息
        {
            index: 0,
            name: "Home",
            icon: markRaw(HomeSvg),
            stroke: "currentColor",
        },
        {
            index: 1,
            name: "Explore",
            icon: markRaw(ExploreSvg),
            stroke: "currentColor",
        },
        {
            index: 6,
            name: "Notifications",
            icon: markRaw(NotificationsSvg),
            stroke: "currentColor",
        },
        {
            index: 2,
            name: "Favorite",
            icon: markRaw(FavoriteSvg),
            stroke: "currentColor",
        },
        {
            index: 3,
            name: "Direct",
            icon: markRaw(DirectSvg),
            stroke: "currentColor",
        },
        {
            index: 4,
            name: "Stars",
            icon: markRaw(StarsSvg),
            stroke: "currentColor",
        },
        {
            index: 5,
            name: "Settings",
            icon: markRaw(SettingsSvg),
            stroke: "currentColor",
        },
    ]
});

// 定義組件的 props，此處使用 defineProps
const props = defineProps({
    isExpanded: {
        type: Boolean,
        required: true
    }
})
</script>

<template>
    <div
        :class="`w-full h-full flex flex-col ${props.isExpanded ? 'p-10 px-5' : 'p-2'} 
                relative overflow-y-auto overflow-x-hidden`">
        <div class="profile flex flex-col justify-center items-center">
            <div :class="`avatar rounded-full bg-ll-base dark:bg-ld-base ${props.isExpanded ? 'w-25 h-25' : 'w-12 h-12'} 
                        border-2 border-ll-border dark:border-ld-border relative`">
                <!-- 大頭貼圖片 -->
                <img 
                    src="https://picsum.photos/seed/picsum/200/300" 
                    class="w-full h-full  rounded-full object-cover" 
                    alt="">
                <!-- 藍勾勾圖示 -->
                <BlueTickIcon :isExpanded="isExpanded"/>
            </div>
            <p v-if="props.isExpanded" class="text-xl font-bold text-gray-800 dark:text-gray-300">Lukebana</p>
            <p class="-mt-1 text-sm" v-if="props.isExpanded">The creator of this platform</p>
        </div>

        <div v-if="isExpanded"
            class="w-full flex justify-between mt-5 pb-5 border-b border-ll-border dark:border-ld-border">
            <!-- <div class="flex flex-col justify-center items-center">
                <p class="text-lg font-bold text-gray-800 dark:text-gray-300">255</p>
                <p class="-mt-1 text-xs">Posts</p>
            </div> -->
            <BriefInfo label="Posts" :number=nPost />
            <!-- <div class="flex flex-col justify-center items-center">
                <p class="text-lg font-bold text-gray-800 dark:text-gray-300">298.45K</p>
                <p class="-mt-1 text-xs">Followers</p>
            </div> -->
            <BriefInfo label="Followers" :number=nFollower />
            <!-- <div class="flex flex-col justify-center items-center">
                <p class="text-lg font-bold text-gray-800 dark:text-gray-300">20.5M</p>
                <p class="-mt-1 text-xs">Following</p>
            </div> -->
            <BriefInfo label="Following" :number=nFollowing />
        </div>
        <ul :class="`flex flex-col w-full pt-5 ${props.isExpanded ? '' : 'justify-center flex '}`">
            <li v-for="(menu, index) in menus.menusList" :key="menu.name"
                :class="`w-full py-2  flex items-center ${props.isExpanded ? 'mb-2' : 'justify-center mb-4'} 
                        ${menu.index == menus.active ? 'text-ll-primary' : ''} cursor-pointer active:scale-95 
                        transform transition-transform select-none`"
                @click="menus.active = menu.index, $.emit('onCloseNavbar', false)">
                <!-- <div v-html="menu.icon"></div>
                <p v-if="props.isExpanded" class="ml-5 text-sm">{{ menu.name }}</p> -->

                <component :is="menu.icon" :stroke="menu.stroke" />
                <p v-if="props.isExpanded" class="ml-5 text-sm">{{ menu.name }}</p>
            </li>
        </ul>

        <button @click="$.emit('onComposePost'), $.emit('onCloseNavbar', false)"
            class="bg-ll-primary dark:bg-ld-primary text-white rounded-lg py-3 px-2 active:scale-95 transform transition-transform flex items-center justify-center">
            <p v-if="props.isExpanded">Share on Space</p>
            <PencilSvg :isExpanded="isExpanded"/>
        </button>
    </div>
</template>
<script lang="ts">
// 在這裏可以添加其他組件邏輯，但您已經在 <script setup> 中導入了所需的內容
// 所以這部分通常不需要額外的代碼
export default {

}
</script>
