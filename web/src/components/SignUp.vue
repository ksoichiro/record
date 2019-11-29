<template>
  <div class="signup">
    <div v-if="error">
      <span class="error">{{ error }}</span>
    </div>
    <form @submit.prevent="signup">
        <input type="text" id="name" v-model="name" placeholder="Username" />
        <input type="password" id="password" v-model="password" placeholder="Password" />
        <button type="submit">Sign up</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'SignUp',
  data () {
    return {
      name: null,
      password: null,
      error: null,
    }
  },
  methods: {
    signup: function() {
      axios
        .post('/api/user/create', {
          name: this.name,
          password: this.password,
        })
        .then(response => {
          if (response.status === 200) {
            this.$router.push({ path: '/' })
          } else {
            this.error = `Failed to sign up: ${response.data.error}`
          }
        })
        .catch(err => {
            this.error = `Failed to sign up: ${err}`
        })
    }
  }
}
</script>

<style scoped>
.error {
  color: #f00;
}
</style>
