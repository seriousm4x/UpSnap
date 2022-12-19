import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const BACKEND_IS_PROXIED = import.meta.env.VITE_BACKEND_IS_PROXIED;
const BACKEND_PORT = import.meta.env.VITE_BACKEND_PORT;

const status = writable('');
const message = writable('');

let socket: WebSocket;

function initSocket() {
	if (BACKEND_IS_PROXIED) {
		const socketUrl = new URL('ws/wol', window.location.href);
		socketUrl.protocol = socketUrl.protocol.replace('http', 'ws');
		socket = new WebSocket(socketUrl);
	} else {
		socket = new WebSocket(`ws://${location.hostname}:${BACKEND_PORT}/ws/wol`);
	}

	// Connection opened
	socket.addEventListener('open', function () {
		status.set('open');
	});

	// Connection closed
	socket.addEventListener('close', function () {
		status.set('close');
		setTimeout(function () {
			initSocket();
		}, 3000);
	});

	// Connection error
	socket.addEventListener('error', function () {
		status.set('error');
		socket.close();
	});

	// Listen for messages
	socket.addEventListener('message', function (event) {
		message.set(JSON.parse(event.data));
	});
}

if (browser) {
	initSocket();
}

const sendMessage = (message: string | object) => {
	if (socket.readyState <= 1) {
		socket.send(JSON.stringify(message));
	}
};

export default {
	subscribeMsg: message.subscribe,
	subscribeStatus: status.subscribe,
	sendMessage
};
