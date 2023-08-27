<template>
    <div class="w-full p-5 bg-ll-neutral dark:bg-ld-neutral rounded-md flex flex-col mb-4">
        <div class="flex justify-between">
            <div class="flex items-center">
                <!-- 貼文中的大頭貼 -->
                <div class="avatar rounded-full bg-ll-base dark:bg-ld-base w-15 h-15 border-2 border-ll-border 
                            dark:border-ld-border relative ">
                    <img 
                        src="https://picsum.photos/seed/picsum/200/300"
                        class="w-full h-full  rounded-full object-cover" 
                        alt="">
                </div>
                <!-- 貼文中的用戶名 與 ID -->
                <div class="flex flex-col ml-2">
                    <p class="text-2xl font-bold text-gray-800 dark:text-gray-300">{{ props.post.author.name }}</p>
                    <p class="-mt-1">@{{ props.post.author.id }}</p>
                </div>
                <TinyBlueTickLabel text="1h" />
            </div>
            <button class="active:scale-95 transform transition-transform">
                <PostSettingsSvg/>
            </button>
        </div>

        <div v-if="props.post.pictures.length > 0"
            :class="`images w-full h-70 bg-ll-neutral dark:bg-ld-neutral rounded-xl my-4 overflow-hidden 
            grid ${(props.post.pictures.length > 1) ? 'grid-cols-2' : 'grid-cols-1'} gap-2`">
            <div class="h-full">
                <img :src="`${props.post.pictures[0]}`" class="w-full h-70   object-cover" alt="">
            </div>
            <div v-if="props.post.pictures.length > 1" :class="`
            
            h-70 grid ${props.post.pictures.length == 2 ? 'grid-cols-1 grid-rows-1' : ''} 
             ${props.post.pictures.length == 3 ? 'grid-cols-1 grid-rows-2' : ''} 
            ${props.post.pictures.length == 4 ? 'grid-cols-2 grid-rows-2' : ''} 
            

            gap-2`">
                <img v-if="props.post.pictures.length > 1" :src="`${props.post.pictures[1]}`"
                    :class="`w-full h-full object-cover 
                    ${props.post.pictures.length == 3 && 'row-span-1 col-span-1 h-full'}`"
                    alt="">
                <img v-if="props.post.pictures.length > 2" :src="`${props.post.pictures[2]}`"
                    :class="`w-full h-full object-cover ${props.post.pictures.length == 3 && 'row-span-2 col-span-1'}`"
                    alt="">
                <img v-if="props.post.pictures.length > 3" :src="`${props.post.pictures[3]}`"
                    :class="`w-full h-full object-cover ${props.post.pictures.length == 4 && 'col-span-2'}`" alt="">
                <img v-if="props.post.pictures.length > 4" :src="`${props.post.pictures[4]}`"
                    :class="`w-full h-2/4 object-cover ${props.post.pictures.length == 5 && 'col-span-3 row-span-1'}`"
                    alt="">
            </div>
        </div>
        <p v-html="generateDescription()" :class="`${props.post.pictures.length == 0 ? ' my-4 text-xl' : ''}`"></p>

        <div class="flex justify-between pt-4 border-t border-ll-border dark:border-ld-border mt-4">
            <ReplyButton :number="post.replys.length"/>
            <ShareButton :number="post.shares.length"/>
            <LikeButton :number="post.likes.length"/>
        </div>
    </div>
</template>

<script setup lang="ts">
import { PropType, ref } from 'vue'
import { IPost } from '../types'
import TinyBlueTickLabel from './Label/TinyBlueTick.vue'
import PostSettingsSvg from './Icon/PostSettingsSvg.vue'
import ReplyButton from './Button/ReplyButton.vue'
import ShareButton from './Button/ShareButton.vue'
import LikeButton from './Button/LikeButton.vue'

// 定義 props
const props = defineProps({
    post: {
        type: Object as PropType<IPost>,
        required: true
    }
})

// 生成描述文字的 HTML
function generateDescription() {
    let description = props.post.description.trim().split('\n').join('<br>');
    description = description.replace(/#(\S*)/g, '<a class="text-ll-primary" href="/search/$1">#$1</a>');
    return description;
}
</script>

<script lang="ts">
// 這裡可以放置在 script setup 之外的 Vue3 Composition API 相關代碼
// 如果你有需要的話
export default {
    // 在這裡可以添加 Vue3 的選項
}
</script>

<style lang="">
</style>