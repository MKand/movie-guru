import './css/main.css'
import {store} from './stores/index'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import {firebaseApp} from './firebase'
import { VueFire, VueFireAuth } from 'vuefire'



createApp(App).use(router).use(store).use(VueFire, {
    firebaseApp,
    modules: [
      VueFireAuth(),
    ],
  }).mount('#app')
