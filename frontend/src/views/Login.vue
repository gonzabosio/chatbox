<script setup>
import { RouterLink, RouterView } from 'vue-router';
import router from '../router';

function goHome() {
 router.push('/home')
}
</script>

<template>
  <fieldset>
    <legend>Login</legend>
      <form @submit.prevent="goHome">
          <label><input type="text" />Username</label><br />
          <label><input type="password" />Password</label><br />
            <button type="submit" onclick="">Sign in</button><br>
          <RouterLink to='/register'>Create an account</RouterLink>
      </form>
  </fieldset>
  <RouterView />
</template>

<style scoped></style>
