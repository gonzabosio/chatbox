<script setup>
import { onMounted, ref } from 'vue';
import router from '../router';
import axios from 'axios';

let userData = ref({})
onMounted(() => {
    axios({
        method: 'get',
        url: `http://localhost:8000/user/id`,
        headers: {
            Authorization: ''
        }
    }).then(res => {
        userData.value = res.data.user_data
    }).catch(err => {
        console.log(err.message)
        console.log('Error details:'+err)
    })
})
function logout() {
    router.replace('/')
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