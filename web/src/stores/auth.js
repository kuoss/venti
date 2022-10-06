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
            const response = (await axios.post('user/login', credentials)).data
            if (response) {
                const token = `Bearer ${response.token}`
                localStorage.setItem('token', token)
                localStorage.setItem('userID', response.userID)
                localStorage.setItem('username', response.username)
                axios.defaults.headers.common['Authorization'] = token
                axios.defaults.headers.common['UserID'] = response.userID
                axios.defaults.headers.common['Username'] = response.username
                this.userID = response.userID
                this.username = response.username
                this.loggedIn = true
                // console.log('logged in')
                // await this.fetchUser()
            }
        },
        async logout() {
            const response = (await axios.post('user/logout')).data
            if (response) {
                localStorage.removeItem('token')
                localStorage.removeItem('userID')
                localStorage.removeItem('username')
                this.$reset()
            }
        },
        // async fetchUser() {
        //     this.user = (await axios.get('api/me')).data
        //     this.loggedIn = true
        // },
    },
})