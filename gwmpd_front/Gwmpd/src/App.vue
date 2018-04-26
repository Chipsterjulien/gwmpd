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
export default {
  name: 'App',
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
      let volTmp = this.player.volume
      volTmp -= 5
      if (volTmp < 0) {
        volTmp = 0
      }

      if (volTmp !== this.player.volume) {
        this.$changeVolume = this.$resource('v1/changeVolume')
        this.$changeVolume.save({volume: volTmp}).then((response) => {
          this.player.volume -= 5
        })
      }
    },
    moreVolume () {
      let volTmp = this.player.volume
      volTmp += 5
      if (volTmp > 100) {
        volTmp = 100
      }
      if (volTmp !== this.player.volume) {
        this.$changeVolume = this.$resource('v1/changeVolume')
        this.$changeVolume.save({volume: volTmp}).then((response) => {
          this.player.volume += 5
        })
      }
    },
    toggleMuteVolume () {
      this.$toggleVolume = this.$resource('v1/toggleMuteVolume')
      this.$toggleVolume.update().then((response) => {
        this.player.volume = response.data.volume
      })
    },
    forwardSong () {
      this.$nextSong = this.$resource('v1/nextSong')
      this.$nextSong.get().then((response) => {
        // success callback
      }, (response) => {
        // error callback
      })
    },
    pauseSong () {
      // Prévoir d'arrêter le ticker créé dans play()
      this.$pauseSong = this.$resource('v1/pauseSong')
      this.$pauseSong.get().then((response) => {
        this.player.state = 'pause'
        this.songPlayed = false
        clearInterval(this.$interval)
      }, (response) => {
        // error callback
      })
    },
    playSong () {
      this.$interval = setInterval(() => {
        this.$stateMPD = this.$resource('v1/stateMPD')
        this.$stateMPD.get().then((response) => {
          this.player = response.data

          if (this.player.state !== 'play') {
            clearInterval(this.$interval)
            this.songPlayed = false
          }
        })
      }, 1000)
      this.$playSong = this.$resource('v1/playSong')
      this.$playSong.get().then((response) => {
        this.player.state = 'play'
        this.songPlayed = true
      }, (response) => {
        // error callback
      })
    },
    previousSong () {
      this.$previousSong = this.$resource('v1/previousSong')
      this.$previousSong.get().then((response) => {
        // success callback
      }, (response) => {
        // error callback
      })
    },
    stopSong () {
      // Prévoir d'arrêter le ticker créé dans play()
      this.$stopSong = this.$resource('v1/stopSong')
      this.$stopSong.get().then((response) => {
        this.player.state = 'stop'
        this.songPlayed = false
        clearInterval(this.$interval)
      }, (response) => {
        // error callback
      })
    }
  },
  mounted () {
    this.$stateMPD = this.$resource('v1/stateMPD')
    this.$stateMPD.get().then((response) => {
      this.player = response.data
      this.connected = true
      if (this.player.state === 'play') {
        this.songPlayed = true
        this.$interval = setInterval(() => {
          this.$stateMPD = this.$resource('v1/stateMPD')
          this.$stateMPD.get().then((response) => {
            this.player = response.data
            if (this.player.state !== 'play') {
              clearInterval(this.$interval)
              this.songPlayed = false
            }
          })
        }, 1000)
      } else {
        this.songPlayed = false
      }
    }, (response) => {
    })
  },
  beforeUpdate () {
  },
  destroy () {
    if (this.$interval) {
      clearInterval(this.$interval)
    }
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
