import {createRouter, createWebHistory, RouteRecordRaw} from "vue-router";

const routes: Readonly<RouteRecordRaw[]> = [
    {
        path: '/',
        redirect: '/self-update',
    },
    {
        path: '/self-update',
        name: 'SelfUpdate',
        component: () => import('../components/SelfUpdate.vue'),
    },
    {
        path: '/library',
        name: 'Library',
        component: () => import('../components/Library.vue'),
    },
    {
        path: '/apps/:id',
        name: 'App',
        component: () => import('../components/App.vue'),
        redirect: '/apps/:id/overview',
        children: [
            {
                path: 'overview',
                name: 'AppOverview',
                component: () => import('../components/AppOverview.vue'),
            },
            {
                path: 'sdk',
                name: 'AppSDK',
                component: () => import('../components/AppSDK.vue'),
            }
        ]
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('../components/NotFound.vue')
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router