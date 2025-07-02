<script lang="ts">
	import { m } from '$lib/paraglide/messages';
	import { backendUrl, permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device } from '$lib/types/device';
	import { faPowerOff } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	export let device: Device;

	let hoverText = '';
	let disabled = false;
	let timeout = 120;
	let interval: number;
	let modalWake: HTMLDialogElement;
	let modalShutdown: HTMLDialogElement;

	$: if (device.status === 'pending' && !interval) {
		// eslint-disable-next-line svelte/infinite-reactive-loop
		countdown(Date.parse(device.updated), 'wake');
	}
	$: minutes = Math.floor(timeout / 60);
	$: seconds = timeout % 60;
	$: if (device.status === 'pending' || device.status === '') {
		disabled = true;
		hoverText = m.device_card_nic_tooltip_pending();
	} else if (device.status === 'online') {
		if (device.shutdown_cmd === '') {
			disabled = true;
			hoverText = m.device_card_nic_tooltip_shutdown_no_cmd();
		} else if (!$pocketbase.authStore.isSuperuser && !$permission.power?.includes(device.id)) {
			disabled = true;
			hoverText = m.device_card_nic_tooltip_shutdown_no_permission();
		} else {
			disabled = false;
			hoverText = m.device_card_nic_tooltip_shutdown();
		}
	} else if (device.status === 'offline') {
		if (!$pocketbase.authStore.isSuperuser && !$permission.power?.includes(device.id)) {
			disabled = true;
			hoverText = m.device_card_nic_tooltip_power_no_permission();
		} else {
			disabled = false;
			hoverText = m.device_card_nic_tooltip_power();
		}
	}

	function wake() {
		countdown(Date.now(), 'wake');
		device.status = 'pending';

		fetch(`${backendUrl}api/upsnap/wake/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then((resp) => resp.json())
			.then(async (data) => {
				if (data.status !== 200) {
					device.status = 'offline';
					return;
				}
				device = data as Device;
				if (device.status === 'online' && device.link && device.link_open !== '') {
					if (device.link_open === 'new_tab') {
						window.open(device.link, '_blank');
					} else {
						window.open(device.link, '_self');
					}
				}
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function shutdown() {
		countdown(Date.now(), 'shutdown');
		device.status = 'pending';

		fetch(`${backendUrl}api/upsnap/shutdown/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then((resp) => resp.json())
			.then(async (data) => {
				if (data.status !== 200) {
					device.status = 'online';
					return;
				}
				device = data as Device;
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function countdown(updated: number, action: 'wake' | 'shutdown') {
		timeout = action === 'wake' ? device.wake_timeout : device.shutdown_timeout;
		if (timeout <= 0) {
			timeout = 120;
		}

		const end = updated + timeout * 1000;

		if (interval) {
			clearInterval(interval);
			interval = 0;
		}

		interval = setInterval(() => {
			timeout = Math.round((end - Date.now()) / 1000);

			if (timeout <= 0 || device.status !== 'pending') {
				clearInterval(interval);
				// eslint-disable-next-line svelte/infinite-reactive-loop
				interval = 0;
			}
		}, 1000);
	}

	function handleClick() {
		if (device.status === 'offline') {
			if (device.wake_confirm) {
				askConfirmation('wake');
			} else {
				wake();
			}
		} else if (device.status === 'online') {
			if (device.shutdown_confirm) {
				askConfirmation('shutdown');
			} else {
				shutdown();
			}
		}
	}

	function askConfirmation(action: string) {
		if (action === 'wake') {
			modalWake.showModal();
		} else {
			modalShutdown.showModal();
		}
	}
</script>

<div
	class={`tooltip ${disabled ? 'cursor-not-allowed' : 'hover:bg-base-300 cursor-pointer'} bg-base-100 rounded-box flex items-start gap-4 p-2`}
	data-tip={hoverText}
	on:click={disabled ? null : handleClick}
	on:keydown={disabled ? null : handleClick}
	role="none"
>
	{#if device.status === 'offline'}
		<button class="btn btn-error btn-circle size-12"><Fa icon={faPowerOff} /></button>
	{:else if device.status === 'online'}
		<button
			class="btn btn-success btn-circle size-12"
			class:cursor-not-allowed={device.shutdown_cmd === ''}><Fa icon={faPowerOff} /></button
		>
	{:else if device.status === 'pending'}
		<button class="btn btn-warning">
			<span class="countdown font-mono">
				<span style="--value:{minutes};"></span>:
				<span style="--value:{seconds};"></span>
			</span>
		</button>
	{:else}
		<div class="btn btn-warning btn-circle size-12">
			<span class="loading loading-ring loading-sm"></span>
		</div>
	{/if}
	<div class="grow">
		<div class="text-lg leading-4 font-bold">{device.ip}</div>
		<div>{device.mac}</div>
		<div class="flex flex-wrap gap-x-4">
			{#if device?.expand?.ports}
				{#each device?.expand?.ports.sort((a, b) => a.number - b.number) as port (port.id)}
					<span class="flex items-center gap-1 break-all">
						{#if port.status}
							<div class="inline-grid *:[grid-area:1/1]">
								<div class="status status-success h-3 w-3 animate-ping"></div>
								<div class="status status-success h-3 w-3"></div>
							</div>
						{:else}
							<div class="inline-grid *:[grid-area:1/1]">
								<div class="status status-error h-3 w-3 animate-ping"></div>
								<div class="status status-error h-3 w-3"></div>
							</div>
						{/if}
						{#if port.link}
							<a
								href={port.link}
								target="_blank"
								class="underline"
								on:click={(e) => e.stopPropagation()}>{port.name} ({port.number})</a
							>
						{:else}
							{port.name} ({port.number})
						{/if}
					</span>
				{/each}
			{/if}
		</div>
	</div>
</div>

<dialog class="modal" bind:this={modalWake}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">
			{m.device_modal_confirm_wake_title({ device: device.name })}
		</h3>
		<p class="py-4">{m.device_modal_confirm_wake_desc({ device: device.name })}</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">{m.buttons_cancel()}</button>
				<button class="btn btn-success" on:click={wake}>{m.buttons_confirm()}</button>
			</form>
		</div>
	</div>
</dialog>

<dialog class="modal" bind:this={modalShutdown}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">
			{m.device_modal_confirm_shutdown_title({ device: device.name })}
		</h3>
		<p class="py-4">{m.device_modal_confirm_shutdown_desc({ device: device.name })}</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">{m.buttons_cancel()}</button>
				<button class="btn btn-success" on:click={shutdown}>{m.buttons_confirm()}</button>
			</form>
		</div>
	</div>
</dialog>
