<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import { permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Port } from '$lib/types/device';
	import toast from 'svelte-french-toast';

	$: if (Object.hasOwn($permission, 'update')) {
		if (!$pocketbase.authStore.isAdmin && !$permission.update.includes($page.params.id)) {
			toast(`You don't have permission to visit ${$page.url.pathname}`, {
				icon: '⛔'
			});
			goto('/');
		}
	}

	async function getDevice(): Promise<Device> {
		const resp = await $pocketbase
			.collection('devices')
			.getOne($page.params.id, { expand: 'ports,groups' });

		let device = resp as Device;

		if (!device.expand) device.expand = {};

		if (!device.expand.ports) {
			device.expand.ports = [] as Port[];
		}

		return resp as Device;
	}
</script>

{#await getDevice()}
	<PageLoading />
{:then device}
	<h1 class="text-3xl font-bold mb-8 sm:break-all">{device.name}</h1>
	<DeviceForm {device} />
{:catch err}
	<div class="container max-w-lg mx-auto">
		{err}
	</div>
{/await}
