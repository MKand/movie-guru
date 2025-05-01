import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import {store} from '../stores/index'
import LoginStatusCheckService from '@/services/LoginStatusCheckService'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
      path: "/login",
      name: "login",
      component: () => import('../views/LoginView.vue'),
    }
  ]
})

router.beforeEach(async (to, from) => {
  if (to.meta.requiresAuth) {
    const loggedIn = store.getters['user/loginStatus']
    if (!loggedIn && to.name !== "login") {
     const user = await LoginStatusCheckService.checkLogin();
     if (user){
      console.log('Valid cookie found, logging in user:', user);
      store.commit('user/logIn', user);
    } else {
      console.log('Invalid cookie, redirecting to login.');
      return { name: 'login' };
    }
     }
    }
})

export default router
