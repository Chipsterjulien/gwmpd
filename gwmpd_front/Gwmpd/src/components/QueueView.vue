<template>
  <div class="" v-if="getConnectionStatus === true">
    <div class="">
      File: {{ currentSong.file }} <br>
      Title song: {{ currentSong.Title }} <br>
      Album: {{ currentSong.Album }} <br>
      Groupe: {{ currentSong.Artist }} <br>
      Consommé: {{ status.elapsed }}s <br>
      Temps total: {{ status.duration }}s <br>
    </div>
    <br>
    <br>
    <div class="" v-if="status.duration !== 0">
      Temps total: {{ status.duration }}<br>
      <input type="range" min="0" :max="status.duration" v-model.number="sliderValue"><br>
      SliderValue: {{ sliderValue }}<br>
      <button type="button" @click="moveBackwardsInTime">-10s</button>
      <button type="button" @click="moveForwardInTime">+10s</button>
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
          <td v-if="currentSong.Id !== k.ID"><button @click="playSong(k.ID, k.Pos)">play</button></td>
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
      connected: false
    }
  },
  computed: {
    ...mapGetters({
      currentPlaylist: 'getCurrentPlaylist',
      currentSong: 'getCurrentSongInfos',
      getConnectionStatus: 'getConnectionStatus',
      status: 'getStatusInfos'
    }),
    sliderValue: {
      get: function () {
        return this.status.elapsed
      },
      set: function (timePosition) {
        this.setPosition(timePosition)
        this.axios.post('v1/setPositionTimeInCurrentSong', {position: timePosition})
          .then(response => {
            this.setPosition(timePosition)
          })
      }
    }
  },
  methods: {
    ...mapActions([
      'setPlaylist',
      'setPosition',
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
    },
    moveBackwardsInTime () {
      var timePosition = this.status.elapsed - 10

      if (timePosition <= 0) {
        this.axios.post('v1/setPositionTimeInCurrentSong', {position: 0})
          .then(response => {
            this.setPosition(0)
          })
      } else {
        this.axios.post('v1/setPositionTimeInCurrentSong', {position: timePosition})
          .then(response => {
            this.setPosition(timePosition)
          })
      }
    },
    moveForwardInTime () {
      var timePosition = this.status.elapsed + 10

      if (timePosition > this.status.duration) {
        this.axios.get('v1/nextSong')
      } else {
        this.axios.post('v1/setPositionTimeInCurrentSong', {position: timePosition})
          .then(response => {
            this.setPosition(timePosition)
          })
      }
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
