<script setup>
import router from '../../router/index';
import getChats from '../../axios-func/calls'
import { computed, onMounted } from 'vue';
import { ref } from 'vue';
import AddChatDialog from '../home/AddChatDialog.vue'

const props = defineProps({
    username: String,
})

const sessionId = localStorage.getItem('session-id')

const chats = ref([])
onMounted(async () => {
    chats.value = await getChats()
})

const filteredChats = computed(() => {
  return Array.isArray(chats.value) ? chats.value.map(chat => ({
    chat, participants: chat.participants.filter(contact => contact.name !== props.username)
  })) : [];
});

const refreshChatList = (newList) => {
    chats.value = newList
}

const emit = defineEmits(['chatIdSelected'])

const setId = (id) => {
    emit('chatIdSelected', id)
}

function logout() {
    axios({
        method: 'post',
        url: `http://localhost:8000/revoke/${sessionId}`,
    }).then(res => {
        console.log('Refresh token revoked: '+ res.data.message)
        axios({
            method: 'delete',
            url: `http://localhost:8000/logout/${sessionId}`
        }).then(res => {
            console.log(res.data.message)
            localStorage.clear()
            router.replace('/')
        }).catch(err => {
            console.log(err.response.data.message)
            console.log(err.response.data.error)
        })
    }).catch(err => {
        console.log('Refresh token was not revoked: '+err.response.data.message)
    })
}
</script>

<template>
    <div id="content">
        <p>{{ props.username }}</p>
        <button @click="logout">Logout</button>
        <hr>
        <div id="chat-header">
            <h3>Chats</h3>
            <AddChatDialog :name="props.username" @chatsUpdated="refreshChatList"/>
        </div>
        <div v-if="filteredChats">
            <div v-for="(item, i) in filteredChats" :key="i">
                <div v-for="(contact, index) in item.participants" :key="index" id="chat-card"
                    @click="console.log('clicked: '+ contact.name+' '+contact.id+' | chat-id: '+item.chat.id), setId(item.chat.id)"
                ><p>{{ contact.name }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
#content {
    width: 200px;
    height: 100vh;
}
#chat-header {
    display: flex;
    justify-content: space-around;
}
#chat-card {
    padding: 3px 6px 3px 6px;
}
#chat-card:hover {
    cursor: pointer;
    background-color: darkgray;
}
</style>