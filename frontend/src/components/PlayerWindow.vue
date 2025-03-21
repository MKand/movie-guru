<template>
    <div
      ref="player"
      class="p-5 relative flex flex-row overflow-y-auto flex-wrap justify-center items-start scrollbar-thin scrollbar-thumb-primary scrollbar-track-accent"
    >
      <img
        :src="poster"
        alt="Movie Poster"
        class="w-full h-full object-cover"
      />
  
      <div class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
        <button
          class="bg-white rounded-full shadow-lg hover:bg-gray-300 transition-transform duration-200 focus:outline-none"
          @click="playMovie"
        >
        <LucidePlay class=" w-20 h-20 bg-green-600 text-text" />
        </button>
      </div>
  
      <button
        class="absolute top-2 right-2 bg-white rounded-full shadow-md hover:bg-gray-300 transition-opacity duration-200 focus:outline-none"
        @click="toggleFullScreen"
      >
        <lucide-icon icon="expand" size="24" />
      </button>
      <button
        class="absolute top-2 left-2 bg-white rounded-full shadow-md hover:bg-gray-300 transition-opacity duration-200 focus:outline-none"
        @click="closePlayer"
      >
      <LucideX class="w-10 h-10 text-text bg-primary rounded-full hover:bg-text hover:text-primary" />
      </button>
    </div>
  </template>
  
  <script>
  import PlayerService from '@/services/PlayerService';
  import { Play as LucidePlay, X as LucideX } from 'lucide-vue-next';
  
  export default {
    components: {
      LucidePlay, LucideX // Register Lucide
    },
    props: {
      poster: {
        type: String,
        required: true,
      },
    },
    methods: {
      playMovie() {
        alert('Play button clicked! Implement video playback logic here.');
      },
      toggleFullScreen() {},
      closePlayer() {
        PlayerService.unsetPoster();
      },
    },
  };
  </script>