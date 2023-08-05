<script lang="ts">
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import NetworkScan from '$lib/components/NetworkScan.svelte';
	import { faBinoculars, faWrench } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';
	import type { Device, Port } from '$lib/types/device';

	let device: Device = {
		expand: {
			ports: [] as Port[]
		}
	} as Device;

	let tabs = [
		{
			name: 'manual',
			icon: faWrench
		},
		{
			name: 'network scan',
			icon: faBinoculars
		}
	];
	let activeTab = 'manual';
</script>

<h1 class="text-3xl font-bold mb-8">New device</h1>
<div class="flex justify-center mb-6">
	<ul class="menu menu-horizontal bg-base-300 rounded-box">
		{#each tabs as tab}
			<li>
				<button
					class="capitalize"
					class:active={activeTab === tab.name}
					on:click={() => (activeTab = tab.name)}
					>{tab.name} <Fa icon={tab.icon} class="ms-2" /></button
				>
			</li>
		{/each}
	</ul>
</div>

{#if activeTab === 'manual'}
	<DeviceForm {device} />
{:else}
	<NetworkScan />
{/if}
