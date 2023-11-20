import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect:"/project"
    },
    {
      path: '/project',
      name: 'project',
      component: () => import('../views/Project.vue')
    },
    {
      path: '/pack',
      name: 'pack',
      component: () => import('../views/Pack.vue')
    },
    {
      path: '/suite',
      name: 'suite',
      component: () => import('../views/Suite.vue')
    },
    {
      path: '/database',
      name: 'database',
      component: () => import('../views/Database.vue')
    },
    {
      path: '/result',
      name: 'result',
      component: () => import('../views/Result.vue')
    },
    {
      path: '/task',
      name: 'task',
      component: () => import('../views/Task.vue')
    },
    {
      path: '/setting',
      name: 'setting',
      component: () => import('../views/Setting.vue')
    },
  ]
})

export default router
