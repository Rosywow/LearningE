import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/reading/addsource', name: 'addsource', component: () => import('../views/addsource.vue') },
    { path: '/reading/articleList', name: 'articleList', component: () => import('../views/articleList.vue') },
    { path: '/reading/article:id', name: 'article', component: () => import('../views/article.vue') },
    { path: '/reading/video', name: 'video', component: () => import('../views/video.vue') },
  ],
});

export default router;