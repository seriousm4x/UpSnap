<script lang="ts">
	import { onMount } from 'svelte';
	import { pocketbase, permission } from '$lib/stores/pocketbase';
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import toast from 'svelte-french-toast';
	import Fa from 'svelte-fa';
	import {
		faChevronCircleLeft,
		faChevronCircleRight,
		faPlus
	} from '@fortawesome/free-solid-svg-icons';
	import { browser } from '$app/environment';
	import type { Device } from '$lib/types/device';

	let loading = true;
	let devices = [] as Device[];
	let devicesWithGroup: {
		[key: string]: Device[];
	};
	let devicesWithoutGroups = [] as Device[];
	let orderBy: 'name' | 'ip';
	let orderExpanded = false;
	let orderByGroups: boolean;
	const gridClass =
		'grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4';

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

	$: if (browser && orderBy !== undefined) {
		localStorage.setItem('orderBy', orderBy);
	}
	$: if (browser && orderByGroups !== undefined) {
		localStorage.setItem('orderByGroups', orderByGroups.toString());
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
		// get soring from local storage
		const lsO = localStorage.getItem('orderBy');
		if (lsO === 'name' || lsO === 'ip') {
			orderBy = lsO;
		} else {
			orderBy = 'name';
		}
		if (localStorage.getItem('orderByGroups') === null) {
			orderByGroups = true;
		} else {
			orderByGroups = localStorage.getItem('orderByGroups') === 'true' ? true : false;
		}

		// get collections and subscribe to changes
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
			{#if orderExpanded}
				<div class="join-item">
					<div class="join">
						<button
							class="join-item btn"
							type="button"
							on:click={() => (orderByGroups = !orderByGroups)}
							>Groups
							<input
								type="checkbox"
								class="checkbox checked:checkbox-primary"
								checked={orderByGroups}
							/>
						</button>
					</div>
				</div>
				<input
					class="join-item btn"
					type="radio"
					name="order"
					aria-label="Name"
					bind:group={orderBy}
					value="name"
				/>
				<input
					class="join-item btn"
					type="radio"
					name="order"
					aria-label="IP"
					bind:group={orderBy}
					value="ip"
				/>
			{/if}
			<div class="tooltip" data-tip="Order">
				<button class="join-item btn" on:click={() => (orderExpanded = !orderExpanded)}>
					{#if orderExpanded}
						<Fa icon={faChevronCircleRight} size="lg" />
					{:else}
						<Fa icon={faChevronCircleLeft} size="lg" />
					{/if}
				</button>
			</div>
		</div>
	</div>
	{#if orderByGroups}
		<div class="space-y-6">
			{#if devicesWithoutGroups.length > 0}
				<div class={gridClass}>
					{#each devicesWithoutGroups.sort( (a, b) => a[orderBy].localeCompare( b[orderBy], undefined, { numeric: true, sensitivity: 'base' } ) ) as device}
						<DeviceCard {device} />
					{/each}
				</div>
			{/if}
			{#if Object.keys(devicesWithGroup).length > 0}
				{#each Object.keys(devicesWithGroup).sort( (a, b) => a.localeCompare( b, undefined, { numeric: true, sensitivity: 'base' } ) ) as group}
					<div>
						<h1 class="text-2xl font-bold mb-3">{group}</h1>
						<div class={gridClass}>
							{#each devicesWithGroup[group].sort( (a, b) => a[orderBy].localeCompare( b[orderBy], undefined, { numeric: true, sensitivity: 'base' } ) ) as device}
								<DeviceCard {device} />
							{/each}
						</div>
					</div>
				{/each}
			{/if}
		</div>
	{:else}
		<div class={gridClass}>
			{#each devices.sort( (a, b) => a[orderBy].localeCompare( b[orderBy], undefined, { numeric: true, sensitivity: 'base' } ) ) as device}
				<DeviceCard {device} />
			{/each}
		</div>
	{/if}
{:else}
	<div class="container text-center">
		<p>No devices here.</p>
		{#if $pocketbase.authStore.isAdmin || $permission.create}
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
