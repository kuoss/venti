import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    loggedIn: localStorage.getItem('token') ? true : false,
    userID: localStorage.getItem('userID'),
    username: localStorage.getItem('username'),
  }),
  actions: {
    async login(credentials) {
      let formData = new FormData()
      formData.append('username', credentials.username)
      formData.append('password', credentials.password)

      const response = await fetch('/auth/login', {
        body: formData,
        method: 'post',
      })
      const data = await response.json()
      if (data) {
        const token = `Bearer ${data.token}`
        localStorage.setItem('token', token)
        localStorage.setItem('userID', data.userID)
        localStorage.setItem('username', data.username)
        axios.defaults.headers.common['Authorization'] = token
        axios.defaults.headers.common['UserID'] = data.userID
        axios.defaults.headers.common['Username'] = data.username
        this.userID = data.userID
        this.username = data.username
        this.loggedIn = true
      }
    },
    logout() {
      fetch('/auth/logout', { method: 'post' }).then(() => console.log('logged out'))
      localStorage.removeItem('token')
      localStorage.removeItem('userID')
      localStorage.removeItem('username')
      this.$reset()
    },
  },
})
