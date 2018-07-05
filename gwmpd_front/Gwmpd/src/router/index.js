import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/components/Login'
import QueueView from '@/components/QueueView'
import SideBarView from '@/components/SideBarView'
import PlaylistView from '@/components/PlaylistView'
import EditPlaylistView from '@/components/EditPlaylistView'
import AboutView from '@/components/AboutView'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'QueueView',
      meta: {auth: true},
      components: {
        default: QueueView,
        SideBar: SideBarView
      }
    }, {
      path: '/login',
      name: 'Login',
      meta: {auth: false},
      component: Login
    }, {
      path: '/about',
      name: 'AboutView',
      components: {
        default: AboutView
      }
    }, {
      path: '/playlist',
      name: 'PlaylistView',
      meta: {auth: true},
      components: {
        default: PlaylistView
      }
    }, {
      path: '/editPlaylist/:playlistName',
      name: 'EditPlaylistView',
      meta: {auth: true},
      components: {
        default: EditPlaylistView
      }
    }, {
      path: '*',
      redirect: '/'
    }
  ]
})
