<template lang="html">
  <div class="">
    <br>
    <br>

    <button type="button" name="button">New playlist</button>
    <table>
      <tr>
        <th>#</th>
        <th>Playlist's name</th>
      </tr>
      <tr v-for="(name, position) in allPlaylists" :key="position">
        <td>{{ position + 1 }}</td>
        <td>{{ name }}</td>
        <td><button type="button" @click="clearAndLoadPlaylist(name)">Replace the playlist</button></td>
        <td><button type="button" @click="loadPlaylist(name)">Append to playlist</button></td>
        <td><button type="button" @click="loadPlaylist(name)">Remove</button></td>
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
    clearAndLoadPlaylist (name) {
      this.$resource('v1/clearCurrentPlaylist').get().then((response) => {
        this.loadPlaylist(name)
      })
    },
    loadPlaylist (name) {
      this.$resource('v1/loadPlaylist{/name}').get({name: name}).then((response) => {
        this.$resource('v1/playSong').get().then((response) => {
          this.setState('play')
          this.songPlayed = true
        })
      })
    }
  },
  mounted () {
    this.$resource('v1/allPlaylists').get().then((response) => {
      this.setAllPlaylists(response.data)
    })
  }
}
</script>

<style lang="scss" scoped>
</style>
