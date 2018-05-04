<template>
  <div class="">
    <div class="">
      Title song: {{ currentSong.Title }} <br>
      Album: {{ currentSong.Album }} <br>
      Groupe: {{ currentSong.Artist }} <br>
      Consommé: {{ status.elapsed }}s <br>
      Temps total: {{ currentSong.Time }}s <br>
    </div>
    <br>
    <br>
    <div class="">
      <table>
        <tr>
          <th>Title</th>
        </tr>
        <tr v-for="(k,v) in currentPlaylist" :key="v">
          <td>{{ k.Title }}</td>
        </tr>
      </table>
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
      currentSong: 'getCurrentSongInfos',
      status: 'getStatusInfos',
      currentPlaylist: 'getCurrentPlaylist'
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
