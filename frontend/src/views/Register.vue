<script setup>
import { RouterLink, RouterView } from 'vue-router';
import router from '../router';

function goHome() {
    router.push('/home')
}
</script>

<template>
    <fieldset>
        <legend>Register</legend>
        <form @submit.prevent="goHome">
            <label><input type="text"/>Username</label><br />
            <label><input type="password"/>Password</label><br />
            <button type="submit">Sign up</button><br/>
            <RouterLink to='/login'>I have an account</RouterLink>
        </form>
    </fieldset>
    <RouterView />
</template>

<style scoped></style>