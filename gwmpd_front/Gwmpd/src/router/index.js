import Vue from 'vue'
import Router from 'vue-router'
import AboutView from '@/components/AboutView'
import PlaylistView from '@/components/PlaylistView'
import QueueView from '@/components/QueueView'
import SideBarView from '@/components/SideBarView'
import EditPlaylistView from '@/components/EditPlaylistView'
import LoginView from '@/components/LoginView'
// import GlobalView from '@/components/GlobalView'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'QueueView',
      components: {
        default: QueueView,
        SideBar: SideBarView
      }
    }, {
      path: '/playlist',
      name: 'PlaylistView',
      components: {
        default: PlaylistView
      }
    }, {
      path: '/about',
      name: 'AboutView',
      components: {
        default: AboutView
      }
    }, {
      path: '/editPlaylist/:playlistName',
      name: 'EditPlaylistView',
      components: {
        default: EditPlaylistView
      }
    }, {
      path: '/login',
      name: 'LoginView',
      component: LoginView
    }, {
      path: '*',
      redirect: '/'
    }
  ]
})
