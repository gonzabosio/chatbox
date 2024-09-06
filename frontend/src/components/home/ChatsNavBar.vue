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

const refreshChatList = (newChat) => {
    if (chats.value === null) {
        chats.value = []
        chats.value.push(newChat)
    } else {
        chats.value.push(newChat)
    }
}
const chatListAfterDelete = async () => {
    chats.value = await loadChats()
}
const emit = defineEmits(['chatSelected'])

const setId = (id, chatname) => {
    emit('chatSelected', id, chatname)
}
</script>

<template>
    <div id="content">
        <div id="top">
            <div id="user-registered">
                <img src="../../assets/user.svg" alt="">
                <p>{{ props.username }}</p>
            </div>
            <button @click="logout" id="logout-btn"><img src="../../assets/logout.svg" alt=""></button>
        </div>
        <hr>
        <div id="chat-header">
            <h3>Chats</h3>
            <AddChatDialog :name="props.username" @chatsUpdated="refreshChatList" />
        </div>
        <div v-if="filteredChats" id="chat-list">
            <div v-for="(item, i) in filteredChats" :key="i">
                <div v-for="(contact, index) in item.participants" :key="index" id="chat-card"
                    @click="setId(item.chat.id, contact.name)">
                    <div>
                        <p>{{ contact.name }}</p>
                    </div>
                    <button @click.stop="async () => {
                        await deleteChat(item.chat.id),
                            chatListAfterDelete()
                    }" id="del-chat-btn">
                        <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none"
                            stroke="#e6e6e6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                            class="icon icon-tabler icons-tabler-outline icon-tabler-trash">
                            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                            <path d="M4 7l16 0" />
                            <path d="M10 11l0 6" />
                            <path d="M14 11l0 6" />
                            <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12" />
                            <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3" />
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
#content {
    padding-left: 16px;
    padding-right: 16px;
    width: 100%;
    max-width: 280px;
    height: 100vh;
    overflow-y: hidden;
    box-sizing: border-box;
    border-right: 6px solid #065464;
    background-color: #2e2e2e;
}

@media (max-width: 768px) {
    #content {
        width: 100%;
        max-width: 200px;
    }
}

#chat-header {
    display: flex;
    justify-content: space-around;
    align-items: center;
}

#user-registered {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 80px;

    img {
        margin-right: 4px;
    }

    p {
        font-size: 18px;
    }
}

#chat-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-top: 3px solid #3b3b3b;
    padding: 3px 6px 3px 6px;
}

#chat-card:hover {
    cursor: pointer;
    background-color: #3b3b3b;
}

#top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 8px;
}

#top button {
    cursor: pointer;
    height: 30px;
}

#del-chat-btn {
    cursor: pointer;
    border: none;
    background-color: transparent;
}

#del-chat-btn:hover .icon-tabler-trash {
    stroke: #d83820;
    background-color: rgba(105, 105, 105, 0.5);
    border-radius: 9px;
}

.icon-tabler-trash {
    padding: 8px;
}

#logout-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 8px;
    border-radius: 8px;
    border: none;
    background-color: transparent;
    cursor: pointer;
    transition: 0.2s;

    &:hover {
        background-color: #126f81;
    }
}

#chat-list {
    overflow-y: auto;
}
</style>