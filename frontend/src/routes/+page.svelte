<script>
	import PocketBase from 'pocketbase';
	import { onMount } from 'svelte';
	import DeviceCard from '@components/DeviceCard.svelte';

	let devices = {};

	onMount(async () => {
		const pb = new PocketBase('http://127.0.0.1:8090');

		// get all devices in pocketbase
		const result = await pb.collection('devices').getFullList();
		result.forEach(device => {
			devices[device.id] = device
		});

		// subscribe to database events
		pb.collection('devices').subscribe('*', function (e) {
			devices[e.record.id] = e.record
		});
	});

	// update device date
	let now = Date.now();
    let clear;
	$: {
        clearInterval(clear)
        clear = setInterval(() => {
            now = Date.now()
        }, 1000);
    }

</script>

<div class="container">
	<div class="row">
		{#each Object.entries(devices) as [_, device]}
			<DeviceCard {device} {now} />
		{/each}
	</div>
</div>
