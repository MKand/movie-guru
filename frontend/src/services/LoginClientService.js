import {fetch as fetchPolyfill} from 'whatwg-fetch'

class LoginClientService {
  async login(user, inviteCode) {
    if(user == ""){
      throw new Error("User cannot be empty");
    }
    try {
      const requestOptions = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'User': user,
        },
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
