<script lang="ts">
	import { browser } from '$app/environment';
	import DeviceCard from '$lib/components/DeviceCard.svelte';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import LL from '$lib/i18n/i18n-svelte';
	import { permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Group } from '$lib/types/device';
	import {
		faChevronCircleLeft,
		faChevronCircleRight,
		faPlus
	} from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

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
			dev?.expand?.groups?.forEach((group: Group) => {
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
							class="btn join-item"
							type="button"
							on:click={() => (orderByGroups = !orderByGroups)}
							>{$LL.home.order_groups()}
							<input type="checkbox" class="checkbox" checked={orderByGroups} />
						</button>
					</div>
				</div>
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
			<div class="tooltip join-item" data-tip={$LL.home.order_tooltip()}>
				<button
					class="btn {orderExpanded ? '' : 'rounded-field'}"
					on:click={() => (orderExpanded = !orderExpanded)}
				>
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
						<h1 class="mb-3 text-2xl font-bold">{group}</h1>
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
	<div class="flex justify-center">
		<div role="alert" class="alert alert-horizontal w-fit">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-6 w-6 shrink-0 stroke-current"
				fill="none"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
				/>
			</svg>
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
					Please ask the admin to grant you permissions to existing devices or to create new ones.
				</span>
			{/if}
		</div>
	</div>
{/if}
