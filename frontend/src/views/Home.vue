<script setup>
import router from '../router';

function logout() {
    router.replace('/')
}
</script>

<template>
    <div>
        <h1>Home</h1>
        <button @click="logout">Logout</button>
    </div>
</template>

<style scoped></style>