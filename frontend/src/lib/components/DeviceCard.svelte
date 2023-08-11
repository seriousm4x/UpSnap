<script lang="ts">
	import { formatDistance, parseISO } from 'date-fns';
	import DeviceCardNic from './DeviceCardNic.svelte';
	import { scale } from 'svelte/transition';
	import Fa from 'svelte-fa';
	import { faCircleArrowDown, faCircleArrowUp, faLock } from '@fortawesome/free-solid-svg-icons';
	import { isAdmin, permission } from '$lib/stores/pocketbase';
	import type { Device } from '$lib/types/device';

	export let device: Device;

	// update device status change
	let now = Date.now();
	let interval: number;
	$: {
		clearInterval(interval);
		interval = setInterval(() => {
			now = Date.now();
		}, 1000);
	}
</script>

<div class="card bg-base-300 shadow-md rounded-3xl" transition:scale={{ delay: 0, duration: 200 }}>
	<div class="card-body p-6">
		{#if device.link.toString() !== ''}
			<a href={device.link.toString()} target="_blank">
				<h1 class="card-title link">{device.name}</h1>
			</a>
		{:else}
			<h1 class="card-title">{device.name}</h1>
		{/if}
		<ul class="menu bg-base-200 rounded-box">
			<!-- TODO: change to nic array once backend supports it -->
			<DeviceCardNic {device} />
		</ul>
		<div class="flex flex-row flex-wrap gap-2">
			{#if device.wake_cron_enabled}
				<div class="tooltip" data-tip="Wake cron">
					<span class="badge badge-success gap-1 p-3"
						><Fa icon={faCircleArrowUp} />{device.wake_cron}</span
					>
				</div>
			{/if}
			{#if device.shutdown_cron_enabled}
				<div class="tooltip" data-tip="Shutdown cron">
					<span class="badge badge-error gap-1 p-3"
						><Fa icon={faCircleArrowDown} />{device.shutdown_cron}</span
					>
				</div>
			{/if}
			{#if device.password}
				<div class="tooltip" data-tip="Wake password">
					<span class="badge gap-1 p-3"><Fa icon={faLock} />Password</span>
				</div>
			{/if}
		</div>
		<div class="card-actions mt-auto">
			<span class="tooltip" data-tip="Last status change: {device.updated}">
				{formatDistance(parseISO(device.updated), now, {
					includeSeconds: true,
					addSuffix: true
				})}
			</span>
			{#if $isAdmin || $permission.update?.includes(device.id)}
				<a class="btn btn-sm btn-neutral ms-auto" href="/device/{device.id}">Edit</a>
			{/if}
		</div>
	</div>
</div>
