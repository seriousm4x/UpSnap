<script lang="ts">
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import LL from '$lib/i18n/i18n-svelte';
	import { permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Group } from '$lib/types/device';
	import {
		faChevronCircleLeft,
		faChevronCircleRight,
		faMagnifyingGlass,
		faPlus,
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

	const gridClass = 'grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4';

	const filteredDevices = () =>
		devices.filter(
			(dev) =>
				dev.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				dev.ip.includes(searchQuery.toLowerCase()) ||
				dev.description.toLowerCase().includes(searchQuery.toLowerCase())
		);

	const devicesWithoutGroups = () => filteredDevices().filter((dev) => dev.groups.length === 0);

	const devicesWithGroup = () => {
		return filteredDevices().reduce(
			(groups, dev) => {
				dev.expand?.groups?.forEach((group: Group) => {
					groups[group.name] = [...(groups[group.name] || []), dev];
				});
				return groups;
			},
			{} as Record<string, Device[]>
		);
	};

	function getAllDevices() {
		$pocketbase
			.collection('devices')
			.getFullList(-1, { sort: 'name', expand: 'ports,groups' })
			.then((resp) => (devices = resp as Device[]))
			.catch((err) => toast.error(err.message))
			.finally(() => (loading = false));
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

{#if loading}
	<PageLoading />
{:else if devices.length > 0}
	<div class="mb-4 flex justify-between">
		<label class="input">
			<Fa icon={faMagnifyingGlass} size="sm" />
			<input type="search" placeholder={$LL.home.search_placeholder()} bind:value={searchQuery} />
		</label>
		<div class="join">
			{#if orderExpanded}
				<button class="btn join-item" on:click={() => (orderByGroups = !orderByGroups)}>
					{$LL.home.order_groups()}
					<input type="checkbox" class="checkbox" checked={orderByGroups} />
				</button>
				<input
					class="btn join-item"
					type="radio"
					name="order"
					aria-label={$LL.home.order_name()}
					bind:group={orderBy}
					value="name"
				/>
				<input
					class="btn join-item"
					type="radio"
					name="order"
					aria-label={$LL.home.order_ip()}
					bind:group={orderBy}
					value="ip"
				/>
			{/if}
			<button
				class="btn join-item tooltip {orderExpanded ? '' : 'rounded-field'}"
				data-tip={$LL.home.order_tooltip()}
				on:click={() => (orderExpanded = !orderExpanded)}
			>
				<Fa icon={orderExpanded ? faChevronCircleRight : faChevronCircleLeft} size="lg" />
			</button>
		</div>
	</div>

	{#if orderByGroups}
		<div class="space-y-6">
			{#if devicesWithoutGroups().length > 0}
				<div class={gridClass}>
					{#each devicesWithoutGroups().sort( (a, b) => a[orderBy].localeCompare(b[orderBy]) ) as device}
						<DeviceCard {device} />
					{/each}
				</div>
			{/if}
			{#each Object.entries(devicesWithGroup()).sort( ([a], [b]) => a.localeCompare(b) ) as [group, groupDevices]}
				<div>
					<h1 class="mb-3 text-2xl font-bold">{group}</h1>
					<div class={gridClass}>
						{#each groupDevices.sort((a, b) => a[orderBy].localeCompare(b[orderBy])) as device}
							<DeviceCard {device} />
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<div class={gridClass}>
			{#each filteredDevices().sort((a, b) => a[orderBy].localeCompare(b[orderBy])) as device}
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
					<h3 class="font-bold">{$LL.home.no_devices()}</h3>
					<div class="text-xs">{$LL.home.add_first_device()}</div>
				</div>
				<a href="/device/new" class="btn btn-sm"
					><Fa icon={faPlus} />{$LL.home.add_first_device()}</a
				>
			{:else}
				<span>
					{$LL.home.grant_permissions()}
				</span>
			{/if}
		</div>
	</div>
{/if}
