import { createRouter, createWebHistory } from "vue-router";
import Main from '../pages/Main.vue'

const routes = [
  {
    path: "/",
    name: "main",
    component: Main,
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
