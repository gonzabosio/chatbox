<script setup>
import { onMounted, ref } from 'vue';
import router from '../router';
import axios from 'axios';
import { axiosInstance } from '../axiosInstance';

let userData = ref({})
const userId = localStorage.getItem('user-id')
const sessionId = localStorage.getItem('session-id')

onMounted(() => {
    axios({
        method: 'get',
        url: `http://localhost:8000/user/${userId}`,
        headers: {
            Authorization: 'Bearer '+localStorage.getItem('access-token')
        },
    }).then(res => {
        userData.value = res.data.user_data
    }).catch(err => {
        if (err.response.status === 401) {
            console.log('Unauthorized: '+ err.response.data.message)
            axiosInstance({
                method: 'post',
                url: '/token/renew',
            }).then(res => {
                localStorage.setItem('access-token', res.data.access_token)
                console.log(res.data.message)
                userData.value = res.data.user_data
            }).catch(err => {
                if (err.response.data.message === "refresh token is not in cookies") {
                    console.log('Refresh token expired: '+ err.response.data.message)
                    logout()
                } else {
                    console.log('Error renewing token: '+err.response.data.message)
                    console.log('Error detail: '+err.response.data.error)
                }
            })
        } else {
            console.log('Error retrieving user data: '+err.response.data.message)
            console.log('Error detail: '+ err.response.data.error)
        }
    })
})
function logout() {
    axios({
        method: 'post',
        url: `http://localhost:8000/revoke/${sessionId}`,
    }).then(res => {
        console.log('Refresh token revoked: '+ res.data.message)
        axios({
            method: 'post',
            url: `http://localhost:8000/logout/${sessionId}`
        }).then(res => {
            console.log(res.data.message)
            localStorage.clear()
            router.replace('/')
        }).catch(err => {
            console.log(err.response.data.message)
            console.log(err.response.data.error)
        })
    }).catch(err => {
        console.log('Refresh token was not revoked: '+err.response.data.message)
    })
}
</script>
<template>
    <div>
        <h1>Home</h1>
        <ul>
            <li>{{ userData.id }}</li>
            <li>{{ userData.name }}</li>
        </ul>
        <button @click="logout">Logout</button>
    </div>
</template>

<style scoped></style>