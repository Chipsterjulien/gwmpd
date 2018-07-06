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
      state.playlist = {}
      state.playlist = newPlaylist
    }
  },
  actions: {
    setPlaylist: ({state, commit}, newPlaylist) => {
      commit('SET_PLAYLIST', newPlaylist)
    }
  }
}

const allPlaylists = {
  state: {
    allPlaylists: []
  },
  getters: {
    getAllPlaylists: state => state.allPlaylists
  },
  mutations: {
    SET_ALL_PLAYLISTS: (state, newPlaylists) => {
      state.allPlaylists = newPlaylists
    }
  },
  actions: {
    setAllPlaylists: ({state, commit}, newPlaylists) => {
      commit('SET_ALL_PLAYLISTS', newPlaylists)
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
      state.song = {}
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
    SET_CONSUME: (state, newConsume) => {
      state.status.consume = newConsume
    },
    SET_RANDOM: (state, newRandom) => {
      state.status.random = newRandom
    },
    SET_REPEAT: (state, newRepeat) => {
      state.status.repeat = newRepeat
    },
    SET_SINGLE: (state, newSingle) => {
      state.status.single = newSingle
    },
    SET_STATE: (state, newState) => {
      state.status.state = newState
    },
    SET_VOLUME: (state, newVolume) => {
      state.status.volume = newVolume
    }
  },
  actions: {
    setAllStatus: ({state, commit}, newStatus) => {
      commit('SET_ALL_STATUS', newStatus)
    },
    setConsume: ({state, commit}, newConsume) => {
      commit('SET_CONSUME', newConsume)
    },
    setRandom: ({state, commit}, newRandom) => {
      commit('SET_RANDOM', newRandom)
    },
    setRepeat: ({state, commit}, newRepeat) => {
      commit('SET_REPEAT', newRepeat)
    },
    setSingle: ({state, commit}, newSingle) => {
      commit('SET_SINGLE', newSingle)
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
    allPlaylists: allPlaylists,
    currentPlaylist: currentPlaylist,
    currentSong: currentSong,
    stat: stat,
    status: status
  },
  strict: true
})
