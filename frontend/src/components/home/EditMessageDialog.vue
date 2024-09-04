<script setup>
import { ref } from 'vue';

const props = defineProps({
    messageToEdit: {default: '', type: String},
    messageId: {default: '', type: String},
    // index: {required:true, type: Number}
})

const message = ref(props.messageToEdit)

const emit = defineEmits(['updateMessage'])

const emitNewMessage = (newMsg) => {
    emit('updateMessage', props.messageId, newMsg)
}
</script>

<template>
    <v-dialog max-width="500">
        <template v-slot:activator="{ props: activatorProps }">
            <v-btn v-bind="activatorProps" class="edit-msg-btn"><img src="../../assets/edit.svg" alt="edit"></v-btn>
        </template>

        <template v-slot:default="{ isActive }">
                <v-card title="Edit message">
                    <v-text-field v-model="message" label="New message" required></v-text-field>
    
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn text="Cancel" @click="isActive.value = false, message = ''"></v-btn>
                        <v-btn class="btn-add" text="Save" @click="isActive.value = false, console.log(message), emitNewMessage(message), message = ''"></v-btn>
                    </v-card-actions>
                </v-card>
        </template>
    </v-dialog>
</template>

<style scoped>
.edit-msg-btn {
    transition: 0.5s !important;
}
.edit-msg-btn:hover {
    background-color: rgb(84, 126, 36) !important;
}

.v-btn {
    cursor: pointer;
    background-color: #2e2e2e69;
    color: #bdbdbd8a; 
}
.v-btn:hover {
    background-color: #6363630e;
    color: #ebebeb71;
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