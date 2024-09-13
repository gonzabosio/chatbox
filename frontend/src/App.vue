<script setup>
import { RouterView, useRoute } from 'vue-router';
import NavBar from './components/NavBar.vue';
import { computed } from 'vue';

const route = useRoute()
const showNavBar = computed(() => route.name !== 'home' && route.name !== 'not-found')

</script>

<template>
  <NavBar v-if="showNavBar" />
  <RouterView />
</template>

<style></style>
