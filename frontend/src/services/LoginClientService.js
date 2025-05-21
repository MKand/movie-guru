import {fetch as fetchPolyfill} from 'whatwg-fetch'

const USE_FIREBASE_AUTH = import.meta.env.VITE_USE_AUTH

class LoginClientService {
  async login(userId, inviteCode="", accessToken="") {
    if(userId == ""){
      throw new Error("User cannot be empty");
    }
    try {
      const headers = {
        'Content-Type': 'application/json',
        'User': userId,
      }
      if (USE_FIREBASE_AUTH=="true"){
        headers['Authorization'] = `Bearer ${accessToken}`
      }
      const requestOptions = {
        method: 'POST',
        headers: headers,
        body: JSON.stringify({ inviteCode }),
        credentials: 'include',
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
          headers: { 'Content-Type': 'application/json'},
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
