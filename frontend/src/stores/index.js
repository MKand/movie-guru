import { createStore } from 'vuex'
import createPersistedState from 'vuex-persistedstate';

import {store as chatStore} from './chatStore'
import {store as userStore} from './userStore'
import {store as preferencesStore} from './preferenesStore'

export const store = createStore({
    modules: {
      chat: chatStore,
      user: userStore,
      preferences: preferencesStore
    },
    plugins: [
      createPersistedState({
        paths: ['user'],  // Persist the entire 'user' module's state
      }),
    ],
  })
export function init(){

  return Promise.all([]);

}
export default store;