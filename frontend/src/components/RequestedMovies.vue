<template>
  <MovieCard :movies="store.getters['chat/movies']"></MovieCard>
</template>

<script>
import store  from '../stores';
import MovieCard from './MovieCard.vue';

export default {
data(){
  return {
    store: store,
  }
},
components: {MovieCard}
}
</script>
