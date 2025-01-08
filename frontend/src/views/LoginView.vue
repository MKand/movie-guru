<script>
import { ref } from 'vue';
import {
  signInWithPopup,
  GoogleAuthProvider,
  getAuth
} from 'firebase/auth'
import store from '../stores';
import router from '../router'
import LoginClientService from '../services/LoginClientService';

let loginFailed = ref(false);

export default {
  data() {
    return {
      loginFailed: ref(false),
    }
  },

  methods: {
    handleSignIn() {
      const provider = new GoogleAuthProvider();
      signInWithPopup(getAuth(), provider)
        .then((result) => {
          let inviteCode = document.querySelector('input[type="text"]').value;

          LoginClientService.login(result.user, inviteCode).then(() => {
            store.commit('user/logIn', result.user)
            router.push('/')
            this.loginFailed = false;
          }).catch(() => {
            this.loginFailed = true;
          })
        })
    }
  }
}


</script>
<template>
  <div class="w-full h-screen flex flex-col justify-center items-center align-middle">
    <div
      class="flex flex-col md:h-1/2 md:w-1/2 lg:h-1/3 lg:w-1/3 xl:h-1/4 xl:w-1/4 w-full justify-center items-center align-middle bg-stars1 bg-cover bg-no-repeat bg-center p-20">
      <input type="text" v-if="loginFailed == true"
        class="text-bold placeholder-accent  p-2 m-2 w-full rounded-lg   bg-gray-300 text-primary border-2 border-negative text-center"
        placeholder="Something went wrong... Enter your email and try again.">
      <input type="text" v-else
        class="text-bold   p-2 m-2 w-full rounded-lg   bg-gray-300 placeholder-accent text-primary text-center"
        placeholder="If this your first time on Movie Guru, enter your Invite Code here...">
      <button type="button" @click="handleSignIn"
        class="text-white w-1/2  bg-[#4285F4] hover:bg-[#4285F4]/90 focus:ring-4 focus:outline-none focus:ring-[#4285F4]/50 font-medium rounded-lg text-sm m-5 px-5 py-2 text-center align-middle inline-flex items-center justify-between mr-2 mb-2">Sign
        in<div></div></button>
    </div>
  </div>
</template>
