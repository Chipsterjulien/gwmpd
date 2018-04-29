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
      <button type="button" name="previousSong" @click="previousSong">Previous</button>
      <button type="button" name="stopSong" @click="stopSong">Stop</button>
      <button type="button" name="playSong" @click="playSong" v-if="!songPlayed">Play</button>
      <button type="button" name="pauseSong" @click="pauseSong" v-if="songPlayed">Pause</button>
      <button type="button" name="forwardSong" @click="forwardSong">Forward</button>
    </div>
    <div class="">
      <!-- Volume du song -->
      <div class="">
        Volume: {{ player.volume }}
      </div>
      <div class="">
        <button type="button" name="muteVolume" @click="toggleMuteVolume">Mute</button>
        <button type="button" name="lessVolume" @click="lessVolume">Less</button>
        <button type="button" name="moreVolume" @click="moreVolume">More</button>
      </div>
    </div>
    <div class="">
      <router-view/>
      <router-view name="SideBar"/>
    </div>
  </div>
</template>

<script>
// import store from './store'
export default {
  name: 'App',
  // store: store,
  data () {
    return {
      appName: 'Gwmpd',
      connected: false,
      songPlayed: false,
      player: {}
    }
  },
  methods: {
    lessVolume () {
      this.$changeVolume = this.$resource('v1/changeVolume')
      this.$changeVolume.save({volume: this.player.volume - 5}).then((response) => {
        this.player.volume -= 5
        if (this.player.volume < 0) {
          this.player.volume = 0
        }
      })
    },
    moreVolume () {
      this.$changeVolume = this.$resource('v1/changeVolume')
      this.$changeVolume.save({volume: this.player.volume + 5}).then((response) => {
        this.player.volume += 5
        if (this.player.volume > 100) {
          this.player.volume = 100
        }
      })
    },
    toggleMuteVolume () {
      this.$toggleVolume = this.$resource('v1/toggleMuteVolume')
      this.$toggleVolume.update().then((response) => {
        this.player.volume = response.data.volume
      })
    },
    forwardSong () {
      this.$nextSong = this.$resource('v1/nextSong')
      this.$nextSong.get()
    },
    pauseSong () {
      if (this.player.state !== 'pause') {
        this.$pauseSong = this.$resource('v1/pauseSong')
        this.$pauseSong.get().then((response) => {
          this.player.state = 'pause'
          this.songPlayed = false
        })
      }
    },
    playSong () {
      if (this.player.state !== 'play') {
        this.$stateMPD = this.$resource('v1/playSong')
        this.$stateMPD.get().then((response) => {
          this.player.state = 'play'
          this.songPlayed = true
        })
      }
    },
    previousSong () {
      this.$previousSong = this.$resource('v1/previousSong')
      this.$previousSong.get()
    },
    stopSong () {
      this.$stopSong = this.$resource('v1/stopSong')
      this.$stopSong.get().then((response) => {
        this.player.state = 'stop'
        this.songPlayed = false
      })
    }
  },
  mounted () {
    // console.log(this.$store)
    // this.$store.commit('CHANGE_DATA', 'coin')
    this.$interval = setInterval(() => {
      this.$stateMPD = this.$resource('v1/stateMPD')
      this.$stateMPD.get().then((response) => {
        this.connected = true
        this.player = response.data
        if (this.player.state === 'play') {
          this.songPlayed = true
        } else {
          this.songPlayed = false
        }
      }, (response) => {
        this.connected = false
      })
    }, 1000)
  },
  beforeUpdate () {
  },
  destroy () {
    clearInterval(this.$interval)
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
