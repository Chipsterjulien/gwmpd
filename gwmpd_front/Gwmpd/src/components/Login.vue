<template lang="html">
  <div class="">
    <b-container fluid>
      <b-form @submit.prevent="login">
        <b-form-group id="fieldURLAPI" description="Put API's location" label="API's url" label-for="urlAPI">
          <b-form-input id="urlAPI" v-model="url" placeholder="https://your_domain.com:8060"></b-form-input>
        </b-form-group>
        <b-form-group id="fieldLogin" description="Put your login" label="Login" label-for="username">
          <b-form-input id="username" v-model="user" placeholder="Login"></b-form-input>
        </b-form-group>
        <b-form-group id="fieldPassword" description="Put your password" label="Password" label-for="userPassword">
          <b-form-input id="userPassword" v-model="password" placeholder="Password" type="password"></b-form-input>
        </b-form-group>
        <b-button type="submit" size="sm" variant="primary">Sign in</b-button>
      </b-form>
    </b-container>
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
