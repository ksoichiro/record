const auth = {
  loggedIn: false,
  login: function (token) {
    localStorage.token = token
    this.loggedIn = true
  },
  logout: function () {
    this.loggedIn = false
    delete localStorage.token
  }
}

export default auth
