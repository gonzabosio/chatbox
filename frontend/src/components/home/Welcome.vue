<script setup>
const props = defineProps({
    username: String
})
</script>

<template>
    <div>
        <h1>Hello, {{ props.username }}!</h1>
        <p>Enter one of your chats or add a new one to talk with whoever you want</p>
    </div>
</template>

<style scoped>
div {
    margin-left: 16px;
}

p {
    color: darkgray;
}
</style>