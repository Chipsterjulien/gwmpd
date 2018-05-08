<template lang="html">
  <div class="">
    <br>
    <br>
    <button type="button">Clear</button>
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
        <td><button type="button">Top</button></td>
        <td><button type="button">Up</button></td>
        <td><button type="button">Down</button></td>
        <td><button type="button">Bottom</button></td>
        <td><button type="button">Remove</button></td>
      </tr>
    </table>
  </div>
</template>

<script>
export default {
  name: 'EditPlaylistView',
  data () {
    return {
      playlist: [],
      playlistName: ''
    }
  },
  computed: {
  },
  methods: {
  },
  mounted () {
    this.playlistName = this.$route.params.playlistName
    this.$resource('v1/playlistSongsList{/playlistName}').get({playlistName: this.playlistName}).then((response) => {
      this.playlist = response.data
      console.log(this.playlist)
    })
  }
}
</script>

<style lang="scss" scoped>
</style>
