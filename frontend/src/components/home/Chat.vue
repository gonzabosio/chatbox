<script setup>
const props = defineProps({
    chatId: String
})
</script>

<template>
    <div>
        <h1>Chat</h1>
        <p>{{ props.chatId }}</p>
    </div>
</template>

<style scoped></style>