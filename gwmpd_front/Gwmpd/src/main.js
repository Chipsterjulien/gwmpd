// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueAuth from '@websanova/vue-auth'
import store from './components/store'

Vue.router = router

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
