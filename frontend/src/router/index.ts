import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('@/pages/home/Home.vue')
        },
        {
            path: '/summary',
            name: 'summary',
            component: () => import('@/pages/summary/SummaryPage.vue')
        },
        {
            path: '/detail/:server/:protocol',
            name: 'detail',
            component: () => import('@/pages/detail/DetailPage.vue')
        },
        {
            path: '/traceroute',
            name: 'traceroute',
            component: () => import('@/pages/traceroute/TraceroutePage.vue')
        }
    ]
})

router.beforeEach((to, from, next) => {
    if (to.path === '/') {
        next({ path: '/summary' })
    } else {
        next()
    }
})

export default router