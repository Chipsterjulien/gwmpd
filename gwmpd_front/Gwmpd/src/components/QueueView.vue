<template>
  <div class="">
    <div class="">
      Title song: {{ currentSong.title }} <br>
      Album: {{ currentSong.album }} <br>
      Groupe: {{ currentSong.artist }} <br>
      Consommé: {{ currentSong.elapsed }}s <br>
      Temps total: {{ currentSong.timeSong }}s <br>
    </div>
    <div class="">
      <router-view name='SideBar'/>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'QueueView',
  data () {
    return {
    }
  },
  computed: {
    ...mapGetters({
      currentSong: 'getCurrentSongInfos'
    })
  },
  methods: {
    ...mapActions([
      'setPlaylist'
    ])
  },
  mounted () {
    this.$resource('v1/currentPlaylist').get().then((response) => {
      console.log(response.data)
      this.setPlaylist(response.data)
    })
  }
}
</script>

<style lang="scss" scoped>
</style>
