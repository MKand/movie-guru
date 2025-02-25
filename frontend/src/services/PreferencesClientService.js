import {fetch as fetchPolyfill} from 'whatwg-fetch'
import store  from '../stores';

class PreferencesClientService {
  async get(){
    const requestOptions = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json'},
        credentials: 'include'
    };
    const response = await fetchPolyfill(import.meta.env.VITE_CHAT_SERVER_URL + '/preferences', requestOptions)
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
  
      const json = await response.json();
      return json
    } catch (error) {
      console.error(error.message);
      throw error;
    }
  
    async update(){
      const requestOptions = {
          method: 'POST',
          headers: { 'Content-Type': 'application/json'},
          body: JSON.stringify({ content: store.getters['preferences/preferences'] }),
          credentials: 'include'
      };
      const response = await fetchPolyfill(import.meta.env.VITE_CHAT_SERVER_URL + '/preferences', requestOptions)
      if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
        }
        const json = await response.json();
        return json
      } catch (error) {
        console.error(error.message);
        throw error;
      }      
}

export default new PreferencesClientService();
