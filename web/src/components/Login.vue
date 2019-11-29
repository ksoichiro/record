<template>
  <div class="login">
    <form @submit.prevent="login">
      <input type="text" id="name" v-model="name" placeholder="Username" />
      <input type="password" id="password" v-model="password" placeholder="Password" />
      <button type="submit">Login</button>
    </form>
    <div><router-link to="/signup">Sign up</router-link></div>
  </div>
</template>

<script>
import axios from 'axios'
import auth from '../auth'

export default {
  name: 'Login',
  data () {
    return {
      name: null,
      password: null,
    }
  },
  methods: {
    login: function () {
      axios
        .post('/api/user/login', {
          name: this.name,
          password: this.password,
        })
        .then(response => {
          var token = response.data.token
          if (token && token !== '') {
            auth.login(token)
            this.$router.push({ path: '/' })
          }
        })
    }
  }
}
</script>
