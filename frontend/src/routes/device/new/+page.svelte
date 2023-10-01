<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import NetworkScan from '$lib/components/NetworkScan.svelte';
	import LL from '$lib/i18n/i18n-svelte';
	import { permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Port } from '$lib/types/device';
	import { faBinoculars, faWrench } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let device: Device = {
		expand: {
			ports: [] as Port[]
		},
		groups: [] as string[]
	} as Device;

	$: if (Object.hasOwn($permission, 'create')) {
		if (!$pocketbase.authStore.isAdmin && !$permission.create) {
			toast($LL.toasts.no_permission({ url: $page.url.pathname }), {
				icon: '⛔'
			});
			goto('/');
		}
	}

	let tabs = [
		{
			name: 'manual',
			ll_name: $LL.device.tabs[0](),
			icon: faWrench,
			show: true
		},
		{
			name: 'scan',
			ll_name: $LL.device.tabs[1](),
			icon: faBinoculars,
			show: $pocketbase.authStore.isAdmin
		}
	];
	let activeTab = 'manual';
</script>

<h1 class="text-3xl font-bold mb-8">{$LL.device.page_title()}</h1>
<div class="flex justify-center mb-6">
	<ul class="menu menu-horizontal bg-base-300 rounded-box gap-1">
		{#each tabs as tab}
			{#if tab.show}
				<li>
					<button class:active={activeTab === tab.name} on:click={() => (activeTab = tab.name)}
						>{tab.ll_name} <Fa icon={tab.icon} class="ms-2" /></button
					>
				</li>
			{/if}
		{/each}
	</ul>
</div>

{#if activeTab === 'manual'}
	<DeviceForm {device} />
{:else}
	<NetworkScan />
{/if}
