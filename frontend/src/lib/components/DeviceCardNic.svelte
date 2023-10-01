<script lang="ts">
	import LL from '$lib/i18n/i18n-svelte';
	import { backendUrl, permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device } from '$lib/types/device';
	import { faCircle, faPowerOff } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	export let device: Device;

	let hoverText = '';
	let disabled = false;
	let timeout = 120;
	var interval: number;

	$: if (device.status === 'pending' && !interval) {
		countdown(Date.parse(device.updated));
	}
	$: minutes = Math.floor(timeout / 60);
	$: seconds = timeout % 60;
	$: if (device.status === 'pending' || device.status === '') {
		disabled = true;
		hoverText = $LL.device.card_nic_tooltip_pending();
	} else if (device.status === 'online') {
		if (device.shutdown_cmd === '') {
			disabled = true;
			hoverText = $LL.device.card_nic_tooltip_shutdown_no_cmd();
		} else if (!$pocketbase.authStore.isAdmin && !$permission.power?.includes(device.id)) {
			disabled = true;
			hoverText = $LL.device.card_nic_tooltip_shutdown_no_permission();
		} else {
			disabled = false;
			hoverText = $LL.device.card_nic_tooltip_shutdown();
		}
	} else if (device.status === 'offline') {
		if (!$pocketbase.authStore.isAdmin && !$permission.power?.includes(device.id)) {
			disabled = true;
			hoverText = $LL.device.card_nic_tooltip_power_no_permission();
		} else {
			disabled = false;
			hoverText = $LL.device.card_nic_tooltip_power();
		}
	}

	// TODO: change wake and shutdown to nic based routes, not device based
	function wake() {
		fetch(`${backendUrl}api/upsnap/wake/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then(async (data) => {
				const dev = (await data.json()) as Device;
				countdown(Date.parse(dev.updated));
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function shutdown() {
		fetch(`${backendUrl}api/upsnap/shutdown/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then(async (data) => {
				const dev = (await data.json()) as Device;
				countdown(Date.parse(dev.updated));
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function countdown(updated: number) {
		timeout = 120;
		const end = updated + 2 * 60 * 1000;
		interval = setInterval(() => {
			timeout = Math.round((end - Date.now()) / 1000);
			if (timeout <= 0 || device.status !== 'pending') {
				clearInterval(interval);
			}
		}, 1000);
	}

	function handleClick() {
		if (device.status === 'offline') {
			wake();
		} else if (device.status === 'online') {
			shutdown();
		}
	}
</script>

<li class="tooltip" class:disabled data-tip={hoverText}>
	<div
		class="flex items-start p-2 gap-4"
		on:click={disabled ? null : handleClick}
		on:keydown={disabled ? null : handleClick}
		role="none"
	>
		{#if device.status === 'offline'}
			<button class="btn btn-error flex-shrink"><Fa icon={faPowerOff} /></button>
		{:else if device.status === 'online'}
			<button
				class="btn btn-success flex-shrink"
				class:cursor-not-allowed={device.shutdown_cmd === ''}><Fa icon={faPowerOff} /></button
			>
		{:else if device.status === 'pending'}
			<button class="btn btn-warning flex-shrink">
				<span class="countdown font-mono">
					<span style="--value:{minutes};" />:
					<span style="--value:{seconds};" />
				</span>
			</button>
		{:else}
			<button class="btn btn-warning flex-shrink">
				<span class="loading loading-ring loading-sm" />
			</button>
		{/if}
		<div class="grow">
			<div class="text-lg font-bold leading-4">{device.ip}</div>
			<div>{device.mac}</div>
			<div class="flex flex-wrap gap-x-4 gap-y-0">
				{#if device?.expand?.ports}
					{#each device?.expand?.ports.sort((a, b) => a.number - b.number) as port}
						<span class="flex items-center gap-1 break-all">
							{#if port.status}
								<Fa icon={faCircle} class="text-success" />
								{port.name} ({port.number})
							{:else}
								<Fa icon={faCircle} class="text-error" />
								{port.name} ({port.number})
							{/if}
						</span>
					{/each}
				{/if}
			</div>
		</div>
	</div>
</li>

<style>
	:global(.menu li.disabled) {
		color: inherit;
	}
</style>
