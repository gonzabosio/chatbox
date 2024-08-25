<script setup>
import { RouterLink, RouterView } from 'vue-router'

</script>

<template>
      <header>
    <nav>
        <RouterLink :to="{ path: '/', replace: true}"><img src="/chat-box.svg" alt="toRoot"></RouterLink>
        <RouterLink to="/login">Login</RouterLink>
        <RouterLink to="/register">Sign Up</RouterLink>
    </nav>
  </header>
  <RouterView/>
</template>

<style scoped>
a { margin-right: 20px; }
</style>