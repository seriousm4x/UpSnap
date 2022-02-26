import {
    writable
} from 'svelte/store';

const messageStore = writable('');
const socket = new WebSocket(`ws://${location.hostname}:${BACKEND_PORT}/wol/`);

// Connection opened
socket.addEventListener('open', function (event) {
    return
});

// Listen for messages
socket.addEventListener('message', function (event) {
    messageStore.set(JSON.parse(event.data));
});

const sendMessage = (message) => {
    if (socket.readyState <= 1) {
        socket.send(JSON.stringify(message));
    }
}

export default {
    subscribe: messageStore.subscribe,
    sendMessage
}
