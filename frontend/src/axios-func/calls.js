import router from "../router";
import { axiosInstance } from "./axiosInstance";
import axios from "axios";

// get user info
export const getUserInfo = async () => await axiosInstance({
    method: 'get',
    url: `/user/${localStorage.getItem('user-id')}`
}).then(res => {
    console.log('User data retrieved')
    return res.data.user_data
}).catch(err => {
    console.log('Error retrieving user data: ' + err.response.data.message)
    console.log('Error detail: ' + err.response.data.error)
    return {}
})

export const savePersonalData = (pEmail, pAge, pCountry) => axiosInstance({
    method: 'put',
    url: `/user/save-personal/${localStorage.getItem('user-id')}`,
    data: {
        email: pEmail,
        country: pCountry,
        age: pAge
    }
}).then(res => {
    console.log(pAge, pCountry, pEmail)
    console.log('Personal data saved')
    return res.status
}).catch(err => {
    console.log('Failed save personal data')
    return err.response.status
})

// FOR CHAT COMPONENTS
export const loadChats = () => axiosInstance({
    method: 'get',
    url: `/chat/${localStorage.getItem('user-id')}`
}).then(res => {
    console.log(res.data.chats)
    console.log(res.data.message)
    return res.data.chats
}).catch(err => {
    console.log(err.response.data.message)
    console.log(err.response.data.error)
    return []
})

export const addChat = (newChat, name) => axiosInstance({
    method: 'post',
    url: '/chat',
    data: {
        username: newChat,
        petitioner_id: localStorage.getItem('user-id'),
        petitioner: name
    }
}).then(async res => {
    console.log(res.data.message)
    return res.data.chat
}).catch(err => {
    console.log(err.response.data.message)
    throw err
})

export const loadMessages = (chatId) => axiosInstance({
    method: 'get',
    url: `/chat/${chatId}/messages`
}).then(res => {
    console.log(res.data.messages)
    return res.data.messages
}).catch(err => {
    console.log(err.response.data.message)
    console.log(err.response.data.error)
    return []
})

export const deleteChat = (chatId) => axiosInstance({
    method: 'delete',
    url: `/chat/${chatId}`
}).then(res => {
    console.log(res.data.message)
}).catch(err => {
    console.log(err.response.data.error)
    console.log(err.response.data.message)
})

const baseUrl = import.meta.env.VITE_BACK_BASE_URL
// SESSION METHODS
export const renewToken = () => {
    return axios({
        method: 'post',
        url: `${baseUrl}/token/renew/${localStorage.getItem("session-id")}`,
    }).then(res => {
        localStorage.setItem('access-token', res.data.access_token)
        console.log(res.data.message)
        return res.data.access_token
    }).catch(err => {
        console.log(err.response.data.message)
        console.log(err.response.data.error)
        logout()
        return null
    })
}

export const logout = () => axios({
    method: 'post',
    url: `${baseUrl}/token/revoke/${localStorage.getItem('session-id')}`,
}).then(res => {
    if (res.status === 204) {
        console.log('Refresh token revoked')
        axios({
            method: 'delete',
            url: `${baseUrl}/logout/${localStorage.getItem('session-id')}`
        }).then(res => {
            console.log(res.data.message)
            localStorage.clear()
            router.replace('/')
        }).catch(err => {
            console.log(err.message)
            console.log(err.response.data.error)
        })
    } else {
        console.log(err.message)
        console.log(err.response.data.error)
    }
}).catch(err => {
    console.log('Refresh token was not revoked: ' + err.response.data.message)
})