<template>
  <div class="task-list">
    <table>
      <tr>
        <td>{{ task }}</td>
      </tr>
    </table>
    <router-link to="/logout">Logout</router-link>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Home',
  data () {
    return {
      task: null
    }
  },
  methods: {
    listTasks: function() {
      console.error("Authorization: " + localStorage.token);
      axios
        .get('/api/task', {
          headers: {
            "Authorization": localStorage.token,
            "Content-Type": "application/json",
          },
          data: {},
        })
        .then(response => {
          console.log(response.data)
          this.task = response.data
        })
    }
  },
  mounted () {
    this.listTasks()
  }
}
</script>
