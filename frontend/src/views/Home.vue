<script setup>
import { onMounted, ref } from 'vue';
import NavBar from '../components/home/ChatsNavBar.vue'
import Welcome from '../components/home/Welcome.vue';
import Chat from '../components/home/Chat.vue'
import { getUserInfo } from '../axios-func/calls';
import Profile from '../components/home/Profile.vue';

const userData = ref({})

onMounted(async () => {
    userData.value = await getUserInfo()
})

const chatID = ref('')
const chatName = ref('')
const setChat = (id, name) => {
    profile.value = false
    chatID.value = id
    chatName.value = name
}
const profile = ref(false)
const openProfile = () => {
    profile.value = true
}
</script>

<template>
    <div id="aligner">
        <NavBar :username="userData.name" @chatSelected="setChat" @goProfile="openProfile" />
        <Profile v-if="profile && userData" :userData="userData.personal" />
        <Chat v-else-if="chatID && userData" :chatName="chatName" :chatId="chatID" :userId="userData.id" />
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