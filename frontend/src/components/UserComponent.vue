<template>
      <div class="flex flex-col justify-center items-center md:m-5 rounded-lg bg-gradient-to-b  from-start via-accent to-secondary  scrollbar-thin scrollbar-thumb-primary scrollbar-track-accent">
        <h1 class="text-gurusilver text-lg p-2 font-bold mt-2 m-1"> {{ store.state.user.email }} </h1>
                <h2 class="text-text text-lg pt-2 m-1 text-center"> The Movie Guru tries to learn your movie preferences.    </h2>
                <h2 class="text-text text-lg pb-2 m-1 text-center">You can instruct it to take your likes and dislikes into account.</h2> 
        <div class="flex flex-row flex-wrap justify-center overflow-y-auto m-5 scrollbar-thin scrollbar-thumb-primary scrollbar-track-accent">

            <div v-for="a in store.getters['preferences/preferences']['likes']['genres']"  class="rounded-full bg-secondary w-auto align-middle m-2 p-2 text-primary shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('likes', 'genres', a)"> ✖ </button>
            </div> 
            <div v-for="a in store.getters['preferences/preferences']['likes']['actors']"  class="rounded-full bg-secondary w-auto align-middle m-2 p-2 text-primary  shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('likes', 'actors', a)"> ✖ </button>
            </div> 
            <div v-for="a in store.getters['preferences/preferences']['likes']['director']"  class="rounded-full bg-secondary w-auto align-middle m-2 p-2 text-primary  shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('likes', 'director', a)"> ✖ </button>
            </div> 
            <div v-for="a in store.getters['preferences/preferences']['likes']['other']"  class="rounded-full bg-secondary w-auto align-middle m-2 p-2 text-primary  shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('likes', 'other', a)"> ✖ </button>
            </div> 

             <div v-for="a in store.getters['preferences/preferences']['dislikes']['genres']"  class="rounded-full bg-negative w-auto align-middle m-2 p-2 text-text  shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('dislikes', 'genres', a)"> ✖ </button>
            </div> 
            <div v-for="a in store.getters['preferences/preferences']['dislikes']['actors']"  class="rounded-full bg-negative w-auto align-middle m-2 p-2 text-text  shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('dislikes', 'actors', a)"> ✖ </button>
            </div> 
            <div v-for="a in store.getters['preferences/preferences']['dislikes']['director']"  class="rounded-full bg-negative w-auto align-middle m-2 p-2 text-text  shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('dislikes', 'director', a)"> ✖ </button>
            </div> 
            <div v-for="a in store.getters['preferences/preferences']['dislikes']['other']"  class="rounded-full bg-negative w-auto align-middle m-2 p-2 text-text  shadow-md shadow-black">{{ a }}
                <button class="text-text hover:text-pop text-lg" @click="deletePref('dislikes', 'other', a)"> ✖ </button>
            </div> 

          </div>
            <div v-if="store.getters['user/loginStatus']==true">
            <button class="bg-accent text-text py-2 px-5 rounded hover:bg-primary m-5 text-lg" @click="handleSignOut">Sign out</button>      
        </div>
      </div>
  </template>
  
<script>
  import store  from '../stores';
  import LoginClientService from '../services/LoginClientService';
  import { getAuth } from "firebase/auth";
  import PreferencesClientService from '../services/PreferencesClientService';
  import {ref } from 'vue';


  export default {
    data(){
      return {
        store: store,
        enableAdd: ref(false)

      }
    },

  methods: {
    userAddClick(){
        this.enableAdd = !this.enableAdd;
      },
     handleSignOut() { 
      getAuth().signOut().then(() => {
            LoginClientService.logout().then(() =>{
              store.commit('user/logOut')
              console.log("signed out")
              window.location.reload()
              console.log("logout successful")
            }).catch((reason) =>{
              console.error('Failed logout', reason)
            })
        })
      },
      deletePref(type, key, value){
        this.enableAdd = false;
        this.store.commit('preferences/delete', {"type":type, "key":key, "value":value})
        PreferencesClientService.update(this.store.getters['preferences/preferences'])
        this.enableAdd = true;
      },

      addPref(type, key, divId){
        if (!this.enableAdd) return;
        this.enableAdd = false;

        // Get the div element by its ID
        const divElement = document.getElementById(divId);

        // Create a new div element with the desired class
        const newDiv = document.createElement("div");
        newDiv.classList.add("rounded-full", "w-auto", "align-middle");

        // Create an input element for text input
        const input = document.createElement("input");
        input.type = "text";
        input.placeholder = "Enter preference or ESC";
        input.classList.add("rounded-full", "bg-primary", "w-auto", "align-middle", "m-2", "p-2", "text-text");

        // Add an event listener to the input to add the preference when Enter is pressed
        input.addEventListener("keyup", (event) => {
          if (event.key === "Enter") {
            // Add the input value to the preferences array
            this.store.commit('preferences/add', { type, key, value: input.value });
            PreferencesClientService.update(this.store.getters['preferences/preferences']);

            // Create a new div element to display the added preference
            const newPreferenceDiv = document.createElement("div");
            newPreferenceDiv.classList.add("rounded-full", "bg-secondary", "w-auto", "align-middle", "m-2", "p-2", "text-primary");
            newPreferenceDiv.textContent = input.value;
            this.enableAdd = true;
            divElement.removeChild(newDiv);
          }
          if (event.key === "Escape") {
            divElement.removeChild(newDiv);
            this.enableAdd = true;
          }
        });

        // Append the input element to the new div
        newDiv.appendChild(input);

        // Append the new div to the divElement
        divElement.appendChild(newDiv);

        input.focus();
      }
  }
  }
</script>