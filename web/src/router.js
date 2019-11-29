import Vue from 'vue'
import VueRouter from 'vue-router'
import Ping from './components/Ping.vue'
import Login from './components/Login.vue'
import Home from './components/Home.vue'
import Logout from './components/Logout.vue'
import SignUp from './components/SignUp.vue'
import auth from './auth'

Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/ping', component: Ping },
    { path: '/login', component: Login },
    { path: '/', component: Home, meta: { requiresAuth: true }},
    { path: '/logout', component: Logout },
    { path: '/signup', component: SignUp },
    { path: '*', redirect: '/' }
  ]
})

router.beforeEach(function (to, from, next) {
  if (to.matched.some(function (record) {
    return record.meta.requiresAuth
  }) && !auth.loggedIn) {
    next({ path: '/login', query: { redirect: to.fullPath }})
  } else {
    next()
  }
})

export default router
