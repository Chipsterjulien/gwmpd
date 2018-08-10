<template>
  <div class="">
    <div class="">
      File: {{ currentSong.file }} <br>
      Title song: {{ currentSong.Title }} <br>
      Album: {{ currentSong.Album }} <br>
      Groupe: {{ currentSong.Artist }} <br>
      Consommé: {{ status.elapsed }}s <br>
      Temps total: {{ currentSong.Time }}s <br>
    </div>
    <br>
    <br>
    <div class="">
      <table v-if="currentPlaylist.length">
        <tr>
          <th>#</th>
          <th>File</th>
          <th>Title</th>
          <th>Duration</th>
        </tr>
        <tr v-for="(k, v) in currentPlaylist" :key="v">
          <td>{{ k.Pos + 1 }}</td>
          <td>{{ k.File }}</td>
          <td>{{ k.Title }}</td>
          <td>{{ k.Time }}</td>
          <td v-if="currentSong.Id !== k.Id"><button @click="playSong(k.Id, k.Pos)">play</button></td>
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
      'setPlaylist',
      'setState',
      'setID'
    ]),
    playSong (id, pos) {
      this.axios.get('v1/playSong', {params: {pos: pos}})
        .then(response => {
          this.setID(id)
          if (this.songPlayed !== true) {
            this.setState('play')
            this.songPlayed = true
          }
        })
    }
  },
  mounted () {
    this.axios.get('v1/currentPlaylist')
      .then(response => {
        this.setPlaylist(response.data)
      })
  }
}
</script>

<style lang="scss" scoped>
</style>
