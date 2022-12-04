<script>
	import { onMount } from "svelte";
	import socketStore from "./socketStore.js";
	import Navbar from "./components/Navbar.svelte";
	import DeviceCard from "./components/DeviceCard.svelte";
	import Toast from "./components/Toast.svelte";

	let visitors = 0;
	let devices = [];
	let settings = {};
	let toast = {
		title: "",
		message: "",
		color: "",
		show: false,
	};

	onMount(() => {
		socketStore.subscribeStatus((status) => {
			if (status == "open") {
				showToast("Websocket", "Connected", "success");
			} else if (status == "close") {
				showToast(
					"Websocket",
					"Connection closed. Trying to reconnect ...",
					"danger"
				);
			} else if (status == "error") {
				showToast(
					"Websocket",
					"Error when connecting to websocket",
					"danger"
				);
			}
		});

		socketStore.subscribeMsg((currentMessage) => {
			if (currentMessage.type == "init") {
				// create devices
				devices = [...currentMessage.message.devices];
				devices = devices;
				devices.sort(compare);
				settings = currentMessage.message.settings;
			} else if (currentMessage.type == "status") {
				// set device array and sort
				const index = devices.findIndex(
					(x) => x.id == currentMessage.message.id
				);
				if (devices.length === 0 || index === -1) {
					devices.push(currentMessage.message);
					devices = devices;
				} else {
					devices[index] = currentMessage.message;
				}
				devices.sort(compare);
				// set device status
				if (currentMessage.message.up == true) {
					setUp(currentMessage.message);
				} else {
					setDown(currentMessage.message);
				}
			} else if (currentMessage.type == "pending") {
				// set device pending
				setPending(currentMessage.message);
			} else if (currentMessage.type == "visitor") {
				// update visitor count
				visitors = currentMessage.message;
			} else if (currentMessage.type == "delete") {
				// delete device
				const devCol = document.querySelector(
					`#device-col-${currentMessage.message}`
				);
				devCol.remove();
			} else if (currentMessage.type == "scan_network") {
				// set scanned network devices
				if (!currentMessage.message) {
					return;
				}
				settings["scan_network"] = currentMessage.message;
				const btnScan = document.querySelector("#btnScan");
				const btnScanSpinner =
					document.querySelector("#btnScanSpinner");
				const btnScanText = document.querySelector("#btnScanText");
				btnScan.disabled = false;
				btnScanSpinner.classList.add("d-none");
				btnScanText.innerText = "Scan";
			} else if (currentMessage.type == "backup") {
				// download backup file
				const now = new Date();
				const fileName = `upsnap_backup_${now.toISOString()}.json`;
				const a = document.createElement("a");
				const file = new Blob(
					[JSON.stringify(currentMessage.message)],
					{ type: "text/plain" }
				);
				a.href = URL.createObjectURL(file);
				a.download = fileName;
				a.click();
			} else if (currentMessage.type == "operationStatus") {
				if (currentMessage.message == "Success") {
					showToast(
						currentMessage.message,
						"Changes were saved",
						"success"
					);
				} else if (currentMessage.message == "Error") {
					showToast(
						currentMessage.message,
						"Error while saving the device. Please check the logs.",
						"danger"
					);
				}
			}
		});
	});

	function setUp(device) {
		const dot = document.querySelector(`#dot-${device.id}`);
		const spinner = document.querySelector(`#spinner-${device.id}`);
		if (dot) {
			if (dot.classList.contains("danger")) {
				showToast(device.name, "Device is up!", "success");
			}
			dot.style.animation = "none";
			dot.offsetWidth;
			if (!spinner.classList.contains("d-none")) {
				spinner.classList.add("d-none");
				dot.classList.remove("d-none", "danger");
				dot.classList.add("success");
			} else {
				dot.style.animation = "on-pulse 1s normal";
			}
		}
	}

	function setDown(device) {
		const dot = document.querySelector(`#dot-${device.id}`);
		const spinner = document.querySelector(`#spinner-${device.id}`);
		if (dot) {
			if (dot.classList.contains("success")) {
				showToast(device.name, "Device is down!", "danger");
			}
			dot.style.animation = "none";
			dot.offsetWidth;
			if (!spinner.classList.contains("d-none")) {
				spinner.classList.add("d-none");
				dot.classList.remove("d-none", "success");
				dot.classList.add("danger");
			} else {
				dot.style.animation = "off-pulse 1s normal";
			}
		}
	}

	function setPending(id) {
		const dot = document.querySelector(`#dot-${id}`);
		const spinner = document.querySelector(`#spinner-${id}`);
		dot.classList.add("d-none");
		spinner.classList.remove("d-none");
	}

	function showToast(title, message, color) {
		if (settings.notifications === false) {
			return;
		}
		toast.title = title;
		toast.message = message;
		toast.color = color;
		toast.show = true;
		setTimeout(() => {
			toast.show = false;
		}, 4000);
	}

	function compare(a, b) {
		if (a.name < b.name) {
			return -1;
		}
		if (a.name > b.name) {
			return 1;
		}
		return 0;
	}
</script>

<main>
	<Navbar {settings} {visitors} />
	<div class="container mb-3">
		<div class="row">
			{#each devices as device}
				<DeviceCard {device} />
			{/each}
		</div>
	</div>
	<Toast {toast} />
</main>

<style lang="scss" global>
	@import "./variables.scss";
	@import "../node_modules/bootstrap/scss/bootstrap";
	@import "../node_modules/@fortawesome/fontawesome-free/scss/fontawesome.scss";
	@import "../node_modules/@fortawesome/fontawesome-free/scss/regular.scss";
	@import "../node_modules/@fortawesome/fontawesome-free/scss/solid.scss";
	@import "theme.scss";

	body {
		color: var(--color-text);
		background-color: var(--color-bg);
	}

	.modal-content {
		background-color: var(--bg-modal);
	}

	.modal-header {
		border-bottom: 1px solid var(--color-text);

		.btn-close {
			background: var(--svg-close);
		}
	}

	.modal-footer {
		border-top: 1px solid var(--color-text);
	}

	.btn,
	button {
		&.btn-light {
			color: var(--color-text);
			background-color: var(--bg-lighter);
			border-width: 0px;
		}

		&.btn-outline-success {
			border-color: $success;
			&:hover {
				background-color: $success;
			}
		}

		&.btn-outline-danger {
			border-color: $danger;
			&:hover {
				background-color: $danger;
			}
		}
	}

	.form-control {
		&:focus {
			border-color: inherit;
			box-shadow: none;
		}
	}

	.callout {
		padding: 1rem;
		border-left-width: 0.25rem;
		border-radius: 0.25rem;

		&.callout-info {
			background-color: var(--info-dark-transparent);
			border: 1px solid var(--info-dark-transparent);
			border-left: 5px solid var(--info);
		}

		&.callout-danger {
			background-color: var(--danger-dark-transparent);
			border: 1px solid var(--danger-dark-transparent);
			border-left: 5px solid var(--danger);
		}
	}

	@keyframes on-pulse {
		0% {
			box-shadow: 0 0 0 0 $success;
		}

		70% {
			box-shadow: 0 0 0 15px rgba(0, 0, 0, 0);
		}

		100% {
			box-shadow: 0 0 0 0 rgba(0, 0, 0, 0);
		}
	}

	@keyframes off-pulse {
		0% {
			box-shadow: 0 0 0 0 $danger;
		}

		70% {
			box-shadow: 0 0 0 15px rgba(0, 0, 0, 0);
		}

		100% {
			box-shadow: 0 0 0 0 rgba(0, 0, 0, 0);
		}
	}
</style>
