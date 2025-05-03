import {ref } from 'vue'
class PlayerService {

  playMovie = ref(false);
  poster = ref(null);

  setPoster(poster){
    this.playMovie.value = true;
    this.poster.value = poster;
  }
  unsetPoster(){
    this.playMovie.value = false;
    this.poster.value = null;
  }

}

export default new PlayerService();
