<script setup>
import { ref, onMounted } from 'vue';
import { axiosInstance } from '../../axios-func/axiosInstance';
import getChats from '../../axios-func/calls';

const props = defineProps({
    name: String
})

let newChat = ''
let chats = {}
function addChat() {
    axiosInstance(({
        method: 'post',
        url: '/chat/add',
        headers: {
            Authorization: 'Bearer '+localStorage.getItem('access-token')
        },
        data: {
            username: newChat,
            petitioner_id: localStorage.getItem('user-id'),
            petitioner: props.name
        }
    })).then(async res => {
        console.log(res.data.message)
        chats = await getChats()
        sendNewChatList()
    }).catch(err => {
        console.log(err.response.data.message)
    })
}

const emit = defineEmits(['chatsUpdated'])

const sendNewChatList = () => {
    console.log(chats)
    emit('chatsUpdated', chats)
}
</script>

<template>
    <v-dialog max-width="500">
        <template v-slot:activator="{ props: activatorProps }">
            <v-btn
            v-bind="activatorProps"
            color="surface-variant"
            text="Add Chat"
            variant="flat"
            ></v-btn>
        </template>

        <template v-slot:default="{ isActive }">
            <v-card title="New chat">
                <v-text-field
                v-model="newChat"
                label="Username"
                required
              ></v-text-field>

                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn text="Cancel" @click="isActive.value = false, newChat=''"></v-btn>
                    <v-btn text="Add" @click="isActive.value = false, addChat(), newChat=''"></v-btn>
                </v-card-actions>
            </v-card>
        </template>
    </v-dialog>
</template>

<style scoped></style>