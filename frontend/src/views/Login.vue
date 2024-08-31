<script setup>
import { RouterLink } from 'vue-router';
import router from '../router';
import axios from 'axios';
let username = ''
let password = ''

async function userLogin() {
  await axios({
    method: 'post',
    url: 'http://localhost:8000/signin',
    data: {
      name: username,
      password: password
    }
  }).then(res => {
    console.log('Message: ' + res.data.message)
    console.log('Status: ' + res.status)
    console.log('StatusText: ' + res.statusText)
    localStorage.setItem('user-id', res.data.user.id)
    localStorage.setItem('access-token', res.data.access_token)
    localStorage.setItem('session-id', res.data.session_id)
    router.push('/home')
  }).catch(err => {
    console.log('Login failed: ' + err.response.data.message)
    console.log('Error details: ' + err.response.data.error)
  })
}
</script>

<template>
  <fieldset>
    <legend>Login</legend>
    <form @submit.prevent="userLogin">
      <label><input type="text" v-model="username" required />Username</label><br />
      <label><input type="password" v-model="password" required />Password</label><br />
      <button type="submit" onclick="">Sign in</button><br>
      <RouterLink to='/register'>Create an account</RouterLink>
    </form>
  </fieldset>
</template>

<style scoped></style>
