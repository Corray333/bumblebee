import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/products',
      alias: "/",
      name: 'products',
      component: () => import('../views/ProductsView.vue'),
    },
  ],
})

export default router
