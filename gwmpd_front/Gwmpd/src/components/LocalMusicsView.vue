<template lang="html">
  <div class="positionning" v-if="getConnectionStatus === true">
    <b-input-group prepend="Location">
      <b-form-input v-model="location" placeholder="/" readonly></b-form-input>
      <b-input-group-append>
          <b-button @click="pathDown" v-if="location === ''" class="icon-undo" disabled></b-button>
          <b-button v-b-tooltip.hover.top title="Go back in the tree" @click="pathDown" v-else class="icon-undo"></b-button>
        </b-input-group-append>
    </b-input-group>

    <!-- Table of folders -->
    <div>
      <b-table stacked="md" striped hover v-if="available.directories.length > 0" :items="available.directories" :fields="directoriesFields">
        <template slot="nameFolder" slot-scope="data">
          <span class="toLongFilenameSong">{{ data.item.Name }}</span>
        </template>
        <template slot="button" slot-scope="data">
          <div class="alignButtonInTable">
            <b-button v-b-tooltip.hover.top title="Go inside" @click="checkFilesList(data.item.Name)" class="icon-visibility buttonMarginRight"></b-button>
            <b-button v-b-tooltip.hover.top title="Add folder to the current playlist" @click="addSongToCurrentPlaylist(data.item.Name)" class="icon-queue_music"></b-button>
            <b-button v-b-tooltip.hover.top title="Replace the current playlist by this folder" @click="replaceCurrentPlaylist(data.item.Name)" class="icon-add"></b-button>
          </div>
        </template>
      </b-table>
    </div>

    <!-- Table of songs -->
    <div>
      <b-table stacked="md" striped hover v-if="available.songs.length > 0" :items="available.songs" :fields="songsFields">
        <template slot="Song" slot-scope="data">
          <span class="toLongFilenameSong">{{ data.item.File }}</span>
        </template>
        <template slot="Duration" slot-scope="data">
          <span>{{ data.item.Duration }}</span>
        </template>
        <template slot="addSong" slot-scope="data">
          <div class="alignButtonInTable">
            <b-button v-b-tooltip.hover.top title="Add music to the current playlist" @click="addSongToCurrentPlaylist(data.item.File)" class="icon-queue_music"></b-button>
            <b-button v-b-tooltip.hover.top title="Replace the current playlist by this music" @click="replaceCurrentPlaylist(data.item.File)" class="icon-add"></b-button>
          </div>
        </template>
      </b-table>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'LocalMusicsView',
  data () {
    return {
      available: {
        directories: [],
        songs: []
      },
      directoriesFields: [{key: 'nameFolder', label: 'Folder'}, {key: 'button', label: ''}],
      location: '',
      songsFields: [{key: 'Song', label: 'Filename'}, 'Duration', {key: 'addSong', label: ''}]
    }
  },
  computed: {
    ...mapGetters({
      getConnectionStatus: 'getConnectionStatus',
      getStatus: 'getStatusInfos'
    })
  },
  methods: {
    ...mapActions([
      'setCurrentView',
      'setPlaylist',
      'setState'
    ]),
    addSongToCurrentPlaylist (filename) {
      if (this.location !== '') {
        filename = this.location + '/' + filename
      }
      this.axios.post('v1/addSongToCurrentPlaylist', {songFilename: filename})
        .then(response => {
          this.getPlaylist()
          if (this.getStatus.state !== 'play') {
            this.axios.get('v1/playSong')
              .then(response => {
                this.setState('play')
              })
          }
        })
    },
    checkFilesList (loc) {
      if (this.location === '') {
        this.location = loc
      } else {
        this.location += '/' + loc
      }

      this.getFilesList()
    },
    clearQueue () {
      this.axios.get('v1/clearCurrentPlaylist')
        .then(response => {
          this.setPlaylist({})
        })
    },
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
    getFilesList () {
      this.axios.get('v1/filesList', {params: {location: this.location}})
        .then(response => {
          this.available = response.data
          var tmp = []
          var i
          for (i = 0; i < this.available.directories.length; i++) {
            tmp.push({'Name': this.available.directories[i]})
          }
          this.available.directories = tmp
          tmp = []

          for (i = 0; i < this.available.songs.length; i++) {
            this.available.songs[i].Duration = this.convertSecondsToString(this.available.songs[i].Duration)
          }
        })
    },
    getPlaylist () {
      this.axios.get('v1/playlistSongsList', {params: {playlistName: this.playlistName}})
        .then(response => {
          var i

          this.playlist = response.data
          for (i = 0; i < this.playlist.length; i++) {
            this.playlist[i].Order = i
            this.playlist[i].Duration = this.convertSecondsToString(this.playlist[i].Duration)
          }
        })
    },
    replaceCurrentPlaylist (filename) {
      // Stop music
      this.axios.get('v1/stopSong')
        .then(response => {
          this.setState('stop')
          // Clear queue
          this.axios.get('v1/clearCurrentPlaylist')
            .then(response => {
              this.setPlaylist({})

              // Add song/folder to current playlist
              if (this.location !== '') {
                filename = this.location + '/' + filename
              }
              this.axios.post('v1/addSongToCurrentPlaylist', {songFilename: filename})
                .then(response => {
                  this.getPlaylist()
                  // Play song if not
                  if (this.getStatus.state !== 'play') {
                    this.axios.get('v1/playSong')
                      .then(response => {
                        this.setState('play')
                      })
                  }
                })
            })
        })
    },
    pathDown () {
      if (location !== '') {
        let locArray = this.location.split('/')
        this.location = locArray.slice(0, locArray.length - 1).join('/')
        this.getFilesList()
      }
    }
  },
  mounted () {
    this.setCurrentView('LocalMusicsView')
    this.getPlaylist()
    this.getFilesList()
  }
}
</script>

<style lang="scss" scoped>
  .addingPlaylist {
    padding-bottom: 5px;
  }

  .alignButtonInTable {
    float: right;
  }

  .buttonMarginRight {
    margin-right: 12px;
  }
</style>
