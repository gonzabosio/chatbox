<script setup>
import { onMounted, ref } from 'vue';
import axios from 'axios';
import { axiosInstance } from '../axios-func/axiosInstance';
import NavBar from '../components/home/ChatsNavBar.vue'
import Welcome from '../components/home/Welcome.vue';
import Chat from '../components/home/Chat.vue'

let userData = ref({})

onMounted(() => {
    axios({
        method: 'get',
        url: `http://localhost:8000/user/${localStorage.getItem('user-id')}`,
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
                axios({
                    method: 'get',
                    url: `http://localhost:8000/user/${localStorage.getItem('user-id')}`,
                    headers: {
                        Authorization: 'Bearer '+localStorage.getItem('access-token')
                    },
                }).then(res => {
                    userData.value = res.data.user_data
                }).catch(err => {
                    console.log('getting data after renew token failed: '+err.response.data.error)
                })
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

const chatID = ref('')
const setChatId = (id) => {
    chatID.value = id
}
</script>

<template>
    <div id="aligner">
        <div id="navigation" v-if="userData">
            <NavBar :username="userData.name" @chatIdSelected="setChatId"/>
        </div>
        <Chat v-if="chatID" :chatId="chatID"/>
        <div id="chat" v-else="userData">
            <Welcome  :username="userData.name" :userId="userId"/>
        </div>
    </div>
</template>

<style scoped>
#all {
    padding: 16px;
}
#navigation {
    background-color: rgb(53, 53, 53);
}
#aligner {
    display: flex;
}
</style>