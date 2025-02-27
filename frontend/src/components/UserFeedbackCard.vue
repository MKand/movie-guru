<template>
    <div class="flex justify-end">
      <button 
        ref="thumbsUp" 
        @click="submitFeedback('positive', $event)" 
        class="p-2 m-2 bg-none text-2xl rounded-md hover:bg-green-300"
      >
        👍
      </button>
      <button 
        @click="submitFeedback('negative', $event)" 
        class="p-2 m-2 bg-none text-2xl rounded-md hover:bg-red-200" 
        :class="{ shake: isShaking }"
      >
        👎
      </button>
      <div 
        v-for="star in stars" 
        :key="star.id" 
        class="shooting-star" 
        :style="{ 
          top: `${star.y}px`, 
          left: `${star.x}px`, 
          width: `${star.size}px`, 
          height: `${star.size}px`, 
          backgroundColor: star.color 
        }"
      ></div>
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
      }
    },
    data() {
      return {
        stars: [],
        isShaking: false
      };
    },
    methods: {
      submitFeedback(valueString, event) {
        if (valueString === 'positive') {
          this.addShootingStars(event);
          setTimeout(() => {
        ChatClientService.submitFeedback(this.traceId, this.spanId, valueString);
      }, 1000); 
        } else {
          this.triggerShake();
          ChatClientService.submitFeedback(this.traceId, this.spanId, valueString);

        }
      },
      addShootingStars(event) {
        this.stars = [];

        const buttonRect = event.target.getBoundingClientRect();
        const buttonCenterX = buttonRect.left + buttonRect.width / 2;
        const buttonCenterY = buttonRect.top + buttonRect.height / 2;
        const colors = ['red', 'blue', 'green', 'yellow', 'purple']; // Array of colors
  
        for (let i = 0; i < 20; i++) {
          const newStar = {
            id: Date.now() + i,
            x: buttonCenterX, // Initial x position at the button center
            y: buttonCenterY, // Initial y position at the button center
            size: 20 + Math.random() * 10, // Vary size between 30px and 20px
            color: colors[Math.floor(Math.random() * colors.length)] // Random color
          };
          this.stars.push(newStar);
        }
  
        setTimeout(() => {
          this.stars = [];
        }, 1000);
      },
      triggerShake() {
        this.isShaking = true;
        setTimeout(() => {
          this.isShaking = false;
        }, 500);
      }
    }
  };
  </script>
  
  <style scoped>
  /* Shooting Stars */
  @keyframes shooting-star {
    0% {
      transform: translate(0, 0) scale(1);
      opacity: 1;
    }
    100% {
      transform: translate(50px, -50px) scale(0.5);
      opacity: 0;
    }
  }
  
  .shooting-star {
    position: absolute;
    background-color: white;
    clip-path: polygon(
      50% 0%,
      61% 35%,
      98% 35%,
      68% 57%,
      79% 91%,
      50% 70%,
      21% 91%,
      32% 57%,
      2% 35%,
      39% 35%
    );
    animation: shooting-star 1s ease-out forwards;
  }
  
  /* Sad Shake Animation */
  @keyframes shake {
    0% { transform: translateX(0); }
    25% { transform: translateX(-5px); }
    50% { transform: translateX(5px); }
    75% { transform: translateX(-5px); }
    100% { transform: translateX(0); }
  }
  
  .shake {
    animation: shake 0.5s ease-in-out;
  }
  </style>