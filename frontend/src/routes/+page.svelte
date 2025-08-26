<script lang="ts">
	import { browser } from '$app/environment';
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import { m } from '$lib/paraglide/messages';
	import { localeStore } from '$lib/stores/locale';
	import { backendUrl, permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Group } from '$lib/types/device';
	import {
		faChevronCircleLeft,
		faChevronCircleRight,
		faPlus,
		faPowerOff,
		faWarning
	} from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let loading = true;
	let devices: Device[] = [];
	let orderBy: 'name' | 'ip' = 'name';
	let orderExpanded = false;
	let orderByGroups = true;
	let searchQuery = '';
	let searchInput: HTMLInputElement;
	const isMac = /Mac|iPhone|iPad/.test(navigator.userAgent);
	const gridClass = 'grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4';

	if (browser) {
		const storedOrderBy = localStorage.getItem('orderBy');
		if (storedOrderBy === 'name' || storedOrderBy === 'ip') {
			orderBy = storedOrderBy;
		}
		orderByGroups = localStorage.getItem('orderByGroups') !== 'false';
	}

	$: if (browser) {
		localStorage.setItem('orderBy', orderBy);
	}
	$: if (browser) {
		localStorage.setItem('orderByGroups', String(orderByGroups));
	}

	$: filteredDevices = devices.filter(
		(dev) =>
			dev.ip.includes(searchQuery.toLowerCase()) ||
			dev.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			dev.mac.toLowerCase().includes(searchQuery.toLowerCase()) ||
			dev.description.toLowerCase().includes(searchQuery.toLowerCase())
	);

	$: devicesWithoutGroups = filteredDevices.filter((dev) => dev.groups.length === 0);

	$: devicesWithGroup = filteredDevices.reduce(
		(groups, dev) => {
			dev.expand?.groups?.forEach((group: Group) => {
				groups[group.id] = [...(groups[group.id] || []), dev];
			});
			return groups;
		},
		{} as Record<string, Device[]>
	);

	function getAllDevices() {
		$pocketbase
			.collection('devices')
			.getFullList(-1, { sort: 'name', expand: 'ports,groups' })
			.then((resp) => (devices = resp as Device[]))
			.catch((err) => toast.error(err.message))
			.finally(() => (loading = false));
	}

	function handleShortcut(event: KeyboardEvent) {
		if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
			event.preventDefault();
			searchInput.focus();
		}
	}

	function wakeGroup(group: string) {
		fetch(`${backendUrl}api/upsnap/wakegroup/${group}`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		}).catch((err) => {
			toast.error(err.message);
		});
	}

	onMount(() => {
		orderBy = (localStorage.getItem('orderBy') as 'name' | 'ip') || 'name';
		orderByGroups = localStorage.getItem('orderByGroups') !== 'false';
		getAllDevices();
		['devices', 'ports', 'permissions'].forEach((collection) =>
			$pocketbase
				.collection(collection)
				.subscribe('*', getAllDevices)
				.then(() => {
					console.log('Subscribed to changes in collection:', collection);
				})
		);
	});
</script>

<svelte:document on:keydown={handleShortcut} />

{#if loading}
	<PageLoading />
{:else if devices.length > 0}
	<div class="mb-4 flex flex-col justify-between gap-4 md:flex-row">
		<label class="input max-md:w-full">
			<svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
				><g
					stroke-linejoin="round"
					stroke-linecap="round"
					stroke-width="2.5"
					fill="none"
					stroke="currentColor"
					><circle cx="11" cy="11" r="8"></circle><path d="m21 21-4.3-4.3"></path></g
				></svg
			>
			<input
				type="search"
				class="grow"
				placeholder={m.home_search_placeholder()}
				bind:value={searchQuery}
				bind:this={searchInput}
			/>
			<kbd class="kbd kbd-sm">{isMac ? '⌘' : 'CTRL'}</kbd>
			<kbd class="kbd kbd-sm">K</kbd>
		</label>
		<div class="join ms-auto">
			{#if orderExpanded}
				<button class="btn join-item !outline-0" on:click={() => (orderByGroups = !orderByGroups)}>
					{m.home_order_groups()}
					<input type="checkbox" class="checkbox" checked={orderByGroups} />
				</button>
				<input
					class="btn join-item !outline-0"
					type="radio"
					name="order"
					aria-label={m.home_order_name()}
					bind:group={orderBy}
					value="name"
				/>
				<input
					class="btn join-item !outline-0"
					type="radio"
					name="order"
					aria-label={m.home_order_ip()}
					bind:group={orderBy}
					value="ip"
				/>
			{/if}
			<button
				class="btn join-item tooltip {orderExpanded ? '' : 'rounded-field'}"
				data-tip={m.home_order_tooltip()}
				on:click={() => (orderExpanded = !orderExpanded)}
			>
				<Fa icon={orderExpanded ? faChevronCircleRight : faChevronCircleLeft} size="lg" />
			</button>
		</div>
	</div>

	{#if orderByGroups}
		<div class="space-y-6">
			{#if devicesWithoutGroups.length > 0}
				<div class={gridClass}>
					{#each structuredClone(devicesWithoutGroups).sort( (a, b) => a[orderBy].localeCompare( b[orderBy], $localeStore, { numeric: true } ) ) as device (device.id)}
						<DeviceCard {device} />
					{/each}
				</div>
			{/if}
			{#each Object.entries(devicesWithGroup).sort( ([a], [b]) => a.localeCompare( b, $localeStore, { numeric: true } ) ) as [group, groupDevices], i (i)}
				<div>
					<h1 class="mb-3 text-2xl font-bold">
						{groupDevices[0].expand.groups.find((grp) => grp.id === group)?.name ||
							'Unknown group name'}
						<button class="btn btn-sm btn-success btn-soft" on:click={() => wakeGroup(group)}
							><Fa icon={faPowerOff} /> {m.home_wake_group()}</button
						>
					</h1>
					<div class={gridClass}>
						{#each groupDevices.sort( (a, b) => a[orderBy].localeCompare( b[orderBy], $localeStore, { numeric: true } ) ) as device (device.id)}
							<DeviceCard {device} />
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<div class={gridClass}>
			{#each structuredClone(filteredDevices).sort( (a, b) => a[orderBy].localeCompare( b[orderBy], $localeStore, { numeric: true } ) ) as device (device.id)}
				<DeviceCard {device} />
			{/each}
		</div>
	{/if}
{:else}
	<div class="flex justify-center">
		<div role="alert" class="alert alert-horizontal w-fit">
			<Fa icon={faWarning} size="lg" />
			{#if $pocketbase.authStore.isSuperuser || $permission.create}
				<div>
					<h3 class="font-bold">{m.home_no_devices()}</h3>
					<div class="text-xs">{m.home_add_first_device()}</div>
				</div>
				<a href="/device/new" class="btn btn-sm"><Fa icon={faPlus} />{m.home_add_first_device()}</a>
			{:else}
				<span>
					{m.home_grant_permissions()}
				</span>
			{/if}
		</div>
	</div>
{/if}
