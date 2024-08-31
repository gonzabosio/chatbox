<script setup>
import { onMounted, ref, watch, watchEffect } from 'vue';
import { loadMessages, sendMessage } from '../../axios-func/calls'

const props = defineProps({
    chatName: String,
    chatId: String,
    userId: String
})

const messages = ref([])
onMounted(async () => {
    messages.value = await loadMessages(props.chatId)
})

watch(() => props.chatId, async (newChatId) => {
    messages.value = await loadMessages(newChatId)
})

const message = ref('')
const sendMsg = () => {
    sendMessage(props.chatId, localStorage.getItem('user-id'), message.value)
    console.log(message.value)
}
</script>

<template>
    <div id="chat-container">
        <h2>{{ props.chatName }}</h2>
        <hr>
        <div class="messages">
            <div v-for="(item, index) in messages" :key="index">
                <div :class="{
                    'sent': item.sender_id === props.userId,
                    'received': item.sender_id !== props.userId
                }">
                    <p class="msg-content">{{ item.content }}</p>
                </div>
            </div>
        </div>
        <div id="message-field">
            <input type="text" v-model="message">
            <button @click="sendMsg" id="send">SEND</button>
        </div>
    </div>
</template>

<style scoped>
#chat-container {
    display: flex;
    width: 60%;
    flex-direction: column;
    height: 100vh;
    border-right: 6px solid rgba(70, 160, 70, 0.493);
}

h2 {
    margin-left: 16px;
    margin-right: 16px;
}

hr {
    width: 99%;
}

.messages {
    flex-grow: 1;
    padding-left: 16px;
    padding-right: 16px;
}

.sent {
    display: flex;
    justify-content: flex-end;
}

.received {
    display: flex;
    justify-content: flex-start;
}

.sent .msg-content,
.received .msg-content {
    padding: 8px;
    border-radius: 9px;
}

.sent .msg-content {
    background-color: #4a7555;
}

.received .msg-content {
    background-color: rgb(122, 122, 122);
}

#message-field {
    display: flex;
    width: 100%;
    justify-content: center;
    margin-bottom: 10px;
}

input {
    width: 90%;
    height: 34px;
    background-color: transparent;
    border-radius: 10px;
}

input:focus-visible {
    outline: 1px solid rgba(169, 169, 169, 0.514);
}

#send {
    cursor: pointer;
    width: 60px;
    height: 40px;
    background-color: rgba(70, 160, 70, 0.493);
    border: none;
    border-radius: 10px;
    margin-left: 4px;
}
</style>