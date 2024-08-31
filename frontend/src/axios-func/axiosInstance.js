import axios from 'axios'
import { renewToken } from './calls'

const axiosInstance = axios.create({
    baseURL: 'http://localhost:8000'
})

axiosInstance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access-token')
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

let refreshingToken = null

axiosInstance.interceptors.response.use(
    (response) => response,
    async (error) => {
        if (error.response.status === 401 && !error.config._retry) {
            // Wait to complete refreshing
            if (!refreshingToken) {
                refreshingToken = renewToken().then(() => {
                    refreshingToken = null
                })
            }
            await refreshingToken

            // Retry request
            error.config._retry = true
            const newAccessToken = localStorage.getItem('access-token')
            error.config.headers.Authorization = `Bearer ${newAccessToken}`
            return axiosInstance(error.config)
        }
        return Promise.reject(error)
    }
)

export { axiosInstance }