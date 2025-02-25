<template>
  <h1 class="font-guru-title text-xl font-bold uppercase tracking-wide text-[#C6D4D2] 
drop-shadow-2xl text-center mb-5">
    Featured films
  </h1>
  <MovieCard :movies="store.getters['chat/placeHolderMovies']"></MovieCard>
</template>

<script>
import store from '../stores';
import MovieCard from './MovieCard.vue';
import ChatClientService from '../services/ChatClientService';
export default {
  data() {
    return {
      store: store,
    }
  },
  components: { MovieCard },
  created() {
    ChatClientService.startup()
      .catch(error => {
        console.error(error);
      })
  }
}
</script>