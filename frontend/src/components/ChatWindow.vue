<template>
    <div class="flex flex-col m-1 px-1 py-5 md:m-5 md:p-5 bg-primary rounded-lg shadow-[0_0_20px_10px_rgba(144,174,173,0.5)] ">
      <div class="flex justify-between items-center mb-4">
        <button @click="clearHistory" class="bg-accent text-text mx-2 my-2 py-1 px-3 rounded hover:bg-primary">Clear Chat</button>
      </div>
      <div id="chat-container" class="m-2 flex-1 bg-stars1 bg-cover bg-no-repeat bg-center p-4 rounded-lg min-h-[560px] max-h-[560px] shadow-inner overflow-anchor-none overflow-y-auto scrollbar-thin scrollbar-thumb-primary scrollbar-track-accent">
        <div
            v-for="m in store.getters['chat/messages']"
            :key="m.id"
            class="flex flex-col justify-center"
          >
            <div 
              v-if="m.sender=='user'" 
              class="shadow-lg m-1 p-2 rounded-lg bg-secondary text-primary self-end"
            >
              <div v-html="renderedMarkdown(m.message)"></div>
            </div>
            
            <div
              v-if="m.sender=='agent'" 
              class="shadow-lg p-2 m-1 rounded-lg bg-accent text-text self-start"
            >
            <img src="../assets/movie-guru.png" class="w-12 h-12 pb-2 object-contain" />

              <div v-html="renderedMarkdown(m.message)"></div>
            </div>

            <div
              v-if="m.sender=='system'" 
              class="shadow-lg p-2 m-1 rounded-lg bg-accent text-text font-bold"
            >
            <img src="../assets/reel-2.jpeg" class="w-12 h-12 pb-2 object-contain" />

              <div v-html="renderedMarkdown(m.message)" class="text-center"></div>
            </div>
          </div>
          <div
              v-if="this.processingRequest==true" 
              class="shadow-lg p-2 m-1 rounded-lg bg-accent text-text  self-start"
            >
            <img src="../assets/movie-guru.png" class="w-12 h-12 pb-2 object-contain" />

            <div id="thinkingDots" class="text-xl font-bold typing-animation w-10"></div>
          </div>
          <div
              v-if="this.errorOccured==true" 
              class="shadow-lg p-2 m-1 rounded-lg border-4 border-negative bg-accent text-text self-start"
            >
            <img src="../assets/movie-guru.png" class="w-12 h-12 pb-2 object-contain" />

            <div id="error_message" class="text-base font-bold "> Oops! Something went wrong. Try again.</div>
          </div>
      </div>
      <div class="mt-4 mx-2">
        <input
          type="text"
          v-model="this.newMessage"  
          v-on:keyup.enter="addUserMessage" 
          placeholder="Type your message and press ENTER to send..." 
          class="w-full min-h-14 py-2 px-3 rounded-lg border bg-gray-300 border-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

    </div>
  </template>
  
  <script>
  import store  from '../stores';
  import ChatClientService from '../services/ChatClientService';
  import { marked } from 'marked';
  import { ref } from 'vue';

  export default {
    data(){
      return {
        store: store,
        processingRequest: ref(false),
        errorOccured: ref(false),
        newMessage: "",
      }
    },
    created(){
      ChatClientService.getHistory().then(response => {
        store.commit('chat/clear')
        if (response.length >= 9) {
        store.commit('chat/add', {
                message: "Older messages are deleted from the system.", 
                sender: "system" 
        });    }    
          for (const message of response) { 
          store.commit('chat/add', {
              message: message.content, 
              sender: message.role 
          });
    }
    }).catch(error => {
        console.error(error);
    });
    this.scrollToBottom();
    },
    methods: {
      scrollToBottom() {
        window.setTimeout(() => {
          const chatContainer = document.getElementById('chat-container');
          if (chatContainer) {
            chatContainer.scrollTo({
              top: chatContainer.scrollHeight,
              behavior: "smooth",
            });
          }
        }, 1);
    },
      addUserMessage(){
        this.errorOccured = false;
        let message = document.querySelector('input[type="text"]').value;
        store.commit('chat/add', {"message":message, "sender":"user"})
        this.newMessage = "";
        this.processingRequest = true;
        this.scrollToBottom();
        ChatClientService.send(message).then((response) => {
          let result = response["result"]
          if(result == "SUCCESS"){
            let answer = response["answer"]
            let context = response["context"]
            if(response["preferences"]){
            store.commit('preferences/update', response["preferences"])
          }
            store.commit('chat/add',{"message":answer, "sender":"agent", "result":result});
            store.commit('chat/addMovies', context)
          }
          else if (result == "ERROR"){
            this.errorOccured = true;
          }
          else if (result == "UNSAFE"){
            store.commit('chat/add',{"message":"That was a naughty query. I cannot answer that question.", "sender":"agent", "result":result});

          }
          this.processingRequest = false;
          this.scrollToBottom();

         
        }).catch((error) => {
            console.log("Error sending chat message:", error);
            this.processingRequest = false;
            this.errorOccured = true;
          });
      },
      clearHistory(){
        this.errorOccured = false;
        ChatClientService.clearHistory().then(() => {
          store.commit('chat/clear')
        }
    ).catch(error => {
        console.error(error);
    });
      },
      renderedMarkdown(markdownText) {
      // Use marked to convert Markdown to HTML
      return marked.parse(markdownText); 
    }
    },
    mounted() {
      this.scrollToBottom();
      // Start the animation when the component is mounted
      const thinkingDots = document.getElementById("thinkingDots");
      if (thinkingDots) {
        thinkingDots.textContent = ""; // Clear initial content
        thinkingDots.classList.add("typing-animation");
      }
    }
  }
</script>


<style>
/* CSS animation for thinking dots */
@keyframes typing {
  0% {
    content: ".";
  }
  25% {
    content: "..";
  }
  50% {
    content: "...";
  }
  75% {
    content: "....";
  }
  100% {
    content: ".";
  }
}

.typing-animation::after {
  content: "";
  animation: typing 1s steps(4, end) infinite;
}
</style>