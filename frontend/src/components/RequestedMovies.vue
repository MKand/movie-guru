<template>
  <MovieCard class="flex justify-center items-center" :movies="store.getters['chat/movies']" :traceId="traceId" :spanId="spanId" ></MovieCard>
</template>

<script>
import ChatClientService from '@/services/ChatClientService';
import store  from '../stores';
import MovieCard from './MovieCard.vue';

export default {
data(){
  return {
    store: store,
    traceId: ChatClientService.traceId,
    spanId: ChatClientService.spanId
  }
},
components: {MovieCard}
}
</script>
