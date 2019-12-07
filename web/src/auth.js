const auth = {
  loggedIn: false,
  login: function (token) {
    localStorage.setItem('token', token)
    this.loggedIn = true
  },
  logout: function () {
    this.loggedIn = false
    localStorage.removeItem('token')
  }
}

export default auth
