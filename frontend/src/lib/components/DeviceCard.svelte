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
	import type { Locale } from 'date-fns';
	import { formatDistance, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale/de';
	import { enUS } from 'date-fns/locale/en-US';
	import { fr } from 'date-fns/locale/fr';
	import { pt } from 'date-fns/locale/pt';
	import { zhCN } from 'date-fns/locale/zh-CN';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';
	import { scale } from 'svelte/transition';
	import DeviceCardNic from './DeviceCardNic.svelte';
	import es from '../i18n/es';

	export let device: Device;

	let moreButtons = [
		{
			text: $LL.device.card_btn_more_edit(),
			icon: faPen,
			onClick: () => goto(`/device/${device.id}`),
			requires: $pocketbase.authStore.isAdmin || $permission.update?.includes(device.id)
		},
		{
			text: $LL.device.card_btn_more_sleep(),
			icon: faBed,
			onClick: () => sleep(),
			requires:
				($pocketbase.authStore.isAdmin || $permission.power?.includes(device.id)) &&
				device.status === 'online' &&
				device.sol_enabled
		},
		{
			text: $LL.device.card_btn_more_reboot(),
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
				case 'es':
				case 'es-ES':
					dateFnsLocale = es;
					break;
				case 'pt':
				case 'pt-PT':
					dateFnsLocale = pt;
					break;
				case 'zh':
				case 'zh-CN':
					dateFnsLocale = zhCN;
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
					<div class="tooltip" data-tip={$LL.device.card_tooltip_wake_cron()}>
						<span class="badge badge-success gap-1 p-3"
							><Fa icon={faCircleArrowUp} />{device.wake_cron}</span
						>
					</div>
				{/if}
				{#if device.shutdown_cron_enabled}
					<div class="tooltip" data-tip={$LL.device.card_tooltip_shutdown_cron()}>
						<span class="badge badge-error gap-1 p-3"
							><Fa icon={faCircleArrowDown} />{device.shutdown_cron}</span
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
				<div class="dropdown dropdown-top dropdown-end bg-base-300 ms-auto">
					<label tabindex="-1" class="btn btn-sm m-1" for="more-{device.id}"
						>{$LL.device.card_btn_more()}</label
					>
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
