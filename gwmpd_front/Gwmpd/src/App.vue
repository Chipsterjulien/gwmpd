<template>
  <div id="app">
    <div class="TopBar">
      {{ appName }}
      <router-link :to="{ name: 'QueueView', params: {} }">Queue</router-link>
      <router-link :to="{ name: 'PlaylistView', params: {} }">Playlist</router-link>
      <router-link :to="{ name: 'AboutView', params: {} }">About</router-link>
    </div>
    <div class="">
      <div class="" v-if="connected">
        Connected
      </div>
      <div class="" v-else>
        Disconnected
      </div>
    </div>
    <div class="">
      <button type="button" name="previous" @click="previous">Arrière</button>
      <button type="button" name="stop" @click="stop">Stop</button>
      <button type="button" name="play" @click="play" v-if="!songPlayed">Jouer</button>
      <button type="button" name="pause" @click="pause" v-if="songPlayed">Pause</button>
      <button type="button" name="forward" @click="forward">Avant</button>
    </div>
    <div class="">
      <router-view/>
      <router-view name="SideBar"/>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data () {
    return {
      appName: 'Gwmpd',
      // stat: [],
      connected: false,
      songPlayed: false,
      player: {}
    }
  },
  methods: {
    previous () {
      this.$previousSong = this.$resource('v1/previousSong')
      this.$previousSong.get().then((response) => {
        // success callback
      }, (response) => {
        // error callback
      })
    },
    stop () {
      // Prévoir d'arrêter le ticker créé dans play()
      this.$stopSong = this.$resource('v1/stopSong')
      this.$stopSong.get().then((response) => {
        this.player.state = 'stop'
        this.songPlayed = false
      }, (response) => {
        // error callback
      })
    },
    play () {
      // Créer un ticker
      this.$playSong = this.$resource('v1/playSong')
      this.$playSong.get().then((response) => {
        this.player.state = 'play'
        this.songPlayed = true
      }, (response) => {
        // error callback
      })
    },
    pause () {
      // Prévoir d'arrêter le ticker créé dans play()
      this.$pauseSong = this.$resource('v1/pauseSong')
      this.$pauseSong.get().then((response) => {
        this.player.state = 'pause'
        this.songPlayed = false
      }, (response) => {
        // error callback
      })
    },
    forward () {
      this.$nextSong = this.$resource('v1/nextSong')
      this.$nextSong.get().then((response) => {
        // success callback
      }, (response) => {
        // error callback
      })
    }
  },
  mounted () {
    this.$stateMPD = this.$resource('v1/stateMPD')
    this.$stateMPD.query().then((response) => {
      this.player = response.data
      this.connected = true
      if (this.player.state === 'play') {
        this.songPlayed = true
      } else {
        this.songPlayed = false
      }
    }, (response) => {
    })
  },
  beforeUpdate () {
  }
}
</script>

<style lang="scss">
#app {
  /* font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px; */
}
</style>
