import './css/main.css'
import {store} from './stores/index'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { initializeApp } from "firebase/app";
import { VueFire, VueFireAuth } from 'vuefire'

const USE_FIREBASE_AUTH = import.meta.env.VITE_USE_AUTH

if(USE_FIREBASE_AUTH=="true"){
const firebaseConfig = {
  projectId: import.meta.env.VITE_GCP_PROJECT_ID,
  appId: import.meta.env.VITE_FIREBASE_APP_ID,
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
};

// Initialize Firebase
const firebaseApp = initializeApp(firebaseConfig);


createApp(App)
.use(router)
.use(store)
.use(VueFire, {
    firebaseApp,
    modules:[
        VueFireAuth(),
    ]
})
.mount('#app')
}
else{
    createApp(App)
    .use(router)
    .use(store)
    .mount('#app')
}

