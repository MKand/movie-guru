<script >
import RequestedMovies from '@/components/RequestedMovies.vue';
import ChatWindow from '@/components/ChatWindow.vue';
import UserComponent from '@/components/UserComponent.vue';
import FeaturedMovies from '@/components/FeaturedMovies.vue';
import PlayerWindow from '@/components/PlayerWindow.vue';
import store from '@/stores';
import PlayerService from '@/services/PlayerService';

export default {
  
  components: { RequestedMovies, PlayerWindow, ChatWindow, FeaturedMovies, UserComponent },
  data() {
    return {
      playMovie: PlayerService.playMovie,
      poster: PlayerService.poster,
      store,
    };
  },
  methods: {
  
  },
  created() {
    PlayerService.unsetPoster();
  }
};

</script>

<template>
    <div class="md:p-5 my-5 py-5 mx-3 md:m-5 flex flex-col-reverse md:flex-row  justify-items-center lg:justify-items-center overflow-y-auto scrollbar scrollbar-thumb-primary scrollbar-track-accent md:h-[860px]">
      <div v-if="store.getters['chat/movies'].length>0" class=" rounded-lg md:w-1/2 xxl:w-2/5 justify-evenly items-center align-middle scrollbar-thin scrollbar-thumb-primary scrollbar-track-accent overflow-y-auto">
              <RequestedMovies  class="overflow-y-clip" />
      </div>
 
        <ChatWindow  class="md:w-1/2 xxl:w-2/5 w-full"/>
    </div>
    <FeaturedMovies />
    <UserComponent class="" />
</template>
