<template lang="html">
  <div class="">
    <b-container fluid>
      <b-form @submit.prevent="login">

        <b-form-group>
          <b-input-group prepend="API's url">
            <b-form-input v-model="url" placeholder="https://your_domain.com:8060"></b-form-input>
          </b-input-group>
        </b-form-group>

        <b-form-group>
          <b-input-group prepend="Login">
            <b-form-input v-model="user"></b-form-input>
          </b-input-group>
        </b-form-group>

        <b-form-group>
          <b-input-group prepend="Password">
            <b-form-input v-model="password" :type="seePassword"></b-form-input>
            <b-input-group-append>
              <b-button :class="visibilityIcon" @click="toggleVisibility"></b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>

        <b-alert :show="dismissCountDown"
          dismissible
          variant="danger"
          @dismissed="dismissCountDown=0"
          @dismiss-count-down="countDownChanged">
            <p>Unable to connect to "<strong>{{ url }}</strong>"</p>
          <b-progress variant="danger"
            :max="dismissSecs"
            :value="dismissCountDown"
            height="4px">
          </b-progress>
        </b-alert>

        <b-button type="submit" size="sm" variant="primary" class="submitButton" :disabled="authorizedSubmitButton === true">Sign in</b-button>
      </b-form>
    </b-container>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
export default {
  name: 'Login',
  data () {
    return {
      data: {
        fetchUser: false,
        rememberMe: false
      },
      dismissSecs: 10,
      dismissCountDown: 0,
      user: '',
      password: '',
      seePassword: 'password',
      showDismissibleAlert: false,
      url: '',
      visibilityBool: false,
      visibilityIcon: 'icon-visibility'
    }
  },
  computed: {
    authorizedSubmitButton: {
      get: function () {
        if (this.url !== '' && this.user !== '' && this.password !== '') {
          return false
        } else {
          // return true if data is missing
          return true
        }
      }
    }
  },
  methods: {
    ...mapActions([
      'setCurrentView'
    ]),
    countDownChanged (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
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
        error (err) {
          this.showAlert()
          console.log(err)
          delete sessionStorage.token
        }
      })
    },
    showAlert () {
      this.dismissCountDown = this.dismissSecs
    },
    toggleVisibility () {
      if (this.visibilityBool === true) {
        this.visibilityBool = false
        this.visibilityIcon = 'icon-visibility'
        this.seePassword = 'password'
      } else {
        this.visibilityBool = true
        this.visibilityIcon = 'icon-visibility_off'
        this.seePassword = 'text'
      }
    }
  },
  mounted () {
    this.setCurrentView('LoginView')
    if (localStorage.url !== 'undefined') {
      this.url = localStorage.url
      this.axios.defaults.baseURL = this.url
    }
  }
}
</script>

<style lang="scss">
  .submitButton {
    margin-top: 10px;
    float: right;
  }
</style>
