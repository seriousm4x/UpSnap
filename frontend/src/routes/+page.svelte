<script lang="ts">
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import { pocketbase, devices } from '$lib/stores/pocketbase';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import { faPlus } from '@fortawesome/free-solid-svg-icons';
	import type { Device } from '$lib/types/device';

	onMount(async () => {
		await getAllDevices();

		$pocketbase.collection('devices').subscribe('*', async () => {
			await getAllDevices();
		});

		$pocketbase.collection('ports').subscribe('*', async () => {
			await getAllDevices();
		});
	});

	async function getAllDevices() {
		const resp = await $pocketbase
			.collection('devices')
			.getFullList(1000, { sort: 'name', expand: 'ports,groups' });
		devices.set(resp as Device[]);
	}
</script>

{#if $devices.length > 0}
	<div
		class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4"
	>
		{#each $devices as device}
			<DeviceCard {device} />
		{/each}
	</div>
{:else}
	<div class="container text-center">
		<p>You have not created any devices.</p>
		<p>
			<a href="/device/new" class="btn btn-ghost"
				><Fa icon={faPlus} class="ms-2" />Add your first device
			</a>
		</p>
	</div>
{/if}
