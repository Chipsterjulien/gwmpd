<template lang="html">
  <div class="" v-if="getConnectionStatus === true">
    <b-form-group>
      <b-form-checkbox-group buttons button-variant="primary" v-model="selected">
        <b-form-checkbox value="consume" @change="toggleConsume" class="icon-whatshot"></b-form-checkbox>
        <b-form-checkbox value="random" @change="toggleRandom" class="icon-shuffle"></b-form-checkbox>
        <b-form-checkbox value="repeat" @change="toggleRepeat" class="icon-repeat"></b-form-checkbox>
      </b-form-checkbox-group>
    </b-form-group>

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
      selected: []
    }
  },
  computed: {
    ...mapGetters({
      getConnectionStatus: 'getConnectionStatus',
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
    updateDB () {
      this.axios.get('v1/updateDB')
    }
  },
  mounted () {
    this.axios.get('v1/statusMPD')
      .then(response => {
        this.setAllStatus(response.data)
        if (response.data.consume === true) {
          this.selected.push('consume')
        }
        if (response.data.random === true) {
          this.selected.push('random')
        }
        if (response.data.repeat === true) {
          this.selected.push('repeat')
        }
      })
  }
}
</script>

<style lang="scss" scoped>
</style>
