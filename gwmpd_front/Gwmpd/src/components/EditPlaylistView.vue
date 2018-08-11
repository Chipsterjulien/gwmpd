<template lang="html">
  <div class="">
    <br>
    <br>
    Nom: {{ playlistName }}
    <br>
    <button type="button" @click="clearPlaylist">Clear</button><br>
    <input type="text" v-model="newPlaylistName">
    <button type="button" @click="renamePlaylist">Rename</button><br>
    <br>
    <div class="">
      <form class="" @submit.prevent="addURL">
        <label for="urlWebRadio">Webradio's URL</label>
        <input v-model="webradioURL" type="text" id="urlWebRadio">
        <button type="submit">add</button>
      </form>
    </div>
    <br>
    <div class="">
      <table v-if="playlist.length">
        <tr>
          <th>#</th>
          <th>File</th>
          <th>Song's name</th>
          <th>Duration</th>
        </tr>
        <tr v-for="(k, v) in playlist" :key="v">
          <td>{{ v + 1 }}</td>
          <td>{{ k.File }}</td>
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
        <table v-if="available.directories">
          <tr>
            <th>Directory's name</th>
          </tr>
          <tr v-for="(k, v) in available.directories" :key="v">
            <td @click="checkFilesList(k)">{{ k }}</td>
            <td><button type="button" @click="addSongToPlaylist(k)">add</button></td>
          </tr>
        </table>
      </div>
      <div class="">
        <table v-if="available.songs.length">
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
      available: {
        directories: [],
        songs: []
      },
      location: '',
      newPlaylistName: '',
      playlist: [],
      playlistName: '',
      webradioURL: ''
    }
  },
  computed: {
  },
  methods: {
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
      this.axios.post('v1/addSongToPlaylist', {
        songFilename: this.webradioURL,
        playlistName: this.playlistName
      })
        .then(response => {
          this.getPlaylist()
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
      this.axios.post('v1/clearPlaylist', {playlistName: this.playlistName})
        .then(response => {
          this.playlist = {}
        })
    },
    getFilesList () {
      this.axios.get('v1/filesList', {params: {location: this.location}})
        .then(response => {
          this.available = response.data
        })
    },
    getPlaylist () {
      this.axios.get('v1/playlistSongsList', {params: {playlistName: this.playlistName}})
        .then(response => {
          this.playlist = response.data
        })
    },
    moveBottom (actualPos) {
      this.axios.post('v1/moveSong', {
        playlistName: this.playlistName,
        oldPos: actualPos,
        newPos: this.playlist.length - 1
      })
        .then(response => {
          this.getPlaylist()
        })
    },
    moveTop (actualPos) {
      this.axios.post('v1/moveSong', {
        playlistName: this.playlistName,
        oldPos: actualPos,
        newPos: 0
      })
        .then(response => {
          this.getPlaylist()
        })
    },
    moveDown (actualPos) {
      this.axios.post('v1/moveSong', {
        playlistName: this.playlistName,
        oldPos: actualPos,
        newPos: actualPos + 1
      })
        .then(response => {
          this.getPlaylist()
        })
    },
    moveUp (actualPos) {
      this.axios.post('v1/moveSong', {
        playlistName: this.playlistName,
        oldPos: actualPos,
        newPos: actualPos - 1
      })
        .then(response => {
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
      this.axios.post('v1/removeSong', {
        playlistName: this.playlistName,
        pos: actualPos
      })
        .then(response => {
          this.getPlaylist()
        })
    },
    renamePlaylist () {
      this.axios.post('v1/renamePlaylist', {
        oldName: this.playlistName,
        newName: this.newPlaylistName
      })
        .then(response => {
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
