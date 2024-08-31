<script setup>
import { ref, onMounted } from 'vue';
import { axiosInstance } from '../../axios-func/axiosInstance';
import { addChat, loadChats } from '../../axios-func/calls';

const props = defineProps({
    name: String
})

let newChat = ''
let chats = {}

const emit = defineEmits(['chatsUpdated'])

//send it when chats change
const sendNewChatList = () => {
    console.log(chats)
    emit('chatsUpdated', chats)
}
</script>

<template>
    <v-dialog max-width="500">
        <template v-slot:activator="{ props: activatorProps }">
            <v-btn v-bind="activatorProps" color="surface-variant" text="Add Chat" variant="flat"></v-btn>
        </template>

        <template v-slot:default="{ isActive }">
                <v-card title="New chat">
                    <v-text-field v-model="newChat" label="Username" required></v-text-field>
    
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn text="Cancel" @click="isActive.value = false, newChat = ''"></v-btn>
                        <v-btn class="btn-add" text="Add" @click="isActive.value = false, chats = addChat(newChat, props.name), newChat = ''"></v-btn>
                    </v-card-actions>
                </v-card>
        </template>
    </v-dialog>
</template>

<style scoped>
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