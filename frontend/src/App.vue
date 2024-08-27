<script setup>
import { RouterView, useRoute } from 'vue-router';
import NavBar from './components/NavBar.vue';
import router from './router'
import { computed } from 'vue';

const route = useRoute()
const showNavBar = computed(() => route.path !== '/home' && route.path !== '/:pathMatch(.*)*')

</script>

<template>
  <NavBar v-if="showNavBar"/>
  <RouterView/>
</template>

<style></style>
