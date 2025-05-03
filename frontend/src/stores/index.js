import { createStore } from 'vuex'

import {store as chatStore} from './chatStore'
import {store as userStore} from './userStore'
import {store as preferencesStore} from './preferenesStore'

export const store = createStore({
    modules: {
      chat: chatStore,
      user: userStore,
      preferences: preferencesStore
    },
  })
export function init(){

  return Promise.all([]);

}
export default store;