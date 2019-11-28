import Vue from 'vue'
import VueRouter from 'vue-router'
import Ping from './components/Ping.vue'

Vue.use(VueRouter)

export default new VueRouter({
  mode: 'history',
  routes: [
    { path: '/ping', component: Ping }
  ]
})
