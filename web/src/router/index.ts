import { createRouter, createWebHistory } from "vue-router";
import Main from '../pages/Main.vue'
import ThePost from '../pages/ThePost.vue'

const routes = [
  {
    path: "/",
    name: "Main",
    component: Main,
  },
  {
    path: "/post/:postId",
    name: 'ThePost',
    component: ThePost,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { left: 0, top: 0 };
  },
});

export default router;
