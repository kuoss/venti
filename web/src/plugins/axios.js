import { useErrorStore } from '../stores/error';
import router from '../router';
import axios from 'axios';

axios.defaults.baseURL = import.meta.env.AXIOS_BASE_URL
axios.defaults.headers.common['Authorization'] = localStorage.getItem('token')
axios.defaults.headers.common['UserID'] = localStorage.getItem('userID')
// axios.defaults.withCredentials = true

axios.interceptors.request.use(
    function (config) {
        useErrorStore().$reset()
        return config
    },
    function (error) {
        return Promise.reject(error)
    }
)

axios.interceptors.response.use(
    function (response) {
        return response
    },
    function (error) {
        switch (error.response.status) {
            case 401:
                console.log(401)
                localStorage.removeItem('token')
                localStorage.removeItem('userID')
                window.location.reload()
                break
            case 403:
            case 404:
                router.push({
                    name: 'error',
                    props: {
                        error: {
                            message: error.response.data.message,
                            status: error.status,
                        },
                    },
                })
                break
            case 422:
                useErrorStore().$state = error.response.data
                break
            default:
                console.log(error.response.data)
        }
        return Promise.reject(error)
    }
)

export default axios