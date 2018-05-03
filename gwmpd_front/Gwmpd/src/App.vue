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
      <!-- song's volume -->
      <div class="">
        Volume: {{ getStatus.volume }}
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

import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'App',
  data () {
    return {
      appName: 'Gwmpd',
      connected: false,
      songPlayed: false
    }
  },
  computed: {
    ...mapGetters({
      getStatus: 'getStatusInfos'
    })
  },
  methods: {
    ...mapActions([
      'setAllStatus',
      'setSong',
      'setVolume',
      'setState'
    ]),
    lessVolume () {
      this.$resource('v1/setVolume').save({volume: this.getStatus.volume - 5}).then((response) => {
        this.setVolume(response.data.volume)
      })
    },
    moreVolume () {
      this.$resource('v1/setVolume').save({volume: this.getStatus.volume + 5}).then((response) => {
        this.setVolume(response.data.volume)
      })
    },
    toggleMuteVolume () {
      this.$resource('v1/toggleMuteVolume').update().then((response) => {
        this.setVolume(response.data.volume)
      })
    },
    forwardSong () {
      this.$resource('v1/nextSong').get()
    },
    pauseSong () {
      if (this.getStatus.state !== 'pause') {
        this.$resource('v1/pauseSong').get().then((response) => {
          this.setState('pause')
          this.songPlayed = false
        })
      }
    },
    playSong () {
      if (this.getStatus.state !== 'play') {
        this.$resource('v1/playSong').get().then((response) => {
          this.setState('play')
          this.songPlayed = true
        })
      }
    },
    previousSong () {
      this.$resource('v1/previousSong').get()
    },
    stopSong () {
      if (this.getStatus.state !== 'stop') {
        this.$resource('v1/stopSong').get().then((response) => {
          this.setState('stop')
          this.songPlayed = false
        })
      }
    }
  },
  mounted () {
    this.$interval = setInterval(() => {
      this.$resource('v1/statusMPD').get().then((response) => {
        this.setAllStatus(response.data)
        this.connected = true
        if (response.data.state === 'play') {
          this.$resource('v1/currentSong').get().then((response) => {
            this.setSong(response.data)
          })
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
