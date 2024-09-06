<script setup>
import { RouterLink } from 'vue-router'
</script>

<template>
  <header>
    <nav>
      <div>
        <RouterLink :to="{ path: '/', replace: true }"><img src="../../public/box_icon.png" alt="to-landing">
        </RouterLink>
        <RouterLink to="/login" class="nav-item">Login</RouterLink>
        <RouterLink to="/register" class="nav-item">Sign Up</RouterLink>
      </div>
    </nav>
  </header>
</template>

<style scoped>
nav {
  background-color: #2e2e2e;
  padding: 10px;
}

nav div {
  display: flex;
  align-items: center;
  margin-left: 8px;
}

.nav-item {
  display: flex;
  margin-left: 20px;
  color: rgb(230, 230, 230);
  font-size: 18px;
  width: 80px;
  height: 40px;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  border: 2px solid transparent;
}

.nav-item:hover {
  border: 2px solid whitesmoke;
  border-radius: 8px;
}
</style>