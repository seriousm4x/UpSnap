<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import NetworkScan from '$lib/components/NetworkScan.svelte';
	import { m } from '$lib/paraglide/messages';
	import { permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Group, Port } from '$lib/types/device';
	import { faBinoculars, faWrench } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let device: Device = {
		name: '',
		ip: '',
		mac: '',
		netmask: '',
		description: '',
		status: '',
		ports: [] as string[],
		link_open: '',
		ping_cmd: '',
		wake_cron: '',
		wake_cron_enabled: false,
		wake_cmd: '',
		wake_confirm: false,
		wake_timeout: 0,
		shutdown_cron: '',
		shutdown_cron_enabled: false,
		shutdown_cmd: '',
		shutdown_confirm: false,
		shutdown_timeout: 0,
		password: '',
		groups: [] as string[],
		expand: {
			ports: [] as Port[],
			groups: [] as Group[]
		},
		created_by: '',
		sol_enabled: false,
		sol_auth: false,
		sol_user: '',
		sol_password: '',
		sol_port: 0
	} as Device;
	let activeTab = $state('manual');

	$effect(() => {
		if (Object.hasOwn($permission, 'create')) {
			if (!$pocketbase.authStore.isSuperuser && !$permission.create) {
				toast(m.toasts_no_permission({ url: page.url.pathname }), {
					icon: '⛔'
				});
				goto(resolve('/'));
			}
		}
	});

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
</script>

<h1 class="mb-8 text-3xl font-bold">{m.device_page_title()}</h1>
<div class="mb-6 flex justify-center">
	<ul class="menu menu-horizontal rounded-box bg-base-200 gap-1">
		{#each tabs as tab (tab)}
			{#if tab.show}
				<li>
					<button class:menu-active={activeTab === tab.name} onclick={() => (activeTab = tab.name)}
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
