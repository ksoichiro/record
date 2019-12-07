import Vue from 'vue'
import VueRouter from 'vue-router'
import Ping from '~/src/components/Ping.vue'
import Login from '~/src/components/Login.vue'
import Home from '~/src/components/Home.vue'
import Logout from '~/src/components/Logout.vue'
import SignUp from '~/src/components/SignUp.vue'
import TaskList from '~/src/components/TaskList.vue'
import auth from '~/src/auth'

Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/ping', component: Ping },
    { path: '/login', component: Login },
    { path: '/', component: Home, meta: { requiresAuth: true }},
    { path: '/logout', component: Logout },
    { path: '/signup', component: SignUp },
    { path: '/task', component: TaskList, meta: { requiresAuth: true }},
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
