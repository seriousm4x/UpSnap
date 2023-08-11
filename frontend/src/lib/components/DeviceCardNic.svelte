<script lang="ts">
	import { pocketbase, backendUrl, permission } from '$lib/stores/pocketbase';
	import Fa from 'svelte-fa';
	import { faCircle, faPowerOff } from '@fortawesome/free-solid-svg-icons';
	import type { Device } from '$lib/types/device';

	export let device: Device;

	let interval: number;
	let timeout = 120;
	$: minutes = Math.floor(timeout / 60);
	$: seconds = timeout % 60;
	$: if (device.status === 'pending') {
		countdown();
	}

	// TODO: change wake and shutdown to nic based routes, not device based
	function wake() {
		fetch(`${backendUrl}api/upsnap/wake/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		});
	}

	function shutdown() {
		fetch(`${backendUrl}api/upsnap/shutdown/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		});
	}

	function countdown() {
		clearInterval(interval);
		const end = Date.parse(device.updated) + 2 * 60 * 1000;
		interval = setInterval(() => {
			const left = Math.round((end - Date.now()) / 1000);
			if (left >= 0) {
				timeout = left;
			} else {
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

<li
	class:disabled={(device.status === 'online' && device.shutdown_cmd === '') ||
		!$permission.power?.includes(device.id)}
	class:tooltip={(device.status === 'online' && device.shutdown_cmd === '') ||
		!$permission.power?.includes(device.id)}
	data-tip={device.status === 'online' && device.shutdown_cmd === ''
		? 'No shutdown command set'
		: !$permission.power?.includes(device.id)
		? 'No permission to shut down this device'
		: null}
>
	<div
		class="flex items-start p-2 gap-4"
		on:click={(device.status === 'online' && device.shutdown_cmd === '') ||
		!$permission.power?.includes(device.id)
			? null
			: handleClick}
		on:keydown={(device.status === 'online' && device.shutdown_cmd === '') ||
		!$permission.power?.includes(device.id)
			? null
			: handleClick}
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
					{#each device?.expand?.ports as port}
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
