<template lang="html">
  <div class="">
    <form class="" @submit.prevent="login">
      <h2>Please sign in</h2>
      <label for="urlAPI">API's url</label>
      <input v-model="url" type="text" placeholder="http://localhost:8060" id="urlAPI" required>
      <label for="username">Login</label>
      <input v-model="user" type="text" id="username" required autofocus>
      <label for="userPassword">Password</label>
      <input v-model="password" type="password" id="userPassword" required>
      <button type="submit">Sign in</button>
    </form>
  </div>
</template>

<script>
export default {
  name: 'Login',
  data () {
    return {
      user: '',
      password: '',
      url: '',
      data: {
        rememberMe: false,
        fetchUser: false
      }
    }
  },
  methods: {
    login () {
      this.axios.defaults.baseURL = this.url
      var redirect = this.$auth.redirect()

      this.$auth.login({
        data: {
          username: this.user,
          password: this.password
        },
        redirect: {name: redirect ? redirect.from.name : 'QueueView'},
        rememberMe: this.data.rememberMe,
        fetchUser: this.data.fetchUser,
        success (response) {
          this.$interval = 5
          this.$auth.token(null, response.data.token)
          sessionStorage.token = response.data.token
          localStorage.url = this.url
        },
        error () {
          delete sessionStorage.token
        }
      })
    }
  },
  mounted () {
    console.log('mounted Login.vue')

    if (localStorage.url !== 'undefined') {
      this.url = localStorage.url
      this.axios.defaults.baseURL = this.url
    }
  }
}
</script>

<style lang="scss">
</style>
