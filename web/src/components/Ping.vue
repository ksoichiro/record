<template>
  <div id="app">{{ ping }}</div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Ping',
  el: '#app',
  data () {
    return {
      ping: null
    }
  },
  mounted () {
    axios
      .get('/api/ping')
      .then(response => (this.ping = response))
  }
}
</script>
