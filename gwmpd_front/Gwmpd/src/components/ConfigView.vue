<template lang="html">
  <div class="" v-if="getConnectionStatus === true">
    ConfigView
    * Utiliser le store
    * Pas de variation pour le réglage du son
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'ConfigView',
  data () {
    return {
    }
  },
  computed: {
    ...mapGetters({
      getConnectionStatus: 'getConnectionStatus'
    })
  },
  methods: {
    ...mapActions([
    ])
  },
  mounted () {
  }
}
</script>

<style lang="scss" scoped>
</style>
