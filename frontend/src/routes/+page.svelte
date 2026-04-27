<script lang="ts">
	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
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

	let loading = $state(true);
	let devices: Device[] = $state([]);
	let orderBy: 'name' | 'ip' = $state('name');
	let orderExpanded = $state(false);
	let orderByGroups = $state(true);
	let searchQuery = $state('');
	let searchInput: HTMLInputElement | undefined = $state();
	let modalShutdownGroup: HTMLDialogElement | undefined = $state();
	let pendingShutdownGroup: Group | null = $state(null);
	let filteredDevices: Device[] = $derived([]);
	let devicesWithoutGroups: Device[] = $derived([]);
	let devicesWithGroup: {
		group: Group | null;
		devices: Device[];
	}[] = $derived([]);

	const isMac = browser && /Mac|iPhone|iPad/.test(navigator.userAgent);
	const gridClass = 'grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4';

	if (browser) {
		const storedOrderBy = localStorage.getItem('orderBy');
		if (storedOrderBy === 'name' || storedOrderBy === 'ip') {
			orderBy = storedOrderBy;
		}
		orderByGroups = localStorage.getItem('orderByGroups') !== 'false';
	}

	$effect(() => {
		if (browser) {
			localStorage.setItem('orderBy', orderBy);
		}
	});
	$effect(() => {
		if (browser) {
			localStorage.setItem('orderByGroups', String(orderByGroups));
		}
	});
	$effect(() => {
		filteredDevices = devices.filter(
			(dev) =>
				dev.ip.includes(searchQuery.toLowerCase()) ||
				dev.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				dev.mac.toLowerCase().includes(searchQuery.toLowerCase()) ||
				dev.description.toLowerCase().includes(searchQuery.toLowerCase())
		);
	});
	$effect(() => {
		devicesWithoutGroups = filteredDevices
			.filter((dev) => dev.groups.length === 0)
			.sort((a, b) => a[orderBy].localeCompare(b[orderBy], $localeStore, { numeric: true }));
	});
	$effect(() => {
		const groups: Record<string, { group: Group; devices: Device[] }> = {};
		filteredDevices.forEach((dev) => {
			if (dev.expand?.groups?.length) {
				for (const g of dev.expand.groups) {
					if (!groups[g.id]) groups[g.id] = { group: g, devices: [] };
					groups[g.id].devices.push(dev);
				}
			}
		});
		devicesWithGroup = Object.values(groups).map((g) => ({
			group: g.group,
			devices: g.devices.sort((a, b) =>
				a[orderBy].localeCompare(b[orderBy], $localeStore, { numeric: true })
			)
		}));
	});

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
			searchInput?.focus();
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

	function askShutdownGroupConfirmation(group: Group) {
		pendingShutdownGroup = group;
		modalShutdownGroup?.showModal();
	}

	function executeShutdownGroup() {
		if (pendingShutdownGroup) {
			fetch(`${backendUrl}api/upsnap/shutdowngroup/${pendingShutdownGroup.id}`, {
				headers: {
					Authorization: $pocketbase.authStore.token
				}
			}).catch((err) => {
				toast.error(err.message);
			});
		}
	}

	onMount(() => {
		orderBy = (localStorage.getItem('orderBy') as 'name' | 'ip') || 'name';
		orderByGroups = localStorage.getItem('orderByGroups') !== 'false';
		getAllDevices();
		['devices', 'ports', 'permissions'].forEach((collection) =>
			$pocketbase.collection(collection).subscribe('*', getAllDevices)
		);
	});
</script>

<svelte:document onkeydown={handleShortcut} />

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
				<button class="btn join-item outline-0!" onclick={() => (orderByGroups = !orderByGroups)}>
					{m.home_order_groups()}
					<input type="checkbox" class="checkbox" checked={orderByGroups} />
				</button>
				<input
					class="btn join-item outline-0!"
					type="radio"
					name="order"
					aria-label={m.home_order_name()}
					bind:group={orderBy}
					value="name"
				/>
				<input
					class="btn join-item outline-0!"
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
				onclick={() => (orderExpanded = !orderExpanded)}
			>
				<Fa icon={orderExpanded ? faChevronCircleRight : faChevronCircleLeft} size="lg" />
			</button>
		</div>
	</div>
	{#if orderByGroups}
		<div class="space-y-6">
			{#if devicesWithoutGroups.length > 0}
				<div class={gridClass}>
					{#each devicesWithoutGroups as device, i (device.id)}
						<DeviceCard bind:device={devicesWithoutGroups[i]} />
					{/each}
				</div>
			{/if}
			{#each devicesWithGroup.sort( (a, b) => (a.group?.name ?? '').localeCompare( b.group?.name ?? '', $localeStore, { numeric: true } ) ) as { group, devices } (group?.id)}
				<div>
					<h1 class="mb-3 flex flex-wrap items-center gap-3 text-2xl font-bold">
						<span>{group?.name ?? 'Unknown group name'}</span>
						<button
							class="btn btn-sm btn-success btn-soft"
							onclick={() => wakeGroup(group?.id ?? '')}
							><Fa icon={faPowerOff} /> {m.home_wake_group()}</button
						>
						<button
							class="btn btn-sm btn-error btn-soft"
							onclick={() => group && askShutdownGroupConfirmation(group)}
							><Fa icon={faPowerOff} /> {m.home_shutdown_group()}</button
						>
					</h1>
					<div class={gridClass}>
						{#each devices as device, i (device.id)}
							<DeviceCard bind:device={devices[i]} />
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{:else}
		{@const sorted = $state
			.snapshot(filteredDevices)
			.sort((a, b) => a[orderBy].localeCompare(b[orderBy], $localeStore, { numeric: true }))}
		<div class={gridClass}>
			{#each sorted as device, i (device.id)}
				<DeviceCard bind:device={sorted[i]} />
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
				<a href={resolve('/device/new')} class="btn btn-sm"
					><Fa icon={faPlus} />{m.home_add_first_device()}</a
				>
			{:else}
				<span>
					{m.home_grant_permissions()}
				</span>
			{/if}
		</div>
	</div>
{/if}

<dialog class="modal" bind:this={modalShutdownGroup}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">
			{m.home_modal_confirm_shutdown_group_title({ group: pendingShutdownGroup?.name ?? '' })}
		</h3>
		<p class="py-4">
			{m.home_modal_confirm_shutdown_group_desc({ group: pendingShutdownGroup?.name ?? '' })}
		</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">{m.buttons_cancel()}</button>
				<button class="btn btn-error" onclick={executeShutdownGroup}>{m.buttons_confirm()}</button>
			</form>
		</div>
	</div>
</dialog>
