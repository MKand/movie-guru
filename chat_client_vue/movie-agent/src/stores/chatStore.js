export const store = {
  namespaced: true,
    state: {
        chatMessageHistory:[],
        movies: [],
        placeHolderMovies: [],
    },
    getters: {
      messages (state) {
        return state.chatMessageHistory
      },
      movies (state) {
        return state.movies
      },
      placeHolderMovies (state) {
        return state.placeHolderMovies
      },
    },
    mutations: {
        add(state, message) {
          // mutate state
          state.chatMessageHistory.push(message)
        },
        clear(state) {
          // mutate state
          state.chatMessageHistory = []
        },
        addMovies(state, movies) {
          if (movies.length > 0) {
          state.movies = []
          movies.forEach(element => {
            if (element.poster=="") {
              element.poster="https://storage.googleapis.com/generated_posters/notfound.png"
            }
            state.movies.push(element)
          });
        }
        },
        addPlaceHolderMovies(state, movies) {
          state.placeHolderMovies = []
          movies.forEach(element => {
            if (element.poster=="") {
              element.poster="https://storage.googleapis.com/generated_posters/notfound.png"
            }
            state.placeHolderMovies.push(element)
          });
        },
        }
    }
  
    export default store;