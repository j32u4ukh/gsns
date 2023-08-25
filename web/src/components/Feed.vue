<template >
    <div :class="`grid-cols-1 w-full grid ${props.oneColumn ? 'md:grid-cols-1 px-20 pt-5' : 'md:grid-cols-2'} transition-all`">
        <div :class="` transition-all ${props.oneColumn ? 'col-span-1' : 'col-span-1 md:col-span-2 mt-2'}  ${props.showPostComposer ? 'h-70 p-5' : 'h-0 p-0'} overflow-hidden mx-2 bg-ll-neutral dark:bg-ld-neutral rounded-md  flex flex-col relative`">
            <textarea 
                class="w-full h-full rounded-md bg-ll-base dark:bg-ld-base p-4 outline-none text-lg"
                placeholder="What's happening?" 
                resize="none"></textarea>
            <div class="w-full flex items-center justify-between pt-3 ">
                <div class="flex">

                    <!-- Picture -->
                    <button 
                        @click="hanldePictureButtonClicked"
                        class="w-10 h-10 mr-2 border rounded-md flex justify-center items-center  border-ll-border dark:border-ld-border bg-ll-base dark:bg-ld-base dark:text-gray-500 active:scale-95 transition-transform transform">
                        <PictureSvg/>
                    </button>

                    <!-- GIF -->
                    <button 
                        @click="hanldeGifButtonClicked"
                        class="w-10 h-10 mr-2 border rounded-md flex justify-center items-center  border-ll-border dark:border-ld-border bg-ll-base dark:bg-ld-base dark:text-gray-500 active:scale-95 transition-transform transform">
                        <GifSvg/>
                    </button>
                </div>

                <div class="flex">
                    <button  
                        @click="handleClick"
                        class=" text-sm px-3 py-2 bg-ll-primary text-white dark:bg-ld-primary rounded-md flex items-center active:scale-95 transform transition-transform">
                        <!-- AddPostSvg 實際上是那個箭頭，而非整個按鈕 -->
                        <AddPostSvg/>
                        Share
                    </button>
                </div>
            </div>
            <button 
                @click="$.emit('onCloseComposePost')"
                class="w-8 h-8 absolute -top-0 -right-1 bg-ll-neutral dark:bg-ld-neutral text-sm  border-ll-border dark:border-ld-border border rounded-full flex items-center justify-center mr-2 active:scale-95 transform transition-transform">
                <CloseSvg/>
            </button>
        </div>

    </div>
</template>

<script setup lang="ts">
import AddPostSvg from './Icon/AddPostSvg.vue'
import CloseSvg from './Icon/CloseSvg.vue'
import PictureSvg from './Icon/PictureSvg.vue'
import GifSvg from './Icon/GifSvg.vue'
import {getPosts} from "../api/posts"
import {useCounterStore} from "../store/count"

const counterStore = useCounterStore();

// 定義 props 屬性
const props = defineProps({
    // 是否使用單列顯示
    oneColumn: Boolean,

    // 是否顯示發文組件，必要屬性
    showPostComposer: {
        type: Boolean,
        required: true
    }
})



const getPostsApi = (): Promise<void> => {
    return new Promise(async (resolve, reject) => {
        await getPosts()
        .then((res:any) => {
            console.log(res);
            resolve();
        }).catch((error) => {
            console.log(error);
            reject("An error has occured");
        });
    })
}

const handleClick = () => {
    let _count = counterStore.getCount();
    console.log("Add post!, count:" + _count);    
    // getPostsApi();
}

const hanldePictureButtonClicked = () =>{
    // this.$emit('onMenuClick');
    counterStore.increment();
}

const hanldeGifButtonClicked = () =>{
    // this.$emit('onMenuClick');
    counterStore.decrement();
}
</script>

<script lang="ts">
export default {

}
</script>