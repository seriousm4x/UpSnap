<script lang="ts">
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
	import type { Device } from '$lib/types/device';
	import Fa from 'svelte-fa';
	import { faCircle, faPowerOff } from '@fortawesome/free-solid-svg-icons';

	export let device: Device;

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

	function handleClick() {
		if (device.status === 'offline') {
			wake();
		} else if (device.statis === 'online') {
			shutdown();
		}
	}
</script>

<li>
	<div
		class="flex items-start p-2 gap-4"
		on:click={handleClick}
		on:keydown={handleClick}
		role="none"
	>
		<div>
			{#if device.status === 'offline'}
				<button class="btn btn-error flex-shrink"><Fa icon={faPowerOff} /></button>
			{:else if device.status === 'online'}
				<button class="btn btn-success flex-shrink"><Fa icon={faPowerOff} /></button>
			{:else}
				<button class="btn btn-warning flex-shrink"
					><span class="loading loading-ring loading-md" /></button
				>
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
