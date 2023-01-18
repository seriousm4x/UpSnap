<script>
	import { onMount } from 'svelte';
	import DeviceCard from '@components/DeviceCard.svelte';
	import { pocketbase } from '@stores/pocketbase';

	let devices = {};

	onMount(async () => {
		let pb;
		pocketbase.subscribe(conn => {
			pb = conn;
		});

		// get all devices in pocketbase
		const result = await pb.collection('devices').getFullList();
		result.forEach((device) => {
			devices[device.id] = device;
		});

		// subscribe to database events
		pb.collection('devices').subscribe('*', function (e) {
			devices[e.record.id] = e.record;
		});
	});

	// update device date
	let now = Date.now();
	let clear;
	$: {
		clearInterval(clear);
		clear = setInterval(() => {
			now = Date.now();
		}, 1000);
	}
</script>

<div class="container text-body-emphasis">
	<div class="row">
		{#each Object.entries(devices) as [_, device]}
			<DeviceCard {device} {now} />
		{/each}
	</div>
</div>
