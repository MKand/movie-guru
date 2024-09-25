// Cookies: https://codesandbox.io/p/sandbox/vuex-persistedstate-with-js-cookie-0rjwk?file=%2Findex.js%3A31%2C7
// https://pusher.com/tutorials/authentication-vue-vuex/#login-action

export const store = {
  namespaced: true,
  state: {
    loggedIn: false,
    email: null,
    accessToken: null
  },
  mutations: {
    logIn(state, email){
        state.email = email
        state.loggedIn = true
    },
    logOut(state){
        state.email = null
        state.loggedIn = false
    }
  },
  getters: {
    loginStatus (state) {
      return state.loggedIn
    },
    email (state) {
      return state.email
    },
  },
}
export default store;
