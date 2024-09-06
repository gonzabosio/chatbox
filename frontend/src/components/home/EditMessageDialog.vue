<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
    messageToEdit: { default: '', type: String },
    messageId: { default: '', type: String }
})

const message = ref(props.messageToEdit)
watch(() => props.messageToEdit, (newMsg) => {
    message.value = newMsg
})
const emit = defineEmits(['updateMessage'])

const emitNewMessage = (newMsg) => {
    emit('updateMessage', props.messageId, newMsg)
}
const errorMsg = ref('')
</script>

<template>
    <v-dialog max-width="500">
        <template v-slot:activator="{ props: activatorProps }">
            <button v-bind="activatorProps" id="edit-msg-btn">
                <img src="../../assets/edit.svg" alt="edit-btn-icon">
            </button>
        </template>

        <template v-slot:default="{ isActive }">
            <v-card title="Edit message">
                <v-text-field v-model="message" label="New message" required></v-text-field>
                <p>{{ errorMsg }}</p>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn text="Cancel" @click="isActive.value = false"></v-btn>
                    <v-btn class="btn-add" text="Save" @click="
                        message !== '' && message !== props.messageToEdit
                            ? (isActive.value = false, emitNewMessage(message))
                            : errorMsg = 'Invalid edition: empty or same message'">
                    </v-btn>
                </v-card-actions>
            </v-card>
        </template>
    </v-dialog>
</template>

<style scoped>
#edit-msg-btn {
    cursor: pointer;
    width: 60px;
    padding: 6px;
    background-color: rgba(46, 46, 46, 0.37);
    transition: 0.3s;
    border: none;
    border-radius: 6px;

    &:hover {
        background-color: #065464;
    }
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

p {
    margin-left: 20px;
    color: rgb(195, 70, 70);
}
</style>