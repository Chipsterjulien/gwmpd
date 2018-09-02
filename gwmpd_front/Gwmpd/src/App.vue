<template>
  <div id="app">
    <div class="" v-if="getConnectionStatus === true">
      <b-navbar toggleable="md" variant="info" type="dark">
        <b-navbar-brand>{{ appName }}</b-navbar-brand>
      </b-navbar>

      <div class="TopBar">
        {{ appName }}
        <router-link :to="{ name: 'QueueView', params: {} }">Queue</router-link>
        <router-link :to="{ name: 'PlaylistView', params: {} }">Playlists</router-link>
        <router-link :to="{ name: 'AboutView', params: {} }">About</router-link>
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
      <div class="" v-if="getStatus.error !== ''">
        Error: {{ getStatus.error }}
      </div>
    </div>
    <div class="" v-else>
      <b-alert show variant="warning">
        <strong>{{ appName }} is disconnected !</strong><br>
        <router-link :to="'Login'" class="alert-link">Please sign in</router-link>
      </b-alert>
    </div>
    <router-view/>
    <router-view name="SideBar"/>
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
      getConnectionStatus: 'getConnectionStatus',
      getCurrentSongInfos: 'getCurrentSongInfos',
      getStatus: 'getStatusInfos'
    })
  },
  methods: {
    refresh () {
      this.$auth.refresh({
        success: function (response) {
          this.$auth.token(null, response.data.token)
          sessionStorage.token = response.data.token
        },
        error: function (e) {
          delete sessionStorage.token
        }
      })
    },
    ...mapActions([
      'setAllStatus',
      'setConnectionStatus',
      'setSong',
      'setVolume',
      'setPlaylist',
      'setState'
    ]),
    lessVolume () {
      if (this.getStatus.volume >= 5) {
        this.axios.post('v1/setVolume', {
          volume: this.getStatus.volume - 5
        })
          .then(response => {
            this.setVolume(response.data.volume)
          })
      }
    },
    moreVolume () {
      if (this.getStatus.volume <= 95) {
        this.axios.post('v1/setVolume', {
          volume: this.getStatus.volume + 5
        })
          .then(response => {
            this.setVolume(response.data.volume)
          })
      }
    },
    toggleMuteVolume () {
      this.axios.put('v1/toggleMuteVolume')
        .then(response => {
          this.setVolume(response.data.volume)
        })
    },
    forwardSong () {
      this.axios.get('v1/nextSong')
    },
    pauseSong () {
      if (this.getStatus.state !== 'pause') {
        this.axios.get('v1/pauseSong')
          .then(response => {
            this.setState('pause')
            this.songPlayed = false
          })
      }
    },
    playSong () {
      if (this.getStatus.state !== 'play') {
        this.axios.get('v1/playSong')
          .then(response => {
            this.setState('play')
            this.songPlayed = true
          })
      }
    },
    previousSong () {
      this.axios.get('v1/previousSong')
    },
    stopSong () {
      if (this.getStatus.state !== 'stop') {
        this.axios.get('v1/stopSong')
          .then(response => {
            this.setState('stop')
            this.songPlayed = false
          })
      }
    }
  },
  mounted () {
    if (this.$auth.watch.authenticated) {
      // refresh token if F5 was sent
      this.refresh()
    }

    this.$refreshTokenInterval = setInterval(() => {
      this.refresh()
    }, 55000)

    this.$refreshMpdDataInterval = setInterval(() => {
      this.axios.get('v1/statusMPD').then((response) => {
        this.setAllStatus(response.data)
        this.setConnectionStatus(true)
        if (response.data.state === 'play') {
          this.axios.get('v1/currentSong').then((response) => {
            if (this.getCurrentSongInfos.file !== response.data.file) {
              this.axios.get('v1/currentPlaylist').then((response) => {
                this.setPlaylist(response.data)
              })
            }
            this.setSong(response.data)
          })
          this.songPlayed = true
        } else {
          this.songPlayed = false
        }
      }, (response) => {
        this.setConnectionStatus(false)
      })
    }, 1000)
  }
}
</script>

<style lang="scss">
// #app {
//   font-family: 'Avenir', Helvetica, Arial, sans-serif;
//   -webkit-font-smoothing: antialiased;
//   -moz-osx-font-smoothing: grayscale;
//   text-align: center;
//   color: #2c3e50;
//   margin-top: 60px;
// }
</style>
