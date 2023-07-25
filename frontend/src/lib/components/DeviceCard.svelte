<script lang="ts">
	import type { Device } from '$lib/types/device';
	import { formatDistanceToNow, parseISO } from 'date-fns';
	import DeviceCardNic from './DeviceCardNic.svelte';

	export let device: Device;
</script>

<div class="card bg-base-300 shadow-md rounded-3xl">
	<div class="card-body p-6">
		{#if device.link !== ''}
			<a href={device.link} target="_blank">
				<h1 class="card-title link">{device.name}</h1>
			</a>
		{:else}
			<h1 class="card-title">{device.name}</h1>
		{/if}
		<ul class="menu bg-base-200 rounded-box">
			<!-- TODO: change to nic array once backend supports it -->
			<DeviceCardNic {device} />
		</ul>
		<div class="card-actions">
			<span class="tooltip" data-tip="Last status change">
				{formatDistanceToNow(parseISO(device.updated), {
					includeSeconds: true,
					addSuffix: true
				})}
			</span>
			<a class="btn btn-sm btn-neutral ms-auto" href="/device/{device.id}">Edit</a>
		</div>
	</div>
</div>
