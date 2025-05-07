import store  from '../stores';

class LoginStatusCheckService {

async checkLogin() {
    const loggedIn = store.getters['user/loginStatus'];
    if (!loggedIn) {
        try {
          const user = await this.checkServerLogin(); 
          if (user) {
            store.commit('user/logIn', user)  // This might be unnecessary          
            console.log('Valid cookie found, logging in user: ', user);
            return true;
          } else {
            console.log('Invalid cookie');
            return false;
          }
        } catch (error) {
          console.log('Error validating cookie:', error);
          return false; 
        }
    }
    return true;
  }
  
async checkServerLogin() {
    try {
    const requestOptions = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json', 'User': store.getters['user/email'] },
        credentials: 'include',
    };

    const response = await fetch(
        `${import.meta.env.VITE_CHAT_SERVER_URL}/checklogin`,
        requestOptions
    );

    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }

    const json = await response.json();
    return json["user"];
    } catch (error) {
    console.error(error.message);
    throw error;
    }
}
}

export default new LoginStatusCheckService();