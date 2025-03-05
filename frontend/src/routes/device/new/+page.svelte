<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import NetworkScan from '$lib/components/NetworkScan.svelte';
	import { m } from '$lib/paraglide/messages';
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
		if (!$pocketbase.authStore.isSuperuser && !$permission.create) {
			toast(m.toasts_no_permission({ url: $page.url.pathname }), {
				icon: '⛔'
			});
			goto('/');
		}
	}

	let tabs = [
		{
			name: 'manual',
			ll_name: m['device_tabs.0'](),
			icon: faWrench,
			show: true
		},
		{
			name: 'scan',
			ll_name: m['device_tabs.1'](),
			icon: faBinoculars,
			show: $pocketbase.authStore.isSuperuser
		}
	];
	let activeTab = 'manual';
</script>

<h1 class="mb-8 text-3xl font-bold">{m.device_page_title()}</h1>
<div class="mb-6 flex justify-center">
	<ul class="menu menu-horizontal rounded-box bg-base-200 gap-1">
		{#each tabs as tab}
			{#if tab.show}
				<li>
					<button class:menu-active={activeTab === tab.name} on:click={() => (activeTab = tab.name)}
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
