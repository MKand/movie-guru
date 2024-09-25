import './css/main.css'
import {store} from './stores/index'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

createApp(App).use(router).use(store).mount('#app')
