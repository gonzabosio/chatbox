<script setup>
import { RouterLink } from 'vue-router'
</script>

<template>
      <header>
    <nav>
        <div>
          <RouterLink :to="{ path: '/', replace: true}"><img src="/chat-box.svg" alt="toRoot" class="nav-item"></RouterLink>
          <RouterLink to="/login" class="nav-item">Login</RouterLink>
          <RouterLink to="/register" class="nav-item">Sign Up</RouterLink>
        </div>
    </nav>
  </header>
</template>

<style scoped>
nav {
  background-color: rgb(51, 51, 51);
  padding: 10px;
}
nav div {
  display: flex;
  align-items: center;
}
a { margin-right: 20px; }
</style>