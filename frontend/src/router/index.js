import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import {store} from '../stores/index'
import LoginStatusCheckService from '@/services/LoginStatusCheckService'


const USE_FIREBASE_AUTH = import.meta.env.VITE_USE_AUTH

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
      path: "/login",
      name: "login",
      component: USE_FIREBASE_AUTH === "true" 
        ? () => import('../views/LoginFirebaseView.vue') 
        : () => import('../views/LoginSimpleView.vue')
    }
  ]
})

router.beforeEach(async (to, from) => {
  if (to.meta.requiresAuth) {
    const loggedIn = store.getters['user/loginStatus']
    if (!loggedIn && to.name !== "login") {
    const serverLoggedIn = await LoginStatusCheckService.checkLogin();
     if (serverLoggedIn){
      console.log('Valid cookie found, logging in user');
     
    } else {
      console.log('Invalid cookie, redirecting to login.');
      return { name: 'login' };
    }
     }
    }
})

export default router
