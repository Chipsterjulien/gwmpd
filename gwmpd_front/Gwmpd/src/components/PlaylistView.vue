<template lang="html">
  <div class="positionning" v-if="getConnectionStatus === true">
    <b-container class="addingPlaylist">
      <b-input-group prepend="New playlist's name">
        <b-form-input v-model="newPlaylist" @keyup.enter.native="addNewPlaylist"></b-form-input>
        <b-input-group-append>
          <b-button variant="info" v-b-tooltip.hover.top title="Create a new playlist" @click="addNewPlaylist">Add</b-button>
        </b-input-group-append>
      </b-input-group>
    </b-container>

    <b-table stacked="md" striped hover :sort-by.sync="sortBy" :sort-desc.sync="sortDesc" :items="allPlaylists" :fields="fields">
      <template slot="removePlaylist" slot-scope="data">
        <div class="alignButtonInTable">
          <b-button v-b-tooltip.hover.top title="Edit playlist" @click="editPlaylist(data.item.name)" class="icon-mode_edit buttonSpaceBetween"></b-button>
          <b-button v-b-tooltip.hover.top title="Remove this playlist" @click="removePlaylist(data.item.name)" class="icon-delete buttonMarginRight"></b-button>
          <b-button v-b-tooltip.hover.top title="Add this to current playlist" @click="loadPlaylist(data.item.name)" class="icon-queue_music buttonSpaceBetween"></b-button>
          <b-button v-b-tooltip.hover.top title="Replace the current playlist by this one" @click="clearAndLoadPlaylist(data.item.name)" class="icon-add"></b-button>
        </div>
      </template>
    </b-table>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'PlaylistView',
  data () {
    return {
      fields: [{key: 'name', sortable: true, label: 'Playlist\'s name'}, {key: 'removePlaylist', label: ''}],
      newPlaylist: '',
      sortBy: 'name',
      sortDesc: false
    }
  },
  computed: {
    ...mapGetters({
      allPlaylists: 'getAllPlaylists',
      getConnectionStatus: 'getConnectionStatus'
    })
  },
  methods: {
    ...mapActions([
      'setCurrentView',
      'setState',
      'setAllPlaylists'
    ]),
    addNewPlaylist () {
      if (this.newPlaylist === '') {
        return
      }
      this.axios.post('v1/savePlaylist', {
        playlistName: this.newPlaylist
      })
        .then(response => {
          this.newPlaylist = ''
          this.loadAllPlaylists()
        })
    },
    clearAndLoadPlaylist (name) {
      this.axios.get('v1/clearCurrentPlaylist')
        .then(response => {
          this.loadPlaylist(name)
        })
    },
    editPlaylist (playlistName) {
      this.$router.push({ name: 'EditPlaylistView', params: {'playlistName': playlistName} })
    },
    loadAllPlaylists () {
      this.axios.get('v1/allPlaylists')
        .then(response => {
          this.setAllPlaylists(response.data)
        })
    },
    loadPlaylist (name) {
      this.axios.get('v1/loadPlaylist', {params: {name: name}})
        .then(response => {
          this.axios.get('v1/playSong')
            .then(response => {
              this.setState('play')
              this.songPlayed = true
            })
        })
    },
    removePlaylist (name) {
      this.axios.post('v1/removePlaylist', {playlistName: name})
        .then(response => {
          this.loadAllPlaylists()
        })
    }
  },
  mounted () {
    if (localStorage.url !== 'undefined') {
      this.url = localStorage.url
      this.axios.defaults.baseURL = this.url
    }

    this.setCurrentView('PlaylistView')
    this.loadAllPlaylists()
  }
}
</script>

<style lang="scss" scoped>
  .addingPlaylist {
    padding-bottom: 5px;
  }

  // J'ai un changement de taille à 396px et ça part en sucette à ce moment là

  .alignButtonInTable {
    display: flex;
    flex-direction: row;
    float: right;
  }

  .buttonSpaceBetween {
    margin-right: 4px;
  }

  .buttonMarginRight {
    margin-right: 12px;
  }
</style>
