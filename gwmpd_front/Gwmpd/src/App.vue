<template>
  <div id="app">
    <div class="TopBar">
      {{ appName }}
      <router-link :to="{ name: 'QueueView', params: {} }">Queue</router-link>
      <router-link :to="{ name: 'PlaylistView', params: {} }">Playlist</router-link>
      <router-link :to="{ name: 'AboutView', params: {} }">About</router-link>
    </div>
    <div class="">
      <div class="" v-if="connected">
        Connected
      </div>
      <div class="" v-else>
        Disconnected
      </div>
    </div>
    <div class="">
      <button type="button" name="back">Arrière</button>
      <button type="button" name="stop">Stop</button>
      <button type="button" name="play">Jouer</button>
      <button type="button" name="forward">Avant</button>
    </div>
    <div class="">
      <router-view/>
      <router-view name="SideBar"/>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data () {
    return {
      appName: 'Gwmpd',
      stat: [],
      connected: false,
      user: {}
    }
  },
  methods: {
  },
  mounted () {
    this.$stateMPD = this.$resource('v1/stateMPD')
    this.$stateMPD.query().then((response) => {
      this.user = response.data
      this.connected = true
      // console.log('Nom: ', this.user.Name)

      // console.log("J'ai reçu une réponse", response)
    }, (response) => {
      // console.log("Je n'ai rien reçu d'intéressant", response)
      this.connected = false
    })
  },
  beforeUpdate () {
  }
}
</script>

<style lang="scss">
#app {
  /* font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px; */
}
</style>
