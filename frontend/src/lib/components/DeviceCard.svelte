<script lang="ts">
	import { goto } from '$app/navigation';
	import { nextCronDate } from '$lib/helpers/cron';
	import { m } from '$lib/paraglide/messages';
	import { dateFnsLocale } from '$lib/stores/locale';
	import { backendUrl, permission, pocketbase } from '$lib/stores/pocketbase';
	import { type Device } from '$lib/types/device';
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

	let modalReboot: HTMLDialogElement;

	$: moreButtons = [
		{
			text: m.device_card_btn_more_sleep(),
			icon: faBed,
			onClick: () => sleep(),
			requires:
				($pocketbase.authStore.isSuperuser || $permission.power?.includes(device.id)) &&
				device.status === 'online' &&
				device.sol_enabled
		},
		{
			text: m.device_card_btn_more_reboot(),
			icon: faRotateLeft,
			onClick: () => askRebootConfirmation(),
			requires:
				($pocketbase.authStore.isSuperuser || $permission.power?.includes(device.id)) &&
				device.status === 'online' &&
				device.shutdown_cmd !== ''
		},
		{
			text: m.device_card_btn_more_edit(),
			icon: faPen,
			onClick: () => goto(`/device/${device.id}`),
			requires: $pocketbase.authStore.isSuperuser || $permission.update?.includes(device.id)
		}
	];

	// update device status change
	let now = Date.now();
	let interval: number;
	$: {
		clearInterval(interval);
		interval = setInterval(() => {
			// eslint-disable-next-line svelte/infinite-reactive-loop
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

	function askRebootConfirmation() {
		if (device.shutdown_confirm) {
			modalReboot.showModal();
		} else {
			reboot();
		}
	}
</script>

<div class="card bg-base-200 shadow-sm" transition:scale={{ delay: 0, duration: 200 }}>
	<div class="card-body p-6">
		{#if device.link !== ''}
			<a href={device.link} target="_blank">
				<h1 class="link card-title">{device.name}</h1>
			</a>
		{:else}
			<h1 class="card-title">{device.name}</h1>
		{/if}
		{#if device.description}
			<p class="grow-0">{device.description}</p>
		{/if}
		<div class="card rounded-box w-full">
			<DeviceCardNic {device} />
		</div>
		{#if device.wake_cron_enabled || device.shutdown_cron_enabled || device.password}
			<div class="mt-1 flex flex-row flex-wrap gap-2">
				{#if device.wake_cron_enabled}
					<div class="tooltip" data-tip={m.device_card_tooltip_wake_cron()}>
						<span class="badge badge-success gap-1 p-3"
							><Fa icon={faCircleArrowUp} />
							{nextCronDate(device.wake_cron)}
						</span>
					</div>
				{/if}
				{#if device.shutdown_cron_enabled}
					<div class="tooltip" data-tip={m.device_card_tooltip_shutdown_cron()}>
						<span class="badge badge-error gap-1 p-3"
							><Fa icon={faCircleArrowDown} />
							{nextCronDate(device.shutdown_cron)}
						</span>
					</div>
				{/if}
				{#if device.password}
					<div class="tooltip" data-tip={m.device_card_tooltip_wake_password()}>
						<span class="badge gap-1 p-3"><Fa icon={faLock} />{m.device_card_password()}</span>
					</div>
				{/if}
			</div>
		{/if}
		<div class="card-actions mt-auto items-center">
			<span
				class="tooltip"
				data-tip="{m.device_card_tooltip_last_status_change()}: {device.updated}"
			>
				{formatDistance(parseISO(device.updated), now, {
					includeSeconds: true,
					addSuffix: true,
					locale: $dateFnsLocale
				})}
			</span>
			{#if moreButtons.filter((btn) => btn.requires).length > 0}
				<div class="ms-auto flex flex-row flex-wrap gap-1">
					{#each moreButtons as btn, i (i)}
						{#if btn.requires}
							<div class="tooltip" data-tip={btn.text}>
								<button class="btn btn-sm btn-circle" on:click={btn.onClick}>
									<Fa icon={btn.icon} />
								</button>
							</div>
						{/if}
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<dialog class="modal" bind:this={modalReboot}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">
			{m.device_modal_confirm_shutdown_title({ device: device.name })}
		</h3>
		<p class="py-4">{m.device_modal_confirm_shutdown_desc({ device: device.name })}</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">{m.buttons_cancel()}</button>
				<button class="btn btn-success" on:click={reboot}>{m.buttons_confirm()}</button>
			</form>
		</div>
	</div>
</dialog>
