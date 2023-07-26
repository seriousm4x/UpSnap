<script lang="ts">
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import { pocketbase, devices } from '$lib/stores/pocketbase';
	import { onMount } from 'svelte';
	import type { Device, Port } from '$lib/types/device';

	onMount(() => {
		$pocketbase.collection('devices').subscribe('*', (e) => {
			const index = $devices.findIndex((device) => device.id === e.record.id);
			if (e.action === 'create') {
				$devices = [...$devices, e.record as Device];
			} else if (e.action === 'update') {
				// get expanded values before updating the store
				$pocketbase
					.collection('devices')
					.getOne(e.record.id, {
						expand: 'ports,groups'
					})
					.then((data) => {
						$devices[index] = data as Device;
					});
			} else if (e.action === 'delete') {
				$devices = $devices.filter((d) => d !== $devices[index]);
			}
		});

		$pocketbase.collection('ports').subscribe('*', (e) => {
			if (e.action === 'update') {
				// replace port in device store by finding both indexes
				let portIndex = -1;
				const deviceIndex = $devices.findIndex((device) =>
					device.expand.ports.map((port) => {
						if (port.id === e.record.id) {
							portIndex = device.expand.ports.indexOf(port);
						}
					})
				);
				if (deviceIndex !== -1 && portIndex !== -1) {
					$devices[deviceIndex].expand.ports[portIndex] = e.record as Port;
				}
			}
		});
	});
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
	<div class="container mx-auto text-center">No devices</div>
{/if}
