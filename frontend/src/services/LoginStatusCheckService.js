import store  from '../stores';

class LoginStatusCheckService {

async checkLogin() {
    const loggedIn = store.getters['user/loginStatus'];
    if (!loggedIn) {
        try {
          const user = await this.checkServerLogin(); 
          if (user) {
            return true;
          } else {
            console.log('Invalid cookie, redirecting to login.');
            return false;
          }
        } catch (error) {
          console.error('Error validating cookie:', error);
          return false; 
        }
    }
  }
  
async checkServerLogin() {
    try {
    const requestOptions = {
        method: 'GET',
        headers: {
        'Content-Type': 'application/json',
        },
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