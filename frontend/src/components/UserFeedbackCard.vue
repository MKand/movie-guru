<template>
    <div class="flex flex-row overflow-y-auto flex-wrap justify-items-end scrollbar-thin scrollbar-thumb-primary scrollbar-track-accent">
      <button @click="submitFeedback('positive')" class="p-2 m-2 bg-none text-2xl rounded-md hover:bg-green-300">
        👍
      </button>
      <button @click="submitFeedback('negative')" class="p-2 m-2 bg-none text-2xl rounded-md hover:bg-red-200">
        👎
      </button>
    </div>
  </template>
  
  <script>
  import ChatClientService from '../services/ChatClientService';
  
  export default {
    props: {
      traceId: {
        type: String,
        required: true
      },
      spanId: {
        type: String,
        required: true
      },
    },
    methods: {
      submitFeedback(valueString) {
        ChatClientService.submitFeedback(this.traceId, this.spanId, valueString)
        .then(() => {})
        .catch(error => {
            console.log("Error sending chat feedback:", error);
          });
      }
    }
  };
  </script>