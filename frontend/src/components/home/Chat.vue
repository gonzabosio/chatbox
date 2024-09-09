<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { deleteMessage, loadMessages } from '../../axios-func/calls'
import EditMessageDialog from './EditMessageDialog.vue';

const props = defineProps({
    chatName: String,
    chatId: String,
    userId: String
})

const messages = ref([])
const message = ref('')
let msgSenderSocket
let msgEditorSocket
const wsUrl = import.meta.env.VITE_WS_URL
onMounted(async () => {
    messages.value = await loadMessages(props.chatId)
    // Message Sender WebSocket
    msgSenderSocket = new WebSocket(`${wsUrl}/send-msg?wsauth=${localStorage.getItem('access-token')}`)
    msgSenderSocket.onopen = () => {
        console.log('Websocket connection opened to send messages')
    }
    msgSenderSocket.onerror = (err) => {
        console.log('Websocket error: ' + err)
    }
    msgSenderSocket.onmessage = (event) => {
        message.value = ''
        if (messages.value === null) {
            messages.value = []
            messages.value.push(JSON.parse(event.data))
        } else {
            messages.value.push(JSON.parse(event.data))
        }
    }
    msgSenderSocket.onclose = () => {
        console.log('Websocket connection closed')
    }
    // Message Editor WebSocket
    msgEditorSocket = new WebSocket(`${wsUrl}/edit-msg?wsauth=${localStorage.getItem('access-token')}`)
    msgEditorSocket.onopen = () => {
        console.log('Websocket connection opened to edit messages')
    }
    msgEditorSocket.onerror = (err) => {
        console.log('Websocket error: ' + err)
    }
    msgEditorSocket.onmessage = (event) => {
        messages.value.splice(msgIndex, 1, JSON.parse(event.data))
    }
    msgEditorSocket.onclose = () => {
        console.log('Websocket connection closed')
    }
})

watch(() => props.chatId, async (newChatId) => {
    messages.value = await loadMessages(newChatId)
})

const sendMsg = () => {
    if (message.value == '') {
        console.log('Empty message')
    } else {
        console.log('Sending message...')
        const msg = {
            chat_id: props.chatId,
            sender_id: localStorage.getItem('user-id'),
            content: message.value
        }
        msgSenderSocket.send(JSON.stringify(msg))
    }
}

let msgIndex = ''
const putMessage = (msgId, newMsg) => {
    const msgUpdate = {
        message_id: msgId,
        new_message: newMsg
    }
    msgEditorSocket.send(JSON.stringify(msgUpdate))
}

const deleteMsg = async (msgId, index) => {
    await deleteMessage(msgId)
    messages.value.splice(index, 1)[0]
}

onBeforeUnmount(() => {
    msgSenderSocket.close()
    msgEditorSocket.close()
})

</script>

<template>
    <div id="chat-container">
        <div class="header">
            <h2>{{ props.chatName }}</h2>
            <hr>
        </div>
        <div class="messages">
            <div v-for="(item, index) in messages" :key="index">
                <div :class="[item.sender_id === props.userId ? 'sent' : 'received']">
                    <div class="sub-msg-container">
                        <div v-if="item.sender_id === props.userId" class="options">
                            <EditMessageDialog :messageToEdit="item.content" :messageId="item.id"
                                @updateMessage="putMessage" @click="msgIndex = index" />
                            <button @click="deleteMsg(item.id, index)" id="del-msg-btn"><img
                                    src="../../assets/trash.svg" alt="delete">
                            </button>
                        </div>
                        <p class="msg-content">{{ item.content }}</p>
                    </div>
                </div>
            </div>
        </div>
        <div id="message-field">
            <input type="text" v-model="message">
            <button @click="sendMsg" id="send"><img src="../../assets/send.svg" alt=""></button>
        </div>
    </div>
</template>

<style scoped>
#chat-container {
    display: flex;
    width: 60%;
    flex-direction: column;
    height: 100vh;
    border-right: 6px solid #065464;
}

.header {
    h2 {
        margin-left: 16px;
    }

    hr {
        width: 97%;
    }
}

.messages {
    flex-grow: 1;
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: #85c3cf transparent;
    padding-top: 16px;
    padding-left: 16px;
    padding-right: 16px;
}

.sent {
    display: flex;
    justify-content: flex-end;
    margin-left: 120px;
}

.received {
    display: flex;
    justify-content: flex-start;
    margin-right: 120px;
}

.sent .msg-content,
.received .msg-content {
    display: flex;
    padding: 8px;
    border-radius: 9px;
}

.sent .msg-content {
    background-color: #2fa3b3;
}

.received .msg-content {
    background-color: #7a7d84;
}

.sent .sub-msg-container {
    display: flex;
    flex-direction: column;
    position: relative;
}

.options {
    display: flex;
    justify-content: flex-end;
    visibility: hidden;
    opacity: 0;
    position: absolute;
    right: 0px;
    top: -12px;
}

.sub-msg-container:hover .options {
    visibility: visible;
    opacity: 1;
    transition-delay: 0.5s;
    display: flex;
    align-items: center;
}

#del-msg-btn {
    cursor: pointer;
    width: 60px;
    padding: 6px;
    background-color: rgba(46, 46, 46, 0.37);
    transition: 0.3s;
    border: none;
    border-radius: 6px;
}

#del-msg-btn:hover {
    background-color: rgb(196, 66, 66);
    transition: 0.3s;
}

#message-field {
    display: flex;
    width: 100%;
    justify-content: center;
    margin-bottom: 10px;
    border: none;
}

input {
    width: 90%;
    height: 34px;
    background-color: transparent;
    border-radius: 10px;
    border: 1px solid darkgray;
}

input:focus-visible {
    outline: 1px solid rgba(169, 169, 169, 0.514);
}

#send {
    cursor: pointer;
    width: 60px;
    height: 40px;
    background-color: transparent;
    border: none;
    border-radius: 10px;
    margin-left: 4px;

    &:hover {
        background-color: #62626383;
    }
}
</style>