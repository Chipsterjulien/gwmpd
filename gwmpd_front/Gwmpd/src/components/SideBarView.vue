<template lang="html">
  <div class="">
    <br>
    <input type="checkbox" :checked="status.consume" @change="toggleConsume">Consume<br>
    <input type="checkbox" :checked="status.random" @change="toggleRandom">Random<br>
    <input type="checkbox" :checked="status.repeat" @change="toggleRepeat">Repeat<br>
    <input type="checkbox" :checked="status.single" @change="toggleSingle">Single<br>
    <br>
    <button type="button" @click="shuffle">Shuffle</button><br>
    <button type="button" @click="updateDB">Update DB</button><br>
    <button type="button" @click="clearQueue">Clear queue</button>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'SideBar',
  data () {
    return {
    }
  },
  computed: {
    ...mapGetters({
      status: 'getStatusInfos'
    })
  },
  methods: {
    ...mapActions([
      'setConsume',
      'setRandom',
      'setRepeat',
      'setSingle',
      'setAllStatus',
      'setPlaylist',
      'setRepeat'
    ]),
    clearQueue () {
      this.$resource('v1/clearCurrentPlaylist').get().then((response) => {
        this.setPlaylist({})
      })
    },
    shuffle () {
      this.$resource('v1/shuffle').get().then((response) => {
        this.$resource('v1/currentPlaylist').get().then((response) => {
          this.setPlaylist(response.data)
        })
      })
    },
    toggleConsume () {
      this.$resource('v1/toggleConsume').update().then((response) => {
        this.setConsume(response.data.consume)
      })
    },
    toggleRandom () {
      this.$resource('v1/toggleRandom').update().then((response) => {
        this.setRandom(response.data.random)
      })
    },
    toggleRepeat () {
      this.$resource('v1/toggleRepeat').update().then((response) => {
        this.setRepeat(response.data.repeat)
      })
    },
    toggleSingle () {
      this.$resource('v1/toggleSingle').update().then((response) => {
        console.log(response.data)
        this.setSingle(response.data.single)
      })
    },
    updateDB () {
      this.$resource('v1/updateDB').get()
    }
  },
  mounted () {
    this.$resource('v1/statusMPD').get().then((response) => {
      this.setAllStatus(response.data)
    })
  }
}
</script>

<style lang="scss" scoped>
</style>
