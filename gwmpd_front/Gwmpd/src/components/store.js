import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const state = {
  dataPlayer: {}
}

const getters = {
  getPlayerInfos: state => state.dataPlayer
}

const mutations = {
  CHANGE_ALL_DATA: (state, data) => {
    state.dataPlayer = data
  }
}

// Les actions sont comme les mutations mais elles permettent
// d'appeler plusieurs mutations à la fois et sourtout
// surtout d'utiliser les promises ($resource)
const actions = {
}

export default new Vuex.Store({
  state: state,
  mutations: mutations,
  getters: getters,
  actions: actions,
  strict: true
})
