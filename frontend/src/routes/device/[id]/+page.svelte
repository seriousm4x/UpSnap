<script lang="ts">
	import { page } from '$app/stores';
	import { pocketbase } from '$lib/stores/pocketbase';
	import PageLoading from '$lib/components/PageLoading.svelte';
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
	<PageLoading />
{:then device}
	<h1 class="text-3xl font-bold mb-8">{device.name}</h1>
	<DeviceForm {device} />
{:catch err}
	<div class="container max-w-lg mx-auto">
		{err}
	</div>
{/await}
