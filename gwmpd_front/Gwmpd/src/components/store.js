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
  },
  CHANGE_VOLUME: (state, volume) => {
    state.dataPlayer.volume = volume
  }
}

const actions = {
  changeAllData: (store, data) => {
    store.commit('CHANGE_ALL_DATA', data)
  },
  changeVolume: (store, volume) => {
    store.commit('CHANGE_VOLUME', volume)
  }
}

export default new Vuex.Store({
  state: state,
  mutations: mutations,
  getters: getters,
  actions: actions,
  strict: true
})
