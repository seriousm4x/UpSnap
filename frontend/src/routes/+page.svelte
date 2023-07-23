<script lang="ts">
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import { devices } from '$lib/stores/pocketbase';
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
