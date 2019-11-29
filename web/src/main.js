import Vue from 'vue'
import router from './router'
import App from './App.vue'

Vue.config.productionTip = false

Vue.config.errorHandler = (err, vm, info) => {
  console.log(`error in Vue.config.errorHandler: ${info}`, err)
}
window.addEventListener("error", event => {
  console.log('error in EventListener', event.error)
})
window.addEventListener('unhandledrejection', event => {
  console.log('unhandledrejection EventListener', event.reason)
})

new Vue({
  router,
  el: '#app',
  render: h => h(App)
}).$mount('#app')
