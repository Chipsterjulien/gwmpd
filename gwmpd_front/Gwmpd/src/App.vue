<template>
  <div id="app">
    <div class="myNavBar" v-if="getConnectionStatus === true">
      <b-navbar toggleable="md" class=navBar type="dark">
        <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
        <b-navbar-brand>{{ appName }}</b-navbar-brand>
        <b-collapse is-nav id="nav_collapse">
          <b-navbar-nav>
            <b-nav-item :to="{ name: 'QueueView', params: {} }">Queue</b-nav-item>
            <b-nav-item :to="{ name: 'PlaylistView', params: {} }">Playlists</b-nav-item>
            <b-nav-item :to="{ name: 'ConfigView', params: {} }" disabled>Config</b-nav-item>
            <b-nav-item :to="{ name: 'AboutView', params: {} }">About</b-nav-item>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>

      <b-container>
        <b-row class="positionning text-center">
          <b-col><b-button size="lg" @click="previousSong" class="icon-skip_previous"></b-button></b-col>
          <b-col><b-button size="lg" @click="stopSong" class="icon-stop"></b-button></b-col>
          <b-col><b-button size="lg" @click="playSong" v-if="!songPlayed" class="icon-play_arrow"></b-button><b-button size="lg" @click="pauseSong" v-if="songPlayed" class="icon-pause"></b-button></b-col>
          <b-col><b-button size="lg" @click="forwardSong" class="icon-skip_next"></b-button></b-col>
        </b-row>

        <b-row class="sound text-center">
          <b-col cols="3">
            <b-button size="lg" @click="toggleMuteVolume" class="icon-volume_off" v-if="volumeValue == 0"></b-button>
            <b-button size="lg" @click="toggleMuteVolume" class="icon-volume_mute" v-if="volumeValue > 0 && volumeValue < 30"></b-button>
            <b-button size="lg" @click="toggleMuteVolume" class="icon-volume_down" v-if="volumeValue > 29 && volumeValue < 60"></b-button>
            <b-button size="lg" @click="toggleMuteVolume" class="icon-volume_up" v-if="volumeValue > 59"></b-button>
          </b-col>
          <b-col cols="9" class="slidecontainer"><b-form-input id="volumeSlider" b-tooltip.hover :title="volumeValue" type="range" min="0" max="100" :step="5" class="slider" v-model.number="volumeValue"></b-form-input></b-col>
        </b-row>
      </b-container>

      <!-- <hr> -->

      <div v-if="getStatus.error !== ''">
        Error: {{ getStatus.error }}
      </div>
    </div>
    <div v-else>
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
    }),
    volumeValue: {
      get: function () {
        return this.getStatus.volume
      },
      set: function (volumePosition) {
        this.setVolume(volumePosition)
        this.axios.post('v1/setVolume', {volume: volumePosition})
          .then(response => {
            this.setVolume(volumePosition)
          })
      }
    }
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
    // lessVolume () {
    //   if (this.getStatus.volume >= 5) {
    //     this.axios.post('v1/setVolume', {
    //       volume: this.getStatus.volume - 5
    //     })
    //       .then(response => {
    //         this.setVolume(response.data.volume)
    //       })
    //   }
    // },
    // moreVolume () {
    //   if (this.getStatus.volume <= 95) {
    //     this.axios.post('v1/setVolume', {
    //       volume: this.getStatus.volume + 5
    //     })
    //       .then(response => {
    //         this.setVolume(response.data.volume)
    //       })
    //   }
    // },
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
  @import "./style";
  html, body {
    background-color: #E0E0E0;
  }

  #app {
    .myNavBar {
      padding-bottom: 20px;
    }
    .navBar {
      background-color: #1E88E5;
    }

    .positionning {
      margin-top: 20px;
    }

    .slider {
      :hover {
        opacity: 1;
      }
      -webkit-appearance: none;
      width: 100%;
      height: 15px;
      border-radius: 5px;
      background: #9E9E9E;
      outline: none;
      opacity: 0.7;
      -webkit-transition: .2s;
      transition: opacity .2s;
    }

    .slider::-moz-range-thumb {
      width: 25px;
      height: 25px;
      border-radius: 50%;
      background: #1E88E5;
      cursor: pointer;
    }

    .slider::-webkit-slider-thumb {
      -webkit-appearance: none;
      appearance: none;
      width: 25px;
      height: 25px;
      border-radius: 50%;
      background: #1E88E5;
      cursor: pointer;
    }

    .slidecontainer {
      margin-top: 12px;
      width: 100%;
    }

    .sound{
      margin-top: 20px;
    }

  // font-family: 'Avenir', Helvetica, Arial, sans-serif;
  //   -webkit-font-smoothing: antialiased;
  //   -moz-osx-font-smoothing: grayscale;
  //   text-align: center;
  //   color: #2c3e50;
  //   margin-top: 60px;
  }
</style>
