<script>
	import { onMount } from 'svelte';
	import store from './store.js';
	import Navbar from "./components/Navbar.svelte";
	import DeviceCard from "./components/DeviceCard.svelte"
	import DummyDeviceCard from "./components/DummyDeviceCard.svelte"

	let visitors = 0;
	let devices = {};
	let dummyDevices = [];
	let ports = [];

	onMount(() => {
		store.subscribe(currentMessage => {
			if (currentMessage.type == "init") {
				// create placeholder devices
				dummyDevices = currentMessage.message;
			} else if (currentMessage.type == "status") {
				// set device status and pulse animation
				dummyDevices -= 1;
				devices[currentMessage.message.device.id] = currentMessage.message.device;
				if (currentMessage.message.device.up == true) {
					setUp(currentMessage.message.device.id);
				} else {
					setDown(currentMessage.message.device.id);
				}
			} else if (currentMessage.type == "wake") {
				// set device waking
				console.log("WAKE:", currentMessage.message);
				setWake(currentMessage.message)
			} else if (currentMessage.type == "visitor") {
				// update visitor count
				visitors = currentMessage.message;
			} else if (currentMessage.type == "get_ports") {
				ports = [...currentMessage.message];
			} else if (currentMessage.type == "delete") {
				const devCol = document.querySelector(`#device-col-${currentMessage.message}`);
				devCol.remove();
			}
		})
	})

	function onSendMessage() {
		 if (message.length > 0) {
			 store.sendMessage(message);
			 message = "";
		 }
	}

	function setUp(id) {
		const dot = document.querySelector(`#dot-${id}`)
		const spinner = document.querySelector(`#spinner-${id}`);
		if (dot) {
			dot.style.animation = "none";
			dot.offsetWidth;
			if (!spinner.classList.contains("d-none")) {
				if (dot.classList.contains("danger")) {
					spinner.classList.add("d-none");
					dot.classList.remove("d-none", "danger");
					dot.classList.add("success");
				}
			} else {
				dot.style.animation = "on-pulse 1s normal";
			}
		}
	}

	function setDown(id) {
		const dot = document.querySelector(`#dot-${id}`)
		const spinner = document.querySelector(`#spinner-${id}`);
		if (dot) {
			dot.style.animation = "none";
			dot.offsetWidth;
			if (!spinner.classList.contains("d-none")) {
				if (dot.classList.contains("success")) {
					spinner.classList.add("d-none");
					dot.classList.remove("d-none", "success");
					dot.classList.add("danger");
				}
			} else {
				dot.style.animation = "off-pulse 1s normal";
			}
		}
	}

	function setWake(id) {
		const dot = document.querySelector(`#dot-${id}`);
		const spinner = document.querySelector(`#spinner-${id}`);
		dot.classList.add("d-none");
		spinner.classList.remove("d-none");
	}
</script>

<main>
	<Navbar visitors={visitors} />
	<div class="container">
		<div class="row">
			{#if dummyDevices >= 1}
				{#each Array(dummyDevices) as _, i}
					<DummyDeviceCard/>
				{/each}
			{/if}
			{#each Object.entries(devices) as [id, device]}
				<DeviceCard device={device} ports={ports}/>
			{/each}
		</div>
	</div>
</main>

<style lang="scss" global>
	@import "variables.scss";
	@import "../node_modules/bootstrap/scss/bootstrap";
	@import "../node_modules/@fortawesome/fontawesome-free/scss/fontawesome.scss";
	@import "../node_modules/@fortawesome/fontawesome-free/scss/regular.scss";
	@import "../node_modules/@fortawesome/fontawesome-free/scss/solid.scss";

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
