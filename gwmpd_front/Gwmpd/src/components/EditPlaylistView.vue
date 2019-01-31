<template lang="html">
  <div class="" v-if="getConnectionStatus === true">
    <b-container>
      <b-input-group>
        <b-form-input v-model="newPlaylistName" @keyup.enter.native="renamePlaylist" :readonly="isReadOnly"></b-form-input>
        <b-input-group-append v-if="isReadOnly">
          <b-button v-b-tooltip.hover.top title="Edit playlist's name" variant="primary" @click="canEditPlaylistName" class="icon-mode_edit"></b-button>
        </b-input-group-append>
        <b-input-group-append v-else>
          <b-button v-b-tooltip.hover.top title="Save" variant="primary" @click="renamePlaylist" class="icon-save"></b-button>
          <b-button v-b-tooltip.hover.top title="Cancel" variant="primary" @click="cancelEditPlaylistName" class="icon-close"></b-button>
        </b-input-group-append>
      </b-input-group>
    </b-container>

    <div class="saveClearButton" align="right">
      <b-button size="lg" v-b-tooltip.hover.top title="Clear playlist" variant="danger" @click="clearPlaylist" class="icon-clear_all"></b-button>
    </div>

    <b-list-group>
      <draggable :list="playlist" :move="onMove" @start="isDragging=true" @end="isDragging=false">
        <b-list-group-item v-for="(element, index) in playlist" :key="element.Order" class="d-flex justify-content-between align-items-center">
          <span class="toLongFilenameSong">{{ element.File }}</span>
          <div class="buttonAlignRightInListGroup">
            <b-badge v-if="element.Duration !== '0:00'" class="addSpace" pill>{{ element.Duration }}</b-badge>
            <b-button v-b-tooltip.hover.top title="Remove from playlist" @click="removeSong(index)" class="icon-delete"></b-button>
          </div>
        </b-list-group-item>
      </draggable>
    </b-list-group>
    <br>

    <hr>

    <b-form-group>
      <b-input-group prepend="web's URL">
        <b-form-input v-model="webradioURL" @keyup.enter.native="addURL"></b-form-input>
        <b-input-group-append>
          <b-button v-b-tooltip.hover.top title="Add a web url" @click="addURL" class="icon-add"></b-button>
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <hr>

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
          <div class="buttonAlignRight">
            <b-button v-b-tooltip.hover.top title="Go inside" @click="checkFilesList(data.item.Name)" class="icon-visibility"></b-button>
            <b-button v-b-tooltip.hover.top title="Add all the songs from the folder" @click="addSongToPlaylist(data.item.Name)" class="icon-create_new_folder"></b-button>
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
          <b-button v-b-tooltip.hover.top title="Add song" @click="addSongToPlaylist(data.item.File)" class="icon-add buttonAlignRight"></b-button>
        </template>
      </b-table>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import draggable from 'vuedraggable'
export default {
  components: {
    draggable
  },
  name: 'EditPlaylistView',
  data () {
    return {
      available: {
        directories: [],
        songs: []
      },
      directoriesFields: [{key: 'nameFolder', label: 'Folder'}, {key: 'button', label: ''}],
      isDragging: false,
      isReadOnly: true,
      location: '',
      newPlaylistName: '',
      playlist: [],
      playlistName: '',
      songsFields: [{key: 'Song', label: 'Filename'}, 'Duration', {key: 'addSong', label: ''}],
      songOnMove: {},
      webradioURL: ''
    }
  },
  computed: {
    ...mapGetters({
      getConnectionStatus: 'getConnectionStatus'
    })
  },
  methods: {
    ...mapActions([
      'setCurrentView'
    ]),
    onMove ({ relatedContext, draggedContext }) {
      this.songOnMove = draggedContext
    },
    addSongToPlaylist (filename) {
      if (this.location !== '') {
        filename = this.location + '/' + filename
      }
      this.axios.post('v1/addSongToPlaylist', {
        songFilename: filename,
        playlistName: this.playlistName
      })
        .then(response => {
          this.getPlaylist()
        })
    },
    addURL () {
      if (this.webradioURL !== '') {
        this.axios.post('v1/addSongToPlaylist', {
          songFilename: this.webradioURL,
          playlistName: this.playlistName
        })
          .then(response => {
            this.getPlaylist()
            this.webradioURL = ''
          })
      }
    },
    canEditPlaylistName () {
      this.isReadOnly = false
    },
    cancelEditPlaylistName () {
      this.isReadOnly = true
      this.newPlaylistName = this.playlistName
    },
    checkFilesList (loc) {
      if (this.location === '') {
        this.location = loc
      } else {
        this.location += '/' + loc
      }

      this.getFilesList()
    },
    clearPlaylist () {
      this.axios.post('v1/clearPlaylist', {playlistName: this.playlistName})
        .then(response => {
          this.playlist = {}
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
    getMusicName (filename) {
      var filenameSplitted = filename.split('/')

      return filenameSplitted[filenameSplitted.length - 1]
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
    pathDown () {
      if (location !== '') {
        let locArray = this.location.split('/')
        this.location = locArray.slice(0, locArray.length - 1).join('/')
        this.getFilesList()
      }
    },
    removeSong (actualPos) {
      this.axios.post('v1/removeSong', {
        playlistName: this.playlistName,
        pos: actualPos
      })
        .then(response => {
          this.getPlaylist()
        })
    },
    renamePlaylist () {
      this.isReadOnly = true
      this.axios.post('v1/renamePlaylist', {
        oldName: this.playlistName,
        newName: this.newPlaylistName
      })
        .then(response => {
          this.playlistName = response.data.newName
          this.newPlaylistName = this.playlistName
          this.$router.push({ name: 'EditPlaylistView', params: {'playlistName': this.newPlaylistName} })
        })
    }
  },
  watch: {
    isDragging (newValue) {
      if (newValue === false) {
        if (this.songOnMove.index !== this.songOnMove.futureIndex) {
          this.axios.post('v1/moveSong', {
            playlistName: this.playlistName,
            oldPos: this.songOnMove.index,
            newPos: this.songOnMove.futureIndex
          })
            .then(response => {
              this.getPlaylist()
            })
        }
      }
    }
  },
  mounted () {
    if (localStorage.url !== 'undefined') {
      this.url = localStorage.url
      this.axios.defaults.baseURL = this.url
    }

    this.setCurrentView('EditPlaylistView')

    this.playlistName = this.$route.params.playlistName
    this.newPlaylistName = this.playlistName
    this.getPlaylist()
    this.getFilesList()
  }
}
</script>

<style lang="scss" scoped>
  .addSpace {
    margin-right: 3px;

    @media screen and (max-width: 768px) {
      margin-bottom: 4px;
    }
  }

  .buttonAlignRight {
      float: right;
  }

  .buttonAlignRightInListGroup {
    display: flex;
    justify-content: flex-end;
    align-items: center;

    @media screen and (max-width: 768px) {
      display: flex;
      flex-direction: column;
      margin-left: 14px;
    }
  }

  .saveClearButton {
    padding-top: 20px;
    padding-bottom: 3px;
  }

  .toLongFilenameSong {
    word-break: break-all;
  }
</style>
