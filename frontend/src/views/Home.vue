<script setup>
import { onMounted, ref } from 'vue';
import NavBar from '../components/home/ChatsNavBar.vue'
import Welcome from '../components/home/Welcome.vue';
import Chat from '../components/home/Chat.vue'
import { getUserInfo } from '../axios-func/calls';

const userData = ref({})

onMounted(async () => {
    userData.value = await getUserInfo()
})

const chatID = ref('')
const chatName = ref('')
const setChat = (id, name) => {
    chatID.value = id
    chatName.value = name
}
</script>

<template>
    <div id="aligner">
        <NavBar :username="userData.name" @chatSelected="setChat" />
        <Chat v-if="chatID && userData" :chatName="chatName" :chatId="chatID" :userId="userData.id" />
        <div id="chat" v-else="userData">
            <Welcome :username="userData.name" :userId="userData.id" />
        </div>
    </div>
</template>

<style scoped>
#all {
    padding: 16px;
}

#aligner {
    display: flex;
}
</style>