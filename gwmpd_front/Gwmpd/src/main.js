import Vue from 'vue'
import VueResource from 'vue-resource'
import App from './App'
import router from './router'

Vue.use(VueResource)
Vue.http.options.root = 'http://127.0.0.1:8060'

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
