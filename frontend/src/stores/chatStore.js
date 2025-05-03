export const store = {
  namespaced: true,
    state: {
        chatMessageHistory:[],
        movies: [],
        placeHolderMovies: [],
        traceId: "",
        spanId: ""
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
      traceAndSpanIds (state) {
        return state.traceId, state.spanId
      },
    },
    mutations: {
        add(state, message) {
          // mutate state
          state.chatMessageHistory.push(message)
        },
        updateTraceSpanIds(state, traceId, spanId) {
          // mutate state
          state.traceId=traceId
          state.spanId= spanId
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
              element.poster= new URL("../assets/movie-guru.png", import.meta.url)
            }
            state.movies.push(element)
          });
        }
        },
        addPlaceHolderMovies(state, movies) {
          state.placeHolderMovies = []
          movies.forEach(element => {
            if (element.poster=="") {
              element.poster=new URL("../assets/movie-guru.png", import.meta.url)
            }
            state.placeHolderMovies.push(element)
          });
        },
        }
    }
  
    export default store;