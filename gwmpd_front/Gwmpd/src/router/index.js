import Vue from 'vue'
import Router from 'vue-router'
import AboutView from '@/components/AboutView'
import PlaylistView from '@/components/PlaylistView'
import QueueView from '@/components/QueueView'
import SideBar from '@/components/SideBar'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'QueueView',
      components: {
        default: QueueView,
        SideBar: SideBar
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
      path: '*',
      redirect: '/'
    }
  ]
})
