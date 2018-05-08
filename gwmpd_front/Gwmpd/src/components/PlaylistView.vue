<template lang="html">
  <div class="">
    <br>
    <br>

    <div class="">
      <input type="text" placeholder="The new playlist's name" v-model="newPlaylist" @keyup.enter="addNewPlaylist">
      <button type="button" @click="addNewPlaylist">Add new playlist</button>
    </div>
    <table>
      <tr>
        <th>#</th>
        <th>Playlist's name</th>
      </tr>
      <tr v-for="(name, position) in allPlaylists" :key="position">
        <td>{{ position + 1 }}</td>
        <td><router-link :to="{ name: 'EditPlaylistView', params: {'playlistName': name } }">{{ name }}</router-link></td>
        <td><button type="button" @click="clearAndLoadPlaylist(name)">Replace the playlist</button></td>
        <td><button type="button" @click="loadPlaylist(name)">Append to playlist</button></td>
        <td><button type="button" @click="removePlaylist(name)">Remove</button></td>
      </tr>
    </table>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'PlaylistView',
  data () {
    return {
      newPlaylist: ''
    }
  },
  computed: {
    ...mapGetters({
      allPlaylists: 'getAllPlaylists'
    })
  },
  methods: {
    ...mapActions([
      'setState',
      'setAllPlaylists'
    ]),
    addNewPlaylist () {
      this.$resource('v1/savePlaylist').save({playlistName: this.newPlaylist}).then((response) => {
        this.newPlaylist = ''
        this.loadAllPlaylists()
      })
    },
    clearAndLoadPlaylist (name) {
      this.$resource('v1/clearCurrentPlaylist').get().then((response) => {
        this.loadPlaylist(name)
      })
    },
    loadAllPlaylists () {
      this.$resource('v1/allPlaylists').get().then((response) => {
        this.setAllPlaylists(response.data)
      })
    },
    loadPlaylist (name) {
      this.$resource('v1/loadPlaylist{/name}').get({name: name}).then((response) => {
        this.$resource('v1/playSong').get().then((response) => {
          this.setState('play')
          this.songPlayed = true
        })
      })
    },
    removePlaylist (name) {
      this.$resource('v1/removePlaylist').save({playlistName: name}).then((response) => {
        this.loadAllPlaylists()
      })
    }
  },
  mounted () {
    this.loadAllPlaylists()
  }
}
</script>

<style lang="scss" scoped>
</style>
