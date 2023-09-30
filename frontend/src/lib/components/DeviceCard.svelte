<script lang="ts">
	import { goto } from '$app/navigation';
	import { backendUrl, permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device } from '$lib/types/device';
	import {
		faBed,
		faCircleArrowDown,
		faCircleArrowUp,
		faLock,
		faPen,
		faRotateLeft
	} from '@fortawesome/free-solid-svg-icons';
	import { formatDistance, parseISO } from 'date-fns';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';
	import { scale } from 'svelte/transition';
	import DeviceCardNic from './DeviceCardNic.svelte';

	export let device: Device;

	let moreButtons = [
		{
			text: 'Edit',
			icon: faPen,
			onClick: () => goto(`/device/${device.id}`),
			requires: $pocketbase.authStore.isAdmin || $permission.update?.includes(device.id)
		},
		{
			text: 'Sleep',
			icon: faBed,
			onClick: () => sleep(),
			requires:
				($pocketbase.authStore.isAdmin || $permission.power?.includes(device.id)) &&
				device.status === 'online' &&
				device.sol_enabled
		},
		{
			text: 'Reboot',
			icon: faRotateLeft,
			onClick: () => reboot(),
			requires:
				($pocketbase.authStore.isAdmin || $permission.power?.includes(device.id)) &&
				device.status === 'online' &&
				device.shutdown_cmd !== ''
		}
	];

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
			{#if moreButtons.filter((btn) => btn.requires).length > 0}
				<div class="dropdown dropdown-top dropdown-end bg-base-300 ms-auto">
					<label tabindex="-1" class="btn btn-sm m-1" for="more-{device.id}">More</label>
					<ul
						id="more-{device.id}"
						tabindex="-1"
						class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-fit"
					>
						{#each moreButtons as btn}
							{#if btn.requires}
								<li>
									<button on:click={btn.onClick}>
										<Fa icon={btn.icon} />{btn.text}
									</button>
								</li>
							{/if}
						{/each}
					</ul>
				</div>
			{/if}
		</div>
	</div>
</div>
