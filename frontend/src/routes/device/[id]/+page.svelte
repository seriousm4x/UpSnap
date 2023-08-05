<script lang="ts">
	import { page } from '$app/stores';
	import { pocketbase } from '$lib/stores/pocketbase';
	import { faTriangleExclamation } from '@fortawesome/free-solid-svg-icons';
	import Alert from '$lib/components/Alert.svelte';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import type { Device, Port } from '$lib/types/device';

	async function getDevice(): Promise<Device> {
		const resp = await $pocketbase
			.collection('devices')
			.getOne($page.params.id, { expand: 'ports,groups' });

		let device = resp as Device;
		if (!device.expand.ports) {
			device.expand.ports = [] as Port[];
		}

		return resp as Device;
	}
</script>

{#await getDevice()}
	<div class="container max-w-lg mx-auto text-center">
		<span class="loading loading-dots loading-lg" />
	</div>
{:then device}
	<h1 class="text-3xl font-bold mb-8">{device.name}</h1>
	<DeviceForm {device} />
{:catch}
	<div class="container max-w-lg mx-auto">
		<Alert
			color="error"
			message="No device with id: {$page.params.id}"
			icon={faTriangleExclamation}
		/>
	</div>
{/await}
