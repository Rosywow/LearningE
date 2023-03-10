import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecord } from 'vue-router'

// load modules routes
// if need children routes in modules, please DO NOT name it *.route.ts
const modules: any = import.meta.glob('@/modules/**/*.route.ts', {
  eager: true,
});
let modulesRoutes: RouteRecord[] = [];

for (const path in modules) {
  const routes = modules[path].default.getRoutes();
  if (routes.length > 0) {
    modulesRoutes = modulesRoutes.concat(routes);
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ()=>import('../views/AboutView.vue'),
    },
    ...modulesRoutes,
  ],
});

export default router;
