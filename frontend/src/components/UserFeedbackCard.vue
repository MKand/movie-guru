<template>
    <div class="flex justify-start"> 
      <div class="flex bg-accent rounded-lg px-2"> 
        <h1 class="text-text font-bold align-middle leading-none pt-5">Was my response useful?</h1>
        <button
          ref="thumbsUp"
          @click="submitFeedback('positive', $event)"
          class="p-2 m-1 bg-none text-xl rounded-md hover:bg-green-300 transform transition-transform duration-300"
          :class="{ 'scale-125': isExpanding }"
        >
          👍
        </button>
        <button
          @click="submitFeedback('negative', $event)"
          class="p-2 m-1 bg-none text-xl rounded-md hover:bg-red-200"
          :class="{ shake: isShaking }"
        >
          👎
        </button>
      </div>
    </div>
  </template>
  
  <script>
  import ChatClientService from '../services/ChatClientService';
  
  export default {
    props: {
      traceId: {
        type: String,
        required: true,
      },
      spanId: {
        type: String,
        required: true,
      },
    },
    data() {
      return {
        isShaking: false,
        isExpanding: false,
      };
    },
    methods: {
      submitFeedback(valueString, event) {
        if (valueString === 'positive') {
          this.triggerExpand();
          setTimeout(() => {
            ChatClientService.submitFeedback(this.traceId, this.spanId, valueString);
          }, 600);
        } else {
          this.triggerShake();
          ChatClientService.submitFeedback(this.traceId, this.spanId, valueString);
        }
      },
      triggerShake() {
        this.isShaking = true;
        setTimeout(() => {
          this.isShaking = false;
        }, 500);
      },
      triggerExpand() {
        this.isExpanding = true;
        setTimeout(() => {
          this.isExpanding = false;
        }, 300);
      },
    },
  };
  </script>
  
  <style scoped>
  /* Sad Shake Animation */
  @keyframes shake {
    0% {
      transform: translateX(0);
    }
    25% {
      transform: translateX(-5px);
    }
    50% {
      transform: translateX(5px);
    }
    75% {
      transform: translateX(-5px);
    }
    100% {
      transform: translateX(0);
    }
  }
  
  .shake {
    animation: shake 0.5s ease-in-out;
  }
  </style>