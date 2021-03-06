<template>
  <div class="" v-if="getConnectionStatus === true">
    <b-container class="backgroundAndRoundedBox" v-if="status.state !== 'stop'">
      <div class="text-center" v-if="currentSong.Title === '' && currentSong.Album === '' && currentSong.Artist === ''">
        <h5>&nbsp;</h5>
        <h2 class="truncateLongText">{{ currentSong.file }}</h2>
        <h5>&nbsp;</h5>
      </div>
      <div class="text-center" v-else>
        <h5 class="truncateLongText"><span v-if="currentSong.Album !== ''">{{ currentSong.Album }}</span><span v-else>&nbsp;</span></h5>
        <h2 class="truncateLongText"><span v-if="currentSong.Title !== ''">{{ currentSong.Title }}</span><span v-else>&nbsp;</span></h2>
        <h5 class="truncateLongText"><span v-if="currentSong.Artist !== ''">{{ currentSong.Artist }}</span><span v-else>&nbsp;</span></h5>
      </div>
    </b-container>

    <div class="musicSliderAndNumber" v-if="status.duration !== 0">
      <b-container fluid>
        <b-row class="currentSongState">
          <b-col class="text-left"><strong>{{ getMusicElapsed }}</strong></b-col>
          <b-col class="text-right"><strong>{{ getMusicDuration }}</strong></b-col>
        </b-row>
      </b-container>
      <!-- Don't use b-form-input otherwise, music will be splited every seconds -->
      <input type="range" class="musicSlider" b-tooltip.hover :title="getMusicElapsed" min="0" :max="status.duration" v-model.number="musicValue"><br>
    </div>

    <!-- Array of playlist -->
    <div v-if="currentPlaylist.length > 0">
      <b-table stacked="md" striped hover :items="currentPlaylist" :fields="fields">
        <template slot="File" slot-scope="data">
          <span class="toLongFilenameSong">{{ data.item.File }}</span>
        </template>
        <template slot="buttonPlayMusic" slot-scope="data">
          <span v-if="currentSong.Id !== data.item.ID" class="buttonAlignRight"><b-button v-b-tooltip.hover.top title="Play song" class="icon-play_arrow" @click="playSong(data.item.ID, data.item.Pos)"></b-button></span>
        </template>
      </b-table>
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
      connected: false,
      fields: [{key: 'File', label: 'Filename'}, 'Duration', {key: 'buttonPlayMusic', label: ''}]
    }
  },
  computed: {
    ...mapGetters({
      currentPlaylist: 'getCurrentPlaylist',
      currentSong: 'getCurrentSongInfos',
      getConnectionStatus: 'getConnectionStatus',
      status: 'getStatusInfos'
    }),
    getMusicDuration () {
      return this.convertSecondsToString(this.status.duration)
    },
    getMusicElapsed () {
      return this.convertSecondsToString(this.status.elapsed)
    },
    musicValue: {
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
      'setCurrentView',
      'setPlaylist',
      'setPosition',
      'setState',
      'setID'
    ]),
    convertSecondsToString (time) {
      // Hours, minutes and seconds
      var hrs = ~~(time / 3600)
      var mins = ~~((time % 3600) / 60)
      var secs = ~~time % 60

      // Output like "1:01" or "4:03:59" or "123:03:59"
      var ret = ''

      if (hrs > 0) {
        ret += '' + hrs + ':' + (mins < 10 ? '0' : '')
      }

      ret += '' + mins + ':' + (secs < 10 ? '0' : '')
      ret += '' + secs

      return ret
    },
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
    if (localStorage.url !== 'undefined') {
      this.url = localStorage.url
      this.axios.defaults.baseURL = this.url
    }

    this.setCurrentView('QueueView')
    this.axios.get('v1/currentPlaylist')
      .then(response => {
        this.setPlaylist(response.data)
      })
  }
}
</script>

<style lang="scss" scoped>
  .backgroundAndRoundedBox {
    background-color: white;
    border-radius: 25px;
  }
  .buttonAlignRight {
    @media screen and (max-width: 768px) {
      float: right;
    }
  }

  .currentSongState {
    padding-bottom: 10px;
  }

  .musicSlider {
    :hover {
      opacity: 1;
    }
    -webkit-appearance: none;
    width: 100%;
    height: 10px;
    border-radius: 5px;
    background: #9E9E9E;
    outline: none;
    opacity: 0.7;
    -webkit-transition: .2s;
    transition: opacity .2s;
  }
  .musicSlider::-moz-range-thumb {
    width: 25px;
    height: 25px;
    border-radius: 50%;
    background: #1E88E5;
    cursor: pointer;
  }
  .musicSlider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 25px;
    height: 25px;
    border-radius: 50%;
    background: #1E88E5;
    cursor: pointer;
  }

  .musicSliderAndNumber {
    padding-bottom: 20px;
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

  .toLongFilenameSong {
    word-break: break-all;
  }

  .truncateLongText {
    width: 100%;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
