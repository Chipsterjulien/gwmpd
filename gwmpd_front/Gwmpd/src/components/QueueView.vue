<template>
  <div class="">
    <div class="">
      Title song: {{ getPlayerInfos.title }} <br>
      Album: {{ getPlayerInfos.album }} <br>
      Groupe: {{ getPlayerInfos.artist }} <br>
      Consommé: {{ getPlayerInfos.elapsed }}s <br>
      Temps total: {{ getPlayerInfos.timeSong }}s <br>
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
    }
  },
  computed: {
    ...mapGetters([
      'getPlayerInfos'
    ])
  },
  methods: {
    ...mapActions([
      'changePlaylist'
    ])
  },
  mounted () {
    this.$resource('v1/getPlaylist').get().then((response) => {
      this.changePlaylist(response.data)
    })
  }
}
</script>

<style lang="scss" scoped>
</style>
