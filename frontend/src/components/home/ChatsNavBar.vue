<script setup>
import { deleteChat, loadChats, logout } from '../../axios-func/calls'
import { computed, onMounted } from 'vue';
import { ref } from 'vue';
import AddChatDialog from '../home/AddChatDialog.vue'

const props = defineProps({
    username: String,
})

const chats = ref([])
onMounted(async () => {
    chats.value = await loadChats()
})

const filteredChats = computed(() => {
    return Array.isArray(chats.value) ? chats.value.map(chat => ({
        chat, participants: chat.participants.filter(contact => contact.name !== props.username)
    })) : [];
});

const refreshChatList = (newList) => {
    chats.value = newList
}

const emit = defineEmits(['chatSelected'])

const setId = (id, chatname) => {
    emit('chatSelected', id, chatname)
}
</script>

<template>
    <div id="content">
        <div id="top">
            <p>{{ props.username }}</p>
            <button @click="logout">Logout</button>
        </div>
        <hr>
        <div id="chat-header">
            <h3>Chats</h3>
            <AddChatDialog :name="props.username" @chatsUpdated="refreshChatList" />
        </div>
        <div v-if="filteredChats">
            <div v-for="(item, i) in filteredChats" :key="i">
                <div v-for="(contact, index) in item.participants" :key="index" id="chat-card"
                    @click="setId(item.chat.id, contact.name)">
                    <div>
                        <p>{{ contact.name }}</p>
                    </div>
                    <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-trash-filled" width="22"
                        height="22" viewBox="0 0 24 24" stroke-width="1.5" stroke="#9e9e9e" fill="none"
                        stroke-linecap="round" stroke-linejoin="round" id="delete-button"
                        @click.stop="console.log('Delete: ' + contact.name), deleteChat(item.chat.id)">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                        <path
                            d="M20 6a1 1 0 0 1 .117 1.993l-.117 .007h-.081l-.919 11a3 3 0 0 1 -2.824 2.995l-.176 .005h-8c-1.598 0 -2.904 -1.249 -2.992 -2.75l-.005 -.167l-.923 -11.083h-.08a1 1 0 0 1 -.117 -1.993l.117 -.007h16z"
                            stroke-width="0" fill="currentColor" />
                        <path
                            d="M14 2a2 2 0 0 1 2 2a1 1 0 0 1 -1.993 .117l-.007 -.117h-4l-.007 .117a1 1 0 0 1 -1.993 -.117a2 2 0 0 1 1.85 -1.995l.15 -.005h4z"
                            stroke-width="0" fill="currentColor" />
                    </svg>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
#content {
    padding-left: 16px;
    padding-right: 16px;
    width: 300px;
    height: 100vh;
    overflow-y: hidden;
    box-sizing: border-box;
    border-right: 6px solid rgba(70, 160, 70, 0.493);
}

#chat-header {
    display: flex;
    justify-content: space-around;
    align-items: center;
}

#chat-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-top: 3px solid rgb(43, 43, 43);
    padding: 3px 6px 3px 6px;
}

#delete-button:hover {
    color: rgb(165, 26, 26);
}

#chat-card:hover {
    cursor: pointer;
    background-color: rgb(85, 85, 85);
}

#top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 8px;
}

#top button {
    cursor: pointer;
    height: 30px;
}
</style>