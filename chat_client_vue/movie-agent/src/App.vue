<template>
  <div id="app" class="flex flex-col h-full  bg-gradient-to-b  from-start via-accent to-start scrollbar scrollbar-thumb-primary scrollbar-track-accent overflow-y-auto">
    
    <!-- <div id="app" class="flex flex-col h-screen w-screen items-center justify-center min-h-screen bg-stars1 bg-cover bg-no-repeat bg-center">  -->
      <Header />
    <router-view />
    <Footer class="sticky bottom-0"/>
  </div>
</template>

<script>
import Header from './components/Header.vue';
import MoviePosters from './components/MoviePosters.vue';
import ChatWindow from './components/ChatWindow.vue';
import Footer from './components/Footer.vue';

export default {
  components: {
    Header,
    MoviePosters,
    ChatWindow,
    Footer,
  },
};
</script>

<style>
/* Hide scrollbar arrows */
::-webkit-scrollbar-button {
  display: none;
}
</style>
