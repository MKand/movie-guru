import {fetch as fetchPolyfill} from 'whatwg-fetch'
import store  from '../stores';

class LoginClientService {
  async login(user, inviteCode) {
    try {
      const requestOptions = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${user.accessToken}`,
        },
        body: JSON.stringify({ inviteCode }),
        credentials: 'include', // Include cookies or authentication credentials
      };
  
      const response = await fetch(
        `${import.meta.env.VITE_CHAT_SERVER_URL}/login`,
        requestOptions
      );
  
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
  
      const json = await response.json();
      return json;
    } catch (error) {
      console.error(error.message);
      throw error;
    }
  }
  
    async logout(){
      const requestOptions = {
          method: 'GET',
          headers: { 'Content-Type': 'application/json', 'user': store.getters["user/email"]},
          credentials: 'include'
        };
      const response = await fetchPolyfill(import.meta.env.VITE_CHAT_SERVER_URL + '/logout', requestOptions)
      if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
        }
        return
      } catch (error) {
        console.error(error.message);
        throw error;
      }
}

export default new LoginClientService();
