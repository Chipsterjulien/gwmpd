<template lang="html">
  <div class="" v-if="getConnectionStatus === true">
    <div class="text-center positionning">
      <b-row>
        <b-col>
          <b-form-group>
            <b-form-checkbox-group buttons button-variant="primary" v-model="selected">
              <b-form-checkbox v-b-tooltip.hover.top title="Consume music (2-state button)" value="consume" @change="toggleConsume" class="icon-whatshot iconSize"></b-form-checkbox>
              <b-form-checkbox v-b-tooltip.hover.top title="Random without changing list order (2-state button)" value="random" @change="toggleRandom" class="icon-shuffle iconSize"></b-form-checkbox>
              <b-form-checkbox v-b-tooltip.hover.top title="Loop playlist (2-state button)" value="repeat" @change="toggleRepeat" class="icon-repeat iconSize"></b-form-checkbox>
            </b-form-checkbox-group>
          </b-form-group>
        </b-col>
        <b-col>
          <b-button-group>
            <b-button v-b-tooltip.hover.top title="Random list order" variant="success" @click="shuffle" class="icon-shuffle iconSize"></b-button>
            <b-button v-b-tooltip.hover.top title="Update music database" variant="success" @click="updateDB" class="icon-autorenew iconSize"></b-button>
            <b-button v-b-tooltip.hover.top title="Clear queue" variant="success" @click="clearQueue" class="icon-clear_all iconSize"></b-button>
          </b-button-group>
        </b-col>
      </b-row>
    </div>
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
  .positionning {
    padding-bottom: 5px;
  }
</style>
