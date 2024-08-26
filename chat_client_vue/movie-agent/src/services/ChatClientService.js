import {fetch as fetchPolyfill} from 'whatwg-fetch'
import store  from '../stores';

class ChatClientService {
  async send(message){
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'user': store.getters["user/email"]},
        body: JSON.stringify({ content: message }),
        credentials: 'include'
    };
    const response = await fetchPolyfill(import.meta.env.VITE_CHAT_SERVER_URL + '/chat', requestOptions)
    
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
      const json = await response.json();
      return json
    } catch (error) {
      console.error(error.message);
      throw error;
    }
    
  async startup(){
    const requestOptions = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json', 'user': store.getters["user/email"]},
        credentials: 'include'
    };
    const response = await fetchPolyfill(import.meta.env.VITE_CHAT_SERVER_URL + '/startup', requestOptions)
    
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
      const json = await response.json();
      return json
    } catch (error) {
      console.error(error.message);
      throw error;
    }
    
    async getHistory(){
      const requestOptions = {
          method: 'GET',
          headers: { 'Content-Type': 'application/json', 'user': store.getters["user/email"]},
          credentials: 'include'
      };
      const response = await fetchPolyfill(import.meta.env.VITE_CHAT_SERVER_URL + '/history', requestOptions)
      
      if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
      }
        const json = await response.json();
        return json
      } catch (error) {
        console.error(error.message);
        throw error;
      }

      async clearHistory(){
      const requestOptions = {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json', 'user': store.getters["user/email"]},
          credentials: 'include'
      };
      const response = await fetchPolyfill(import.meta.env.VITE_CHAT_SERVER_URL + '/history', requestOptions)
      
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
