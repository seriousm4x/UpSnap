<script lang="ts">
	import { onMount } from 'svelte';
	import { pocketbase } from '$lib/stores/pocketbase';
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import Alert from '$lib/components/Alert.svelte';
	import Fa from 'svelte-fa';
	import { faPlus, faTriangleExclamation } from '@fortawesome/free-solid-svg-icons';
	import type { Device } from '$lib/types/device';

	let devices = [] as Device[];
	let loading = true;
	let err = false;

	function getAllDevices() {
		$pocketbase
			.collection('devices')
			.getFullList(1000, { sort: 'name', expand: 'ports,groups' })
			.then((resp) => {
				devices = resp as Device[];
			})
			.catch(() => {
				err = true;
			})
			.finally(() => {
				loading = false;
			});
	}

	onMount(() => {
		getAllDevices();

		$pocketbase.collection('devices').subscribe('*', () => {
			getAllDevices();
		});

		$pocketbase.collection('ports').subscribe('*', () => {
			getAllDevices();
		});
	});
</script>

{#if loading}
	<div class="container max-w-lg mx-auto text-center">
		<span class="loading loading-dots loading-lg" />
	</div>
{:else if err}
	<div class="container max-w-lg mx-auto">
		<Alert color="error" message="Failed to get devices from api." icon={faTriangleExclamation} />
	</div>
{:else if devices.length > 0}
	<div
		class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4"
	>
		{#each devices as device}
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
