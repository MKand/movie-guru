import { fetch } from 'whatwg-fetch'
import store from '../stores';
import { ref } from 'vue';

class ChatClientService {

  handleAddedUserMessage = null;
  handleAgentMessage = null;
  handleErrorMessage = null;
  handleSubmittedFeedback = null;
  traceId = ref(null);
  spanId = ref(null);
  featureAccepted = ref(false)

  clearTraceIds() {
    this.traceId.value = null;
    this.spanId.value = null;
  }

  setTraceIds(traceId, spanId) {
    this.traceId.value = traceId;
    this.spanId.value = spanId;
  }

  async send(message) {
    var requestOptions = {}
    try{
      this.handleAddedUserMessage();
      this.clearTraceIds();
      requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: message }),
        credentials: 'include'
      };
      store.commit('chat/add', { "message": message, "sender": "user" })

    }
    catch (error) {
      this.handleErrorMessage("I've had trouble sending your message. Try again.")
    }

    try {
      const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/chat', requestOptions)

      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const json = await response.json();
      const result = json["result"];
      if (result == "SUCCESS") {
        store.commit('chat/add', { "message": json["answer"], "sender": "agent", "result": result });
        store.commit('chat/addMovies', json["context"])
        this.setTraceIds(json["traceId"], json["spanId"])
        this.handleAgentMessage();
      }
      else if (result == "ERROR" || result == "QUOTALIMIT" || result == "UNSAFE" || result == "BAD_QUERY" || result == "TOO_LONG") {
        this.handleErrorMessage(json["answer"]|| "unknown error occurred")
      }
      if (json["preferences"]) {
        store.commit('preferences/update', json["preferences"])
      }

      return json
    } catch (error) {
      this.handleErrorMessage("I've had trouble connecting to the server. Try again.")
    }
  }


  async startup() {
    const requestOptions = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include'
    };
    const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/startup', requestOptions)

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }
    const json = await response.json();
    let context = json["context"]
    let result = json["result"]
    let preferences = json["preferences"]
    if (result == "SUCCESS") {
      store.commit('chat/addPlaceHolderMovies', context)
      store.commit('preferences/update', preferences)
    }
  } 

  async getHistory() {
    const requestOptions = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include'
    };
    const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/history', requestOptions)

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }
    const json = await response.json();
    return json
  } 

  async submitFeedback(traceId, spanId, valueString) {
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ traceId: traceId, spanId: spanId, name: "chatFlow", feedbackExperience: valueString }),

    };
    const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/feedback', requestOptions)

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }
    // this.handleSubmittedFeedback();
    const json = await response.json();
    return json
  }

  async submitFeatureAcceptance(traceId, spanId, accepted) {
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ traceId: traceId, spanId: spanId, name: "chatFlow", accepted: accepted }),

    };
    const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/acceptance', requestOptions)

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }
    // this.handleSubmittedFeedback();
    const json = await response.json();
    return json
  }


  async clearHistory() {
    try {
      this.clearTraceIds();
      const requestOptions = {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include'
      };
      const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/history', requestOptions)

      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      return;
    } catch (error) {
      console.error(error.message);
      throw error;
    }

  }
}
export default new ChatClientService();
