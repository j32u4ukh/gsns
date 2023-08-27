<template >
    <div v-if="posts[0]">
        <Post :post="posts[0]"></Post>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import {getThePost} from "../api/post"
import { useRoute } from 'vue-router';
import * as ptc from "../protocol/post";
import Post from '../components/Post.vue';
import { IPost } from '../types';

const postId = ref();
const route = useRoute();
const posts = ref<IPost[]>([]);

onMounted(async () => { 
    // 在這裡可以使用 $route
    postId.value = route.params.postId;
    console.log("post id: " + route.params.postId);
    await getPostsApi();
    console.log("Layout created");
});


const getPostsApi = (): Promise<void> => {
    return new Promise(async (resolve, reject) => {
        await getThePost(postId.value)
        .then((res: any) => {
            let response: ptc.GetThePostResponse =  <ptc.GetThePostResponse>res; 
            console.log("response:");
            console.log(response);

            // Loop through the posts array and log each post's information
            for (const post of response.pms) {
                let ipost: IPost = {
                        id: post.id,
                        author: {
                            id: post.user_id,
                            name: "Henry",
                        },
                        title: "",
                        description: post.content,
                        pictures: [],
                        categories: [""],
                        likes: [],
                        shares: [],
                        replys: [],
                    };
                posts.value = posts.value.concat(
                    ipost
                );
              console.log(`Post ID: ${post.id}, Title: ${post.content}`);
            }
            resolve();
        }).catch((error: ptc.GetThePostResponse) => {
            console.log(error);
            reject("An error has occured");
        });
    })
}

// onMounted(() => {
//   // 在這裡可以使用 $route
//   postId.value = route.params.postId;
//   console.log("post id: " + route.params.postId);
//   getPostsApi();
//   console.log("Layout created");
// });

</script>

<script lang="ts">
export default {

}
</script>

<style lang="">

</style>