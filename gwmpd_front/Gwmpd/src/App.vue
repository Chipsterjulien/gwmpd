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
      <button type="button" name="back" @click="back">Arrière</button>
      <button type="button" name="stop" @click="stop">Stop</button>
      <button type="button" name="play" @click="play" v-if="!songPlayed">Jouer</button>
      <button type="button" name="pause" @click="pause" v-if="songPlayed">Pause</button>
      <button type="button" name="forward" @click="forward">Avant</button>
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
      songPlayed: false,
      user: {}
    }
  },
  methods: {
    back () {
      console.log('back')
    },
    stop () {
      console.log('stop')
    },
    play () {
      console.log('play')
    },
    pause () {
      console.log('pause')
    },
    forward () {
      console.log('forward')
    }
  },
  mounted () {
    this.$stateMPD = this.$resource('v1/stateMPD')
    this.$stateMPD.query().then((response) => {
      this.user = response.data
      this.connected = true
      if (this.user.state === 'play') {
        this.songPlayed = true
      } else {
        this.songPlayed = false
      }
    }, (response) => {
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
