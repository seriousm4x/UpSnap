<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import NetworkScan from '$lib/components/NetworkScan.svelte';
	import { permission, isAdmin } from '$lib/stores/pocketbase';
	import { faBinoculars, faWrench } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';
	import type { Device, Port } from '$lib/types/device';

	let device: Device = {
		expand: {
			ports: [] as Port[]
		},
		groups: [] as string[]
	} as Device;

	$: if (Object.hasOwn($permission, 'create')) {
		if (!$isAdmin && !$permission.create) {
			toast(`You don't have permission to visit ${$page.url.pathname}`, {
				icon: '⛔'
			});
			goto('/');
		}
	}

	let tabs = [
		{
			name: 'manual',
			icon: faWrench,
			show: true
		},
		{
			name: 'network scan',
			icon: faBinoculars,
			show: $isAdmin
		}
	];
	let activeTab = 'manual';
</script>

<h1 class="text-3xl font-bold mb-8">New device</h1>
<div class="flex justify-center mb-6">
	<ul class="menu menu-horizontal bg-base-300 rounded-box gap-1">
		{#each tabs as tab}
			{#if tab.show}
				<li>
					<button
						class="capitalize"
						class:active={activeTab === tab.name}
						on:click={() => (activeTab = tab.name)}
						>{tab.name} <Fa icon={tab.icon} class="ms-2" /></button
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
