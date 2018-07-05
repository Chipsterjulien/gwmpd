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
      this.axios.get('v1/clearCurrentPlaylist')
        .then(response => {
          this.setPlaylist({})
        })
    },
    shuffle () {
      this.axios.get('v1/shuffle')
        .then(response => {
          this.axios.get('v1/currentPlaylist')
            .then(response => {
              this.setPlaylist(response.data)
            })
        })
    },
    toggleConsume () {
      this.axios.put('v1/toggleConsume')
        .then(response => {
          this.setConsume(response.data.consume)
        })
    },
    toggleRandom () {
      this.axios.put('v1/toggleRandom')
        .then(response => {
          this.setRandom(response.data.random)
        })
    },
    toggleRepeat () {
      this.axios.put('v1/toggleRepeat')
        .then(response => {
          this.setRepeat(response.data.repeat)
        })
    },
    toggleSingle () {
      this.axios.put('v1/toggleSingle')
        .then(response => {
          this.setSingle(response.data.single)
        })
    },
    updateDB () {
      this.axios.get('v1/updateDB')
    }
  },
  mounted () {
    this.axios.get('v1/statusMPD')
      .then(response => {
        this.setAllStatus(response.data)
      })
  }
}
</script>

<style lang="scss" scoped>
</style>
