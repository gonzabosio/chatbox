import axios from 'axios'

const axiosInstance = axios.create({
    baseURL: 'http://localhost:8000',
    withCredentials: true,
})

let refreshToken = ''
// set the refresh token in cookies
axiosInstance.interceptors.request.use(
    (config) => {
        if (refreshToken) {
            config.headers['Cookie'] = refreshToken
        }
        return config;
    },
    (error) => {
        return Promise.reject(error)
    }
)

function setTokenInCookie(value) {
    console.log(value)
    refreshToken = value
}

export { axiosInstance, setTokenInCookie }