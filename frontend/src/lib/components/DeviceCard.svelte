<script lang="ts">
	import { pocketbase, isAdmin, backendUrl, permission } from '$lib/stores/pocketbase';
	import DeviceCardNic from './DeviceCardNic.svelte';
	import Fa from 'svelte-fa';
	import {
		faBed,
		faCircleArrowDown,
		faCircleArrowUp,
		faLock,
		faPen,
		faRotateLeft
	} from '@fortawesome/free-solid-svg-icons';
	import { scale } from 'svelte/transition';
	import { formatDistance, parseISO } from 'date-fns';
	import toast from 'svelte-french-toast';
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

	function sleep() {
		fetch(`${backendUrl}api/upsnap/sleep/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		}).catch((err) => {
			toast.error(err.message);
		});
	}

	function reboot() {
		fetch(`${backendUrl}api/upsnap/reboot/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		}).catch((err) => {
			toast.error(err.message);
		});
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
		{#if device.wake_cron_enabled || device.shutdown_cron_enabled || device.password}
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
		{/if}
		<div class="card-actions mt-auto items-center">
			<span class="tooltip" data-tip="Last status change: {device.updated}">
				{formatDistance(parseISO(device.updated), now, {
					includeSeconds: true,
					addSuffix: true
				})}
			</span>
			{#if $isAdmin || $permission.update?.includes(device.id) || $permission.power?.includes(device.id)}
				<div class="dropdown dropdown-top dropdown-end bg-base-300 ms-auto">
					<label tabindex="-1" class="btn btn-sm m-1" for="more-{device.id}">More</label>
					<ul
						id="more-{device.id}"
						tabindex="-1"
						class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-fit"
					>
						{#if ($isAdmin || $permission.power?.includes(device.id)) && device.status === 'online'}
							{#if device.shutdown_cmd !== ''}
								<li><button on:click={() => reboot()}><Fa icon={faRotateLeft} />Reboot</button></li>
							{/if}
							{#if device.sol_enabled}
								<li><button on:click={() => sleep()}><Fa icon={faBed} />Sleep</button></li>
							{/if}
						{/if}
						{#if $isAdmin || $permission.update?.includes(device.id)}
							<li>
								<a href="/device/{device.id}"><Fa icon={faPen} />Edit</a>
							</li>
						{/if}
					</ul>
				</div>
			{/if}
		</div>
	</div>
</div>
