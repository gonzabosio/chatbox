<script setup>
import { ref, onMounted } from 'vue';
import { axiosInstance } from '../../axios-func/axiosInstance';
import { addChat, loadChats } from '../../axios-func/calls';

const props = defineProps({
    name: String
})

let newChat = ''
let chat = {}
const emit = defineEmits(['chatsUpdated'])

const sendNewChatList = (chat) => {
    console.log(chat)
    emit('chatsUpdated', chat)
}
</script>

<template>
    <v-dialog max-width="500">
        <template v-slot:activator="{ props: activatorProps }">
            <v-btn v-bind="activatorProps" text="Add Chat" variant="flat" id="open-dialog-btn"></v-btn>
        </template>

        <template v-slot:default="{ isActive }">
                <v-card title="New chat">
                    <v-text-field v-model="newChat" label="Username" required></v-text-field>
    
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn text="Cancel" @click="isActive.value = false, newChat = ''"></v-btn>
                        <v-btn class="btn-add" text="Add" @click="async () => {
                            isActive.value = false, 
                            chat = await addChat(newChat, props.name), 
                            sendNewChatList(chat),
                            newChat = ''
                        }"></v-btn>
                    </v-card-actions>
                </v-card>
        </template>
    </v-dialog>
</template>

<style scoped>
#open-dialog-btn {
    background-color: rgb(66, 90, 61);
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
</style>