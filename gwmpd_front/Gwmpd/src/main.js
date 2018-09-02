// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'

import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import bButton from 'bootstrap-vue/es/components/button/button'
import bAlert from 'bootstrap-vue/es/components/alert/alert'
import bForm from 'bootstrap-vue/es/components/form/form'
import bFormGroup from 'bootstrap-vue/es/components/form-group/form-group'
import bFormInput from 'bootstrap-vue/es/components/form-input/form-input'
import bContainer from 'bootstrap-vue/es/components/layout/container'
import bImg from 'bootstrap-vue/es/components/image/img'
import bNavbar from 'bootstrap-vue/es/components/navbar/navbar'
import bNavbarBrand from 'bootstrap-vue/es/components/navbar/navbar-brand'

import App from './App'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueAuth from '@websanova/vue-auth'
import store from './components/store'

Vue.router = router

Vue.component('b-button', bButton)
Vue.component('b-alert', bAlert)
Vue.component('b-form', bForm)
Vue.component('b-form-group', bFormGroup)
Vue.component('b-form-input', bFormInput)
Vue.component('b-container', bContainer)
Vue.component('b-img', bImg)
Vue.component('b-navbar', bNavbar)
Vue.component('b-navbar-brand', bNavbarBrand)

Vue.use(VueAxios, axios)
Vue.axios.defaults.baseURL = 'http://localhost:8060'

Vue.use(VueAuth, {
  authRedirect: 'Login',
  auth: require('@websanova/vue-auth/drivers/auth/bearer.js'),
  http: require('@websanova/vue-auth/drivers/http/axios.1.x.js'),
  router: require('@websanova/vue-auth/drivers/router/vue-router.2.x.js'),
  loginData: {url: 'login', method: 'POST', redirect: '/'},
  refreshData: {url: 'v1/refresh', enabled: false},
  fetchData: {enabled: false}
})

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
