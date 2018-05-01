import Vue from 'vue'
import Vuex from 'vuex'

// Si je veux plusieurs stores, il faut que je regarde du coté des modules:
// https://vuex.vuejs.org/fr/modules.html

Vue.use(Vuex)

const state = {
  dataPlayer: {},
  playlist: {}
}

const getters = {
  getPlayerInfos: state => state.dataPlayer,
  getPlaylist: state => state.playlist
}

const mutations = {
  CHANGE_ALL_DATA: (state, data) => {
    state.dataPlayer = data
  },
  CHANGE_VOLUME: (state, volume) => {
    state.dataPlayer.volume = volume
  },
  CHANGE_PLAYLIST: (state, playlist) => {
    state.playlist = playlist
  }
}

const actions = {
  changeAllData: (store, data) => {
    store.commit('CHANGE_ALL_DATA', data)
  },
  changeVolume: (store, volume) => {
    store.commit('CHANGE_VOLUME', volume)
  },
  changePlaylist: (store, playlist) => {
    store.commit('CHANGE_PLAYLIST', playlist)
  }
}

export default new Vuex.Store({
  state: state,
  mutations: mutations,
  getters: getters,
  actions: actions,
  strict: true
})
