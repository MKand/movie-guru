<template>
  <div
    class="relative flex flex-row overflow-y-auto flex-wrap justify-center items-start scrollbar-thin scrollbar-thumb-primary scrollbar-track-accent"
  >
    <div 
      v-for="m in movies"
      :key="m.title"
      class="mb-8 mx-4 my-4 w-64 md:w-72 lg:w-80 rounded-lg border-2 border-black shadow-md transition-transform duration-300 hover:scale-105 relative"
      @mouseenter="showButtons(m.title)"
      @mouseleave="hideButtons(m.title)"
    >
      <div class="bg-gray-800 rounded-t-lg py-2">
        <p class="text-center text-text font-bold mb-2 text-xl truncate">
          {{ m.title }}
        </p>
      </div>
      <img
        :src="m.poster"
        :alt="m.title"
        class="w-full h-auto rounded-b-lg shadow-lg hover:brightness-100 border-4 border-accent"
      />
      <div
        class="absolute bottom-4 left-0 right-0 flex justify-center transition-opacity duration-300"
        :class="{ 'opacity-100': buttonsToShow[m.title], 'opacity-0': !buttonsToShow[m.title] }"
      >
        <button
          class="bg-accent text-text hover:bg-primary hover:text-secondary font-semibold py-2 px-4 rounded-lg m-2 shadow-md transition-colors duration-300 hover:scale-105 flex items-center"
          @click="tellMeMore(m.title)"
        >
          <Info class="w-6 h-6" /></button>
        <button
          class="bg-green-500 text-white hover:bg-green-600 font-semibold py-2 px-4 rounded-lg m-2 shadow-md transition-colors duration-300 hover:scale-105 flex items-center"
          @click="wantToWatch(m.title,m.poster)"
        >
          <Play class="w-6 h-6" /></button>
      </div>
    </div>
  </div>
</template>

<script>
import { Play, Info } from "lucide-vue-next";
import ChatClientService from "@/services/ChatClientService";
import PlayerService from "@/services/PlayerService";

export default {
  components: {
    Play,
    Info,
  },
  props: {
    movies: {
      required: true,
    },
    traceId: {
      required: false,
    },
    spanId: {
      required: false,
    },
  },
  data() {
    return {
      buttonsToShow: {},
    };
  },
  methods: {
    tellMeMore(title) {
      const message = "Tell me more about the movie: " + title;
      ChatClientService.send(message)
        .then(() => {})
        .catch((error) => {
          console.log("Error sending chat message:", error);
        });
    },
    wantToWatch(title, poster) {
      console.log(`Want to watch: ${title}`);
      if(this.traceId && this.spanId){
        ChatClientService.submitFeatureAcceptance(this.traceId, this.spanId, 'accepted')
      }
      console.log("Current value", PlayerService.playMovie)
      PlayerService.setPoster(poster);
    },
    showButtons(title) {
      this.buttonsToShow[title] = true;
    },
    hideButtons(title) {
      this.buttonsToShow[title] = false;
    },
  },
};
</script>
