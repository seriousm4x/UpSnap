<script lang="ts">
	import { page } from '$app/stores';
	import { devices } from '$lib/stores/pocketbase';
	import { faTriangleExclamation } from '@fortawesome/free-solid-svg-icons';
	import Alert from '$lib/components/Alert.svelte';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import type { Device } from '$lib/types/device';

	let waitingForDevices = true;
	let device: Device | undefined;
	let deviceClone: Device | undefined;

	devices.subscribe((data) => {
		if (waitingForDevices === false) return;
		if (data.length === 0) return;

		device = data.find((dev) => dev.id === $page.params.id);
		if (device !== undefined) {
			deviceClone = device;
			if (deviceClone.expand.ports === undefined) {
				deviceClone.expand.ports = [];
			}
		}
		waitingForDevices = false;
	});
</script>

{#if waitingForDevices}
	<div class="container max-w-lg mx-auto text-center">
		<span class="loading loading-dots loading-lg" />
	</div>
{:else if device === undefined || deviceClone === undefined}
	<div class="container max-w-lg mx-auto">
		<Alert
			color="error"
			message="No device with id: {$page.params.id}"
			icon={faTriangleExclamation}
		/>
	</div>
{:else}
	<h1 class="text-3xl font-bold mb-8">{device.name}</h1>
	<DeviceForm device={deviceClone} />
{/if}
