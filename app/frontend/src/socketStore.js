import {
    writable
} from 'svelte/store';

const status = writable('');
const message = writable('');
let socket;

function initSocket() {
    socket = new WebSocket(`ws://${location.hostname}:${BACKEND_PORT}/wol/`);

    // Connection opened
    socket.addEventListener('open', function () {
        status.set("open");
    });

    // Connection closed
    socket.addEventListener('close', function () {
        status.set("close");
        setTimeout(function () {
            initSocket();
        }, 3000);
    });

    // Connection error
    socket.addEventListener('error', function () {
        status.set("error");
        socket.close()
    });

    // Listen for messages
    socket.addEventListener('message', function (event) {
        message.set(JSON.parse(event.data));
    });
}

initSocket()

const sendMessage = (message) => {
    if (socket.readyState <= 1) {
        socket.send(JSON.stringify(message));
    }
}

export default {
    subscribeMsg: message.subscribe,
    subscribeStatus: status.subscribe,
    sendMessage
}
