<template>
  <div>logout</div>
</template>

<script>
import auth from '../auth'

export default {
  name: 'Logout',
  created: function() {
    auth.logout()
    this.$router.push({ path: '/' }).catch(() => {
      // This works, but promise rejection occurs without any reasons.
      // https://stackoverflow.com/questions/57493516/uncaught-in-promise-undefined-vue-router
    })
  }
}
</script>
