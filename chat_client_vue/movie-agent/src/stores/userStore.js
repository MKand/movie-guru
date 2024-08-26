// Cookies: https://codesandbox.io/p/sandbox/vuex-persistedstate-with-js-cookie-0rjwk?file=%2Findex.js%3A31%2C7
// https://pusher.com/tutorials/authentication-vue-vuex/#login-action

export const store = {
  namespaced: true,
  state: {
    loggedIn: false,
    userName: null,
    email: null,
    accessToken: null
  },
  mutations: {
    logIn(state, result){
        state.userName = result.displayName
        state.email = result.email
        state.loggedIn = true
        state.accessToken = result.accessToken
    },
    logOut(state){
        state.userName = null
        state.email = null
        state.loggedIn = false
        state.accessToken = null
    }
  },
  getters: {
    loginStatus (state) {
      return state.loggedIn
    },
    userName (state) {
      return state.userName
    },
    email (state) {
      return state.email
    },
    accessToken (state) {
      return state.accessToken
    },
  },
}
export default store;
