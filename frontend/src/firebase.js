
import { initializeApp } from "firebase/app";

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_GCP_PROJECT_ID,
  appId: import.meta.env.VITE_FIREBASE_APPID,
};


// Initialize Firebase
export const firebaseApp = initializeApp(firebaseConfig);