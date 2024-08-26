<script>
import { ref } from 'vue';
import {
  signInWithPopup,
  GoogleAuthProvider, 
  getAuth
} from 'firebase/auth'
import store  from '../stores';
import router from '../router'
import LoginClientService from '../services/LoginClientService';

let loginFailed = ref(false);

export default {
  data(){
      return {
        loginFailed: ref(false),
      }
},

methods: {
 handleGoogleSignIn() {
const provider = new GoogleAuthProvider();
  signInWithPopup(getAuth(), provider)
  .then((result) => {
    let inviteCode = document.querySelector('input[type="text"]').value;

    LoginClientService.login(result.user, inviteCode).then(() =>{
      store.commit('user/logIn', result.user)
      router.push('/')
      this.loginFailed = false;
    }).catch(() =>{
      this.loginFailed = true;
    })
  })
 }
 
}
}

</script>

<template>
    <div class="w-full h-screen flex flex-col justify-center items-center align-middle">
        <div class="flex flex-col md:h-1/2 md:w-1/2 lg:h-1/3 lg:w-1/3 xl:h-1/4 xl:w-1/4 w-full justify-center items-center align-middle bg-stars1 bg-cover bg-no-repeat bg-center p-20">  
          <input type="text" v-if="loginFailed==true" class="text-bold placeholder-accent  p-2 m-2 w-full rounded-lg   bg-gray-300 text-primary border-2 border-negative text-center" placeholder="Wrong Invite Code... Enter the correct code and try again." >  
          <input type="text" v-else class="text-bold   p-2 m-2 w-full rounded-lg   bg-gray-300 placeholder-accent text-primary text-center" placeholder="If this your first time on Movie Guru, enter your Invite Code here..." >            
          <button type="button" @click="handleGoogleSignIn" class="text-white w-1/2  bg-[#4285F4] hover:bg-[#4285F4]/90 focus:ring-4 focus:outline-none focus:ring-[#4285F4]/50 font-medium rounded-lg text-sm m-5 px-5 py-2.5 text-center inline-flex items-center justify-between mr-2 mb-2"><svg class="mr-2 -ml-1 w-4 h-4" aria-hidden="true" focusable="false" data-prefix="fab" data-icon="google" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 488 512"><path fill="currentColor" d="M488 261.8C488 403.3 391.1 504 248 504 110.8 504 0 393.2 0 256S110.8 8 248 8c66.8 0 123 24.5 166.3 64.9l-67.5 64.9C258.5 52.6 94.3 116.6 94.3 256c0 86.5 69.1 156.6 153.7 156.6 98.2 0 135-70.4 140.8-106.9H248v-85.3h236.1c2.3 12.7 3.9 24.9 3.9 41.4z"></path></svg>Sign in with Google<div></div></button>
        </div>         
    </div>      
 </template>
