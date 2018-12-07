import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

function convertSecondsToString (time) {
  // Hours, minutes and seconds
  var hrs = ~~(time / 3600)
  var mins = ~~((time % 3600) / 60)
  var secs = ~~time % 60

  // Output like "1:01" or "4:03:59" or "123:03:59"
  var ret = ''

  if (hrs > 0) {
    ret += '' + hrs + ':' + (mins < 10 ? '0' : '')
  }

  ret += '' + mins + ':' + (secs < 10 ? '0' : '')
  ret += '' + secs

  return ret
}

function getMusicName (filename) {
  var filenameSplitted = filename.split('/')

  return filenameSplitted[filenameSplitted.length - 1]
}

const currentPlaylist = {
  state: {
    playlist: {}
  },
  getters: {
    getCurrentPlaylist: state => state.playlist
  },
  mutations: {
    SET_PLAYLIST: (state, newPlaylist) => {
      if (newPlaylist !== {}) {
        var i
        for (i = 0; i < newPlaylist.length; i++) {
          newPlaylist[i].Duration = convertSecondsToString(newPlaylist[i].Duration)
          newPlaylist[i].File = getMusicName(newPlaylist[i].File)
        }
      }
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
      var table = []
      var i
      for (i = 0; i < newPlaylists.length; i++) {
        table.push({'name': newPlaylists[i]})
      }
      state.allPlaylists = table
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

const connectionStatus = {
  state: {
    status: {}
  },
  getters: {
    getConnectionStatus: state => state.status
  },
  mutations: {
    SET_CONNECTION_STATUS: (state, newStatus) => {
      state.status = newStatus
    }
  },
  actions: {
    setConnectionStatus: ({state, commit}, newStatus) => {
      commit('SET_CONNECTION_STATUS', newStatus)
    }
  }
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
    SET_POSITION: (state, newPosition) => {
      state.status.elapsed = Math.round(newPosition * 100) / 100
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
    setPosition: ({state, commit}, newPosition) => {
      commit('SET_POSITION', newPosition)
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
    // stat: stat,
    connectionStatus: connectionStatus,
    status: status
  },
  strict: true
})
