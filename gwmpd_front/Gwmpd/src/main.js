import Vue from 'vue'
import VueResource from 'vue-resource'
import App from './App'
import router from './router'
import store from './components/store'
// import store from './store'

Vue.use(VueResource)
Vue.http.options.root = 'http://127.0.0.1:8060'

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
