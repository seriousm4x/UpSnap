<script lang="ts">
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import { pocketbase, devices } from '$lib/stores/pocketbase';
	import type { Device } from '$lib/types/device';
	import { onMount } from 'svelte';

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
				delete $devices[index];
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
					$devices[deviceIndex].expand.ports[portIndex] = e.record;
				}
			}
		});
	});
</script>

{#each $devices as device}
	<div
		class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4"
	>
		<DeviceCard {device} />
	</div>
{:else}
	<div class="container mx-auto text-center">No devices</div>
{/each}
