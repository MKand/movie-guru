<template>
  <div id="app" class="flex flex-col h-full  bg-gradient-to-b  from-start via-accent to-start scrollbar scrollbar-thumb-primary scrollbar-track-accent overflow-y-auto">
    <Header />
    <router-view />
    <Footer class="sticky bottom-0"/>
  </div>
</template>

<script>
import Header from './components/Header.vue';
import Footer from './components/Footer.vue';

export default {
  components: {
    Header,
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
