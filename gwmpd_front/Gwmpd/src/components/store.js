import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const currentPlaylist = {
  state: {
    playlist: {}
  },
  getters: {
    getCurrentPlaylist: state => state.playlist
  },
  mutations: {
    SET_PLAYLIST: (state, newPlaylist) => {
      state.playlist = newPlaylist
    }
  },
  actions: {
    setPlaylist: ({state, commit}, newPlaylist) => {
      commit('SET_PLAYLIST', newPlaylist)
    }
  }
}

const currentSong = {
  state: {
    song: {}
  },
  getters: {
    getCurrentSongInfos: state => state.song
  },
  mutations: {
    SET_SONG: (state, newSong) => {
      state.song = newSong
    },
    SET_ID: (state, newID) => {
      state.song.Id = newID
    }
  },
  actions: {
    setSong: ({state, commit}, newSong) => {
      commit('SET_SONG', newSong)
    },
    setID: ({state, commit}, newID) => {
      commit('SET_ID', newID)
    }
  }
}

const stat = {
  state: {},
  getters: {},
  mutations: {},
  actions: {}
}

const status = {
  state: {
    status: {}
  },
  getters: {
    getStatusInfos: state => state.status
  },
  mutations: {
    SET_ALL_STATUS: (state, newStatus) => {
      state.status = newStatus
    },
    SET_VOLUME: (state, newVolume) => {
      state.status.volume = newVolume
    },
    SET_STATE: (state, newState) => {
      state.status.state = newState
    }
  },
  actions: {
    setAllStatus: ({state, commit}, newStatus) => {
      commit('SET_ALL_STATUS', newStatus)
    },
    setState: ({state, commit}, newState) => {
      commit('SET_STATE', newState)
    },
    setVolume: ({state, commit}, newVolume) => {
      commit('SET_VOLUME', newVolume)
    }
  }
}

export default new Vuex.Store({
  modules: {
    currentPlaylist: currentPlaylist,
    currentSong: currentSong,
    stat: stat,
    status: status
  },
  strict: true
})
