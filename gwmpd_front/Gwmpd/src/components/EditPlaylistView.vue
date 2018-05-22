<template lang="html">
  <div class="">
    <br>
    <br>
    Nom: {{ playlistName }}
    <br>
    <button type="button" @click="clearPlaylist">Clear</button><br>
    <input type="text" v-model="newPlaylistName">
    <button type="button" @click="renamePlaylist">Rename</button>
    <div class="">
      <table v-if="playlist.length">
        <tr>
          <th>#</th>
          <th>Song's name</th>
          <th>Duration</th>
        </tr>
        <tr v-for="(k, v) in playlist" :key="v">
          <td>{{ v + 1 }}</td>
          <td>{{ k.Title }}</td>
          <td>{{ k.Time }}</td>
          <td><button type="button" v-if="v !== 0" @click="moveTop(v)">Top</button></td>
          <td><button type="button" v-if="v < playlist.length - 1" @click="moveBottom(v)">Bottom</button></td>
          <td><button type="button" v-if="v !== 0" @click="moveUp(v)">Up</button></td>
          <td><button type="button" v-if="v < playlist.length - 1" @click="moveDown(v)">Down</button></td>
          <td><button type="button" @click="removeSong(v)">Remove</button></td>
        </tr>
      </table>
    </div>
    <div class="">
      <div class="">
        <br>
        Location: <span @click="pathDown">{{ location }}</span>
      </div>
      <div class="">
        <table v-if="Object.keys(available).length">
          <tr>
            <th>Directory's name</th>
          </tr>
          <tr v-for="(k, v) in available.directories" :key="v">
            <td @click="checkFilesList(k)">{{ k }}</td>
          </tr>
        </table>
      </div>
      <div class="">
        <table>
          <tr>
            <th>File's name</th>
            <th>Artist</th>
            <th>Album</th>
            <th>Duration</th>
          </tr>
          <tr v-for="(k, v) in available.songs" :key="v">
            <td>{{ k.File }}</td>
            <td>{{ k.Artist }}</td>
            <td>{{ k.Album }}</td>
            <td>{{ k.Time }}</td>
            <td><button type="button" @click="addSongToPlaylist(k.File)">Add</button></td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EditPlaylistView',
  data () {
    return {
      available: {},
      location: '',
      newPlaylistName: '',
      playlist: [],
      playlistName: ''
    }
  },
  computed: {
  },
  methods: {
    addSongToPlaylist (filename) {
      if (this.location !== '') {
        filename = this.location + '/' + filename
      }
      this.$resource('v1/addSongToPlaylist').save({songFilename: filename, playlistName: this.playlistName}).then((response) => {
        console.log('coin')
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
    clearPlaylist () {
      this.$resource('v1/clearPlaylist').save({playlistName: this.playlistName}).then((response) => {
        this.playlist = {}
      })
    },
    getFilesList () {
      this.$resource('v1/filesList{/location}').get({location: this.location}).then((response) => {
        this.available = response.data
      })
    },
    getPlaylist () {
      this.$resource('v1/playlistSongsList{/playlistName}').get({playlistName: this.playlistName}).then((response) => {
        this.playlist = response.data
      })
    },
    moveBottom (actualPos) {
      this.$resource('v1/moveSong').save({playlistName: this.playlistName, oldPos: actualPos, newPos: this.playlist.length - 1}).then((response) => {
        this.getPlaylist()
      })
    },
    moveTop (actualPos) {
      this.$resource('v1/moveSong').save({playlistName: this.playlistName, oldPos: actualPos, newPos: 0}).then((response) => {
        this.getPlaylist()
      })
    },
    moveDown (actualPos) {
      this.$resource('v1/moveSong').save({playlistName: this.playlistName, oldPos: actualPos, newPos: actualPos + 1}).then((response) => {
        this.getPlaylist()
      })
    },
    moveUp (actualPos) {
      this.$resource('v1/moveSong').save({playlistName: this.playlistName, oldpos: actualPos, newpos: actualPos - 1}).then((response) => {
        this.getPlaylist()
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
      this.$resource('v1/removeSong').save({playlistName: this.playlistName, pos: actualPos}).then((response) => {
        this.getPlaylist()
      })
    },
    renamePlaylist () {
      this.$resource('v1/renamePlaylist').save({oldName: this.playlistName, newName: this.newPlaylistName}).then((response) => {
        this.playlistName = response.data.newName
        this.newPlaylistName = this.playlistName
      })
    }
  },
  mounted () {
    this.playlistName = this.$route.params.playlistName
    this.newPlaylistName = this.playlistName
    this.getPlaylist()
    this.getFilesList()
  }
}
</script>

<style lang="scss" scoped>
</style>
