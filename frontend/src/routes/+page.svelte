<script lang="ts">
	import { onMount } from 'svelte';
	import { pocketbase, permission, isAdmin } from '$lib/stores/pocketbase';
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import toast from 'svelte-french-toast';
	import Fa from 'svelte-fa';
	import { faPlus } from '@fortawesome/free-solid-svg-icons';
	import { browser } from '$app/environment';
	import type { Device } from '$lib/types/device';

	let devices = [] as Device[];
	let loading = true;
	let order: string;
	let devicesWithGroup: {
		[key: string]: Device[];
	};
	let devicesWithoutGroups = [] as Device[];

	$: {
		devicesWithoutGroups = [];
		devicesWithoutGroups = devices.filter((dev) => dev.groups.length === 0);
	}

	$: {
		devicesWithGroup = {};
		devices.forEach((dev) => {
			dev?.expand?.groups?.forEach((group) => {
				if (!devicesWithGroup[group.name]) {
					devicesWithGroup[group.name] = [];
				}
				devicesWithGroup[group.name] = [...devicesWithGroup[group.name], dev];
			});
		});
	}

	$: if (browser) {
		if (order === undefined) {
			order = localStorage.getItem('order') || 'group';
		}
		localStorage.setItem('order', order || 'group');
	}

	function getAllDevices() {
		$pocketbase
			.collection('devices')
			.getFullList(1000, { sort: 'name', expand: 'ports,groups' })
			.then((resp) => {
				devices = resp as Device[];
			})
			.catch((err) => {
				toast.error(err.message);
			})
			.finally(() => {
				loading = false;
			});
	}

	onMount(() => {
		getAllDevices();

		$pocketbase.collection('devices').subscribe('*', () => {
			getAllDevices();
		});

		$pocketbase.collection('ports').subscribe('*', () => {
			getAllDevices();
		});

		$pocketbase.collection('permissions').subscribe('*', () => {
			getAllDevices();
		});
	});
</script>

{#if loading}
	<PageLoading />
{:else if devices.length > 0}
	<div class="flex justify-end">
		<div class="join mb-4">
			<input
				class="join-item btn"
				type="radio"
				name="order"
				aria-label="Groups"
				bind:group={order}
				value="group"
			/>
			<input
				class="join-item btn"
				type="radio"
				name="order"
				aria-label="Name"
				bind:group={order}
				value="name"
			/>
			<input
				class="join-item btn"
				type="radio"
				name="order"
				aria-label="IP"
				bind:group={order}
				value="ip"
			/>
		</div>
	</div>
	{#if order === 'group'}
		<div
			class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4"
		>
			{#if devicesWithoutGroups.length > 0}
				{#each devicesWithoutGroups.sort( (a, b) => a.name.localeCompare( b.name, undefined, { numeric: true, sensitivity: 'base' } ) ) as device}
					<DeviceCard {device} />
				{/each}
			{/if}
		</div>
		{#each Object.keys(devicesWithGroup).sort( (a, b) => a.localeCompare( b, undefined, { numeric: true, sensitivity: 'base' } ) ) as group}
			<h1 class="mt-6 mb-4 text-2xl font-bold">{group}</h1>
			<div
				class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4"
			>
				{#each devicesWithGroup[group] as device}
					<DeviceCard {device} />
				{/each}
			</div>
		{/each}
	{:else if order === 'name'}
		<div
			class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4"
		>
			{#each devices.sort( (a, b) => a.name.localeCompare( b.name, undefined, { numeric: true, sensitivity: 'base' } ) ) as device}
				<DeviceCard {device} />
			{/each}
		</div>
	{:else if order === 'ip'}
		<div
			class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4"
		>
			{#each devices.sort( (a, b) => a.ip.localeCompare( b.ip, undefined, { numeric: true } ) ) as device}
				<DeviceCard {device} />
			{/each}
		</div>
	{/if}
{:else}
	<div class="container text-center">
		<p>No devices here.</p>
		{#if $isAdmin || $permission.create}
			<p>
				<a href="/device/new" class="btn btn-ghost"
					><Fa icon={faPlus} class="ms-2" />Add your first device
				</a>
			</p>
		{:else}
			<p>
				Please ask the admin to grant you permissions to existing devices or to create new ones.
			</p>
		{/if}
	</div>
{/if}
