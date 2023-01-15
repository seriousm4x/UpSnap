<script>
	import PocketBase from 'pocketbase';
	import { onMount } from 'svelte';

	let devices = [];

	onMount(async () => {
		const pb = new PocketBase('http://127.0.0.1:8090');
		const result = await pb.collection('devices').getFullList();
		devices = result;
	});
</script>

{#each devices as device}
	{device.name}
{/each}
