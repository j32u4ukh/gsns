<template >
    <AppShell :navbarExpanded="showLeftNavbar" :rightNavbarExpanded="showRightNavbar" :full-sidebar="true"
        @on-change-theme="warn">
        <template #header>
             <!-- 頭部組件 -->
            <Header @on-menu-click="
    showLeftNavbar = !showLeftNavbar;
            " @on-right-menu-click="
    showRightNavbar = !showRightNavbar;
            ">
            </Header>
        </template>
        <template #navbar>
            <!-- 左側導航欄組件 -->
            <Navbar :is-expanded="showLeftNavbar" @on-compose-post="showComposePost = !showComposePost"
                @on-close-navbar="(v) => { showLeftNavbar = v }">
            </Navbar>
        </template>
        <template #rightNavbar>
            <!-- 右側導航欄組件 -->
            <NavbarRight>
            </NavbarRight>
        </template>
        <template #body>
            <!-- 帖子內容展示組件 -->
            <Feed :oneColumn="showLeftNavbar && showRightNavbar" :showPostComposer="showComposePost"
                @on-close-compose-post="showComposePost = !showComposePost"></Feed>

            <div class="flex flex-col p-2 ">
                <Post v-for="(post, index) in rightPosts" :post="post" :key="index"></Post>
            </div>
            <div class="flex flex-col p-2 ">
                <Post v-for="(post, index) in leftPosts" :post="post" :key="index"></Post>
            </div>
        </template>
    </AppShell>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import AppShell from '../layouts/AppShell.vue';
import Header from '../components/Header.vue';
import Navbar from '../components/Navbar.vue';
import NavbarRight from '../components/NavbarRight.vue';
import Feed from '../components/Feed.vue';
import Post from '../components/Post.vue';
import { Posts } from '../seed/processFeedPosts';
import { IPost } from '../types';

const warn = (val: boolean) => {
}

// 創建 ref 變量，用於控制左側導航欄、右側導航欄和發帖組件的展示狀態
const showLeftNavbar = ref(true);
const showRightNavbar = ref(true);
const showComposePost = ref(true);

// 定義 feedPosts 作為 IPost 陣列
const feedPosts: IPost[] = reactive(Posts);

// 建立兩個響應式的陣列
let leftPosts: IPost[] = reactive([]);
let rightPosts: IPost[] = reactive([]);

// 監視 feedPosts 變化，觸發排序函式
watch(feedPosts, (val) => {
    sortList();
});

// 在組件被掛載後執行排序函式
onMounted(sortList);

// 定義排序函式
function sortList() {
    // 遍歷 feedPosts 陣列
    feedPosts.forEach((post, index) => {
        if ((index % 2) != 0) {
            // 如果索引是奇數，將 post 放入 leftPosts
            leftPosts.push(post);
        } else {
            // 如果索引是偶數，將 post 放入 rightPosts
            rightPosts.push(post);
        }
    });
}
</script>

<script lang="ts">
export default {
    // 引入的組件
    components: { AppShell, Header, Navbar, Feed, NavbarRight }
}
</script>

<style lang="">

</style>