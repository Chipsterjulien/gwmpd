<template lang="html">
  <div class="">
    Config.vue
    <br>
    <br>
    <form class="" @submit.prevent="connect">
      <label for="urlAPI">API's url</label>
      <input v-model="url" type="text" placeholder="http://localhost:8060" id="urlAPI" text="url" required>
      <label for="autoConnect">Auto connect</label>
      <input type="checkbox" id="autoConnect" v-model="autoconnect">
      <button type="submit">Connect</button>
    </form>
  </div>
</template>

<script>
export default {
  name: 'Config',
  data () {
    return {
      url: '',
      autoconnect: false
    }
  },
  methods: {
    connect () {
      localStorage.url = this.url
      localStorage.autoconnect = this.autoconnect
      this.axios.defaults.baseURL = this.url
      this.$router.go('Login')
    }
  },
  mounted () {
    // if (typeof variable === 'undefined' || variable === null) {}
    if (localStorage.autoconnect !== 'undefined') {
      this.autoconnect = localStorage.autoconnect
    }
    if (localStorage.url !== 'undefined') {
      this.url = localStorage.url
    }

    if (this.autoconnect && (this.url !== 'undefined')) {
      this.axios.defaults.baseURL = this.url
      this.$router.go('Login')
    }
  }
}
</script>

<style lang="css">
</style>
