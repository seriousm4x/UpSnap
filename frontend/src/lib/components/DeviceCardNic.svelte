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
	let modalWake: HTMLDialogElement;
	let modalShutdown: HTMLDialogElement;

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
		} else if (!$pocketbase.authStore.isSuperuser && !$permission.power?.includes(device.id)) {
			disabled = true;
			hoverText = $LL.device.card_nic_tooltip_shutdown_no_permission();
		} else {
			disabled = false;
			hoverText = $LL.device.card_nic_tooltip_shutdown();
		}
	} else if (device.status === 'offline') {
		if (!$pocketbase.authStore.isSuperuser && !$permission.power?.includes(device.id)) {
			disabled = true;
			hoverText = $LL.device.card_nic_tooltip_power_no_permission();
		} else {
			disabled = false;
			hoverText = $LL.device.card_nic_tooltip_power();
		}
	}

	function wake() {
		fetch(`${backendUrl}api/upsnap/wake/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then((resp) => resp.json())
			.then(async (data) => {
				device = data as Device;
				await countdown(Date.parse(device.updated));
				if (device.status === 'online' && device.link && device.link_open !== '') {
					console.log('here');
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
		fetch(`${backendUrl}api/upsnap/shutdown/${device.id}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then((resp) => resp.json())
			.then((data) => {
				device = data as Device;
				countdown(Date.parse(device.updated));
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function countdown(updated: number) {
		return new Promise((resolve, reject) => {
			try {
				timeout = 120;
				const end = updated + 2 * 60 * 1000;

				if (interval) {
					clearInterval(interval);
					interval = 0;
				}

				interval = setInterval(() => {
					timeout = Math.round((end - Date.now()) / 1000);

					if (timeout <= 0 || device.status !== 'pending') {
						clearInterval(interval);
						interval = 0;
						resolve(interval);
					}
				}, 1000);
			} catch (error) {
				reject(error);
			}
		});
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

<li class="tooltip" class:disabled data-tip={hoverText}>
	<div
		class="flex items-start gap-4 p-2"
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
			<div class="flex flex-wrap gap-x-4">
				{#if device?.expand?.ports}
					{#each device?.expand?.ports.sort((a, b) => a.number - b.number) as port}
						<span class="flex items-center gap-1 break-all">
							{#if port.status}
								<Fa icon={faCircle} class="text-success" />
							{:else}
								<Fa icon={faCircle} class="text-error" />
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
</li>

<dialog class="modal" bind:this={modalWake}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">
			{$LL.device.modal_confirm_wake_title({ device: device.name })}
		</h3>
		<p class="py-4">{$LL.device.modal_confirm_wake_desc({ device: device.name })}</p>
		<div class="modal-action">
			<form method="dialog" class="flex flex-row flex-wrap gap-2">
				<button class="btn btn-neutral">{$LL.buttons.cancel()}</button>
				<button class="btn btn-success" on:click={wake}>{$LL.buttons.confirm()}</button>
			</form>
		</div>
	</div>
</dialog>

<dialog class="modal" bind:this={modalShutdown}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">
			{$LL.device.modal_confirm_shutdown_title({ device: device.name })}
		</h3>
		<p class="py-4">{$LL.device.modal_confirm_shutdown_desc({ device: device.name })}</p>
		<div class="modal-action">
			<form method="dialog" class="flex flex-row flex-wrap gap-2">
				<button class="btn btn-neutral">{$LL.buttons.cancel()}</button>
				<button class="btn btn-success" on:click={shutdown}>{$LL.buttons.confirm()}</button>
			</form>
		</div>
	</div>
</dialog>

<style>
	:global(.menu li.disabled) {
		color: inherit;
	}
</style>
