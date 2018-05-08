<template lang="html">
  <div class="">
    <br>
    <br>
    Edit playlist View
    <br>
    <table>
      <tr>
        <th>#</th>
        <th>Song's name</th>
        <th>Duration</th>
      </tr>
      <tr>
        <td></td>
      </tr>
    </table>
  </div>
</template>

<script>
export default {
  name: 'EditPlaylistView',
  data () {
    return {
    }
  },
  computed: {
  },
  methods: {
  },
  mounted () {
    // console.log(this.$route.params.playlistName)
    this.$resource('v1/playlistSongsList{/playlistName}').get({playlistName: this.$route.params.playlistName}).then((response) => {
    })
    // console.log(this.$route.params)
  }
}
</script>

<style lang="scss" scoped>
</style>
