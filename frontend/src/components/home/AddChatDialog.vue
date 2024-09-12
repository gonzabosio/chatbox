<script setup>
import { ref } from 'vue';
import { addChat } from '../../axios-func/calls';

const props = defineProps({
    name: String
})

let newChat = ref('')
let chat = {}
const emit = defineEmits(['chatsUpdated'])
const errorMsg = ref('')
const handleAddClick = async () => {
    if (newChat.value !== '') {
        try {
            chat = await addChat(newChat.value, props.name)
            sendNewChatList(chat)
            newChat = ''
            errorMsg.value = ''
            return true
        }
        catch (err) {
            if (err.response.data.message == 'mongo: no documents in result') {
                errorMsg.value = 'Non-existent user'
            } else if (err.response.data.message == 'contact is already in the chats') {
                errorMsg.value = 'Already have this contact in your chats'
            } else {
                errorMsg.value = 'Non-existent user'
            }
            return false
        }
    } else {
        errorMsg.value = 'Username field is empty'
    }
}
const sendNewChatList = (chat) => {
    console.log(chat)
    emit('chatsUpdated', chat)
}

</script>

<template>

    <v-dialog max-width="500">
        <template v-slot:activator="{ props: activatorProps }">
            <button v-bind="activatorProps" id="open-dialog-btn"><img src="../../assets/user-plus.svg"
                    alt="add-chat-icon"></button>
        </template>

        <template v-slot:default="{ isActive }">
            <v-card title="New chat">
                <v-text-field v-model="newChat" label="Username" required></v-text-field>
                <p>{{ errorMsg }}</p>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn text="Cancel" @click="isActive.value = false, newChat = '', errorMsg = ''"></v-btn>
                    <v-btn class="btn-add" text="Add" @click="async () => {
                        let closeDialog = await handleAddClick()
                        closeDialog ? isActive.value = false : null
                    }"></v-btn>
                </v-card-actions>
            </v-card>
        </template>
    </v-dialog>
</template>

<style scoped>
#open-dialog-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 8px;
    border-radius: 8px;
    border: none;
    background-color: #065464;
    cursor: pointer;
    transition: 0.2s;

    &:hover {
        background-color: #126f81;
    }
}

.v-btn:hover {
    background-color: #6363630e;
    color: #ebebeb71;
}

.v-btn {
    cursor: pointer;
    background-color: #2e2e2e69;
    color: #bdbdbd8a;
}

.btn-add {
    cursor: pointer;
    background-color: rgba(97, 97, 97, 0.747);
    color: rgba(228, 228, 228, 0.562);
}

.btn-add:hover {
    color: rgba(228, 228, 228, 0.562);
}

p {
    margin-left: 20px;
    color: rgb(195, 70, 70);
}
</style>