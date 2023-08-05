<script lang="ts">
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
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

<li>
	<div
		class="flex items-start p-2 gap-4"
		class:hover:bg-inherit={device.status === 'online' && device.shutdown_cmd === ''}
		class:hover:cursor-default={device.status === 'online' && device.shutdown_cmd === ''}
		on:click={handleClick}
		on:keydown={handleClick}
		role="none"
	>
		<div>
			{#if device.status === 'offline'}
				<button class="btn btn-error flex-shrink"><Fa icon={faPowerOff} /></button>
			{:else if device.status === 'online'}
				<button class="btn btn-success flex-shrink"><Fa icon={faPowerOff} /></button>
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
		</div>
		<div class="grow">
			<div class="text-lg font-bold leading-4">{device.ip}</div>
			<div class="">{device.mac}</div>
			<div class="flex flex-wrap gap-4">
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
