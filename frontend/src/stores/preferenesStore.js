import { createStore } from 'vuex'

export const store = {
  namespaced: true,
    state: {
        preferences:{"likes":{
            "genres":[],
            "actors":[],
            "director":[],
            "other":[]
        }, "dislikes":{
            "genres":[],
            "actors":[],
            "director":[],
            "other":[]
        }},
    },
    getters: {
      preferences (state) {
        return state.preferences
      },
      
      checkMovieExists: (state) => (value) => {
        let ls= state.preferences["likes"]["other"]
        if(ls){
          return ls.includes(value)
        }
        return false
      },
    },
    mutations: {
        update(state, preferences) {
          // mutate state
          state.preferences= preferences
        },

        add(state, target) {
          if(!state.preferences[target.type]){
            state.preferences[target.type] = {};
          }
          if(!state.preferences[target.type][target.key]){
            state.preferences[target.type][target.key] = [];
          }
          state.preferences[target.type][target.key].push(target.value)

        },
        checkMovieExists (state, value) {
          let ls= state.preferences["likes"]["other"]
          if(ls){
            return ls.includes(value)
          }
          return false
        },

        delete(state, target) {
          // mutate state
          // Check if the target key exists in the preferences object
          if (state.preferences[target.type][target.key]) {
            // Use filter to create a new array without the target value
            state.preferences[target.type][target.key] = state.preferences[target.type][target.key].filter(
              (value) => value !== target.value
            );
          }
        },
      
      }


    }
  
    export default store;