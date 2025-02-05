<script lang="ts">
	import { goto } from '$app/navigation';
	import LL, { locale } from '$lib/i18n/i18n-svelte';
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
	import cronParser from 'cron-parser';
	import type { Locale } from 'date-fns';
	import { formatDistance, formatRelative, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale/de';
	import { enUS } from 'date-fns/locale/en-US';
	import { es } from 'date-fns/locale/es';
	import { fr } from 'date-fns/locale/fr';
	import { it } from 'date-fns/locale/it';
	import { ja } from 'date-fns/locale/ja';
	import { nl } from 'date-fns/locale/nl';
	import { pl } from 'date-fns/locale/pl';
	import { pt } from 'date-fns/locale/pt';
	import { zhCN } from 'date-fns/locale/zh-CN';
	import { zhTW } from 'date-fns/locale/zh-TW';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';
	import { scale } from 'svelte/transition';
	import DeviceCardNic from './DeviceCardNic.svelte';

	export let device: Device;

	let modalReboot: HTMLDialogElement;

	$: moreButtons = [
		{
			text: $LL.device.card_btn_more_sleep(),
			icon: faBed,
			onClick: () => sleep(),
			requires:
				($pocketbase.authStore.isSuperuser || $permission.power?.includes(device.id)) &&
				device.status === 'online' &&
				device.sol_enabled
		},
		{
			text: $LL.device.card_btn_more_reboot(),
			icon: faRotateLeft,
			onClick: () => askRebootConfirmation(),
			requires:
				($pocketbase.authStore.isSuperuser || $permission.power?.includes(device.id)) &&
				device.status === 'online' &&
				device.shutdown_cmd !== ''
		},
		{
			text: $LL.device.card_btn_more_edit(),
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
			now = Date.now();
		}, 1000);
	}

	// dynamically import locales
	// it's ugly but the only working solution i found in ~ 3 hours
	var dateFnsLocale: Locale;
	$: if ($locale !== undefined) {
		(async () => {
			switch ($locale) {
				case 'de':
				case 'de-DE':
					dateFnsLocale = de;
					break;
				case 'en':
				case 'en-US':
					dateFnsLocale = enUS;
					break;
				case 'fr':
				case 'fr-FR':
					dateFnsLocale = fr;
					break;
				case 'it':
				case 'it-IT':
					dateFnsLocale = it;
					break;
				case 'ja':
				case 'ja-JP':
					dateFnsLocale = ja;
					break;
				case 'es':
				case 'es-ES':
					dateFnsLocale = es;
					break;
				case 'nl':
				case 'nl-NL':
					dateFnsLocale = nl;
					break;
				case 'pl':
				case 'pl-PL':
					dateFnsLocale = pl;
					break;
				case 'pt':
				case 'pt-PT':
					dateFnsLocale = pt;
					break;
				case 'zh':
				case 'zh-CN':
					dateFnsLocale = zhCN;
					break;
				case 'zh-TW':
					dateFnsLocale = zhTW;
					break;
				default:
					dateFnsLocale = enUS;
					break;
			}
		})();
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

	function getNextCronRelativeTime(expression: string, now: number) {
		const cron = cronParser.parseExpression(expression);
		return formatRelative(cron.next().toISOString(), now, {
			locale: dateFnsLocale
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
		{#if device.link.toString() !== ''}
			<a href={device.link.toString()} target="_blank">
				<h1 class="link card-title">{device.name}</h1>
			</a>
		{:else}
			<h1 class="card-title">{device.name}</h1>
		{/if}
		{#if device.description}
			<p>{device.description}</p>
		{/if}
		<div class="card rounded-box w-full">
			<DeviceCardNic {device} />
		</div>
		{#if device.wake_cron_enabled || device.shutdown_cron_enabled || device.password}
			<div class="mt-1 flex flex-row flex-wrap gap-2">
				{#if device.wake_cron_enabled}
					<div class="tooltip" data-tip={$LL.device.card_tooltip_wake_cron()}>
						<span class="badge badge-success gap-1 p-3"
							><Fa icon={faCircleArrowUp} />{getNextCronRelativeTime(device.wake_cron, now)}</span
						>
					</div>
				{/if}
				{#if device.shutdown_cron_enabled}
					<div class="tooltip" data-tip={$LL.device.card_tooltip_shutdown_cron()}>
						<span class="badge badge-error gap-1 p-3"
							><Fa icon={faCircleArrowDown} />{getNextCronRelativeTime(
								device.shutdown_cron,
								now
							)}</span
						>
					</div>
				{/if}
				{#if device.password}
					<div class="tooltip" data-tip={$LL.device.card_tooltip_wake_password()}>
						<span class="badge gap-1 p-3"><Fa icon={faLock} />{$LL.device.card_password()}</span>
					</div>
				{/if}
			</div>
		{/if}
		<div class="card-actions mt-auto items-center">
			<span
				class="tooltip"
				data-tip="{$LL.device.card_tooltip_last_status_change()}: {device.updated}"
			>
				{#if dateFnsLocale !== undefined}
					{formatDistance(parseISO(device.updated), now, {
						includeSeconds: true,
						addSuffix: true,
						locale: dateFnsLocale
					})}
				{/if}
			</span>
			{#if moreButtons.filter((btn) => btn.requires).length > 0}
				<div class="ms-auto flex flex-row flex-wrap gap-1">
					{#each moreButtons as btn}
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
			{$LL.device.modal_confirm_shutdown_title({ device: device.name })}
		</h3>
		<p class="py-4">{$LL.device.modal_confirm_shutdown_desc({ device: device.name })}</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">{$LL.buttons.cancel()}</button>
				<button class="btn btn-success" on:click={reboot}>{$LL.buttons.confirm()}</button>
			</form>
		</div>
	</div>
</dialog>
