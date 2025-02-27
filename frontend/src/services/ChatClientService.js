import {fetch} from 'whatwg-fetch'
import store  from '../stores';

class ChatClientService {

  handleAddedUserMessage = null;
  handleAgentMessage = null;
  handleErrorMessage = null;
  handleSubmittedFeedback = null;

  async send(message){
   this.handleAddedUserMessage();
    
    store.commit('chat/add', {"message":message, "sender":"user"})
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json'},
        body: JSON.stringify({ content: message }),
        credentials: 'include'
    };
    try{
      const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/chat', requestOptions)
    
      if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
      }
        const json = await response.json();
        const result = json["result"];
        if(result == "SUCCESS"){
          store.commit('chat/add',{"message":json["answer"], "sender":"agent", "result":result});
          store.commit('chat/addMovies', json["context"])
          this.handleAgentMessage(json["traceId"], json["spanId"])
        }
        else if (result == "ERROR" || result == "QUOTALIMIT" || result == "UNSAFE"){
          this.handleErrorMessage(json["answer"])
        }
        if(json["preferences"]){
          store.commit('preferences/update', json["preferences"])
        }
        return json
      }   catch (error) {
        this.handleErrorMessage("I've had trouble connecting to the server. Try again.")
        throw error;
      }
    }
 
    
  async startup(){
    const requestOptions = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json'},
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
      if (result == "SUCCESS"){
        store.commit('chat/addPlaceHolderMovies', context)
        store.commit('preferences/update', preferences)
        }
    } catch (error) {
      console.error(error.message);
      throw error;
    }
    
    async getHistory(){
      const requestOptions = {
          method: 'GET',
          headers: { 'Content-Type': 'application/json'},
          credentials: 'include'
      };
      const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/history', requestOptions)
      
      if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
      }
        const json = await response.json();
        return json
      } catch (error) {
        console.error(error.message);
        throw error;
      }

      async submitFeedback(traceId, spanId, valueString){
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify({traceId: traceId, spanId: spanId,  name: "chatFlow", feedbackExperience: valueString}),
        };
        const response = await fetch(import.meta.env.VITE_CHAT_SERVER_URL + '/feedback', requestOptions)

        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
          // this.handleSubmittedFeedback();
          const json = await response.json();
          return json
        } catch (error) {
          console.error(error.message);
          throw error;
        }
  

      async clearHistory(){
      const requestOptions = {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json'},
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

export default new ChatClientService();
