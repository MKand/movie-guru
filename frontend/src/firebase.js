
import { initializeApp } from "firebase/app";

const USE_AUTH = import.meta.env.VITE_USE_AUTH

if(USE_AUTH=="true"){
  const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_GCP_PROJECT_ID,
    appId: import.meta.env.VITE_FIREBASE_APP_ID,
  };

// Initialize Firebase
   const firebaseApp = initializeApp(firebaseConfig);
}



