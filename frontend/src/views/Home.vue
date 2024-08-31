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
        <div id="navigation" v-if="userData">
            <NavBar :username="userData.name" @chatSelected="setChat" />
        </div>
        <Chat v-if="chatID" :chatName="chatName" :chatId="chatID" :userId="userData.id"/>
        <div id="chat" v-else="userData">
            <Welcome :username="userData.name" :userId="userId" />
        </div>
    </div>
</template>

<style scoped>
#all {
    padding: 16px;
}

#navigation {
    background-color: rgb(53, 53, 53);
}

#aligner {
    display: flex;
}
</style>