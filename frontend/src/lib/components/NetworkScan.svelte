<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
	import { settingsPriv } from '$lib/stores/settings';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import Fa from 'svelte-fa';
	import { faMagnifyingGlass, faPlus, faX } from '@fortawesome/free-solid-svg-icons';
	import toast from 'svelte-french-toast';
	import type { SettingsPrivate } from '$lib/types/settings';
	import type { ScanResponse, ScannedDevice } from '$lib/types/scan';
	import type { Device } from '$lib/types/device';

	let scanRange = '';
	let scanRunning = false;
	let scanResponse: ScanResponse = {
		netmask: '',
		devices: []
	};
	let addAllCheckbox = true;
	let replaceNetmaskCheckbox = false;
	let replaceNetmask = '';

	onMount(() => {
		if (!$settingsPriv) {
			$pocketbase
				.collection('settings_private')
				.getFirstListItem('')
				.then((res) => {
					settingsPriv.set(res as SettingsPrivate);
					scanRange = $settingsPriv.scan_range;
				})
				.catch((err) => {
					toast.error(err.message);
				});
		} else {
			scanRange = $settingsPriv.scan_range;
		}
	});

	function saveSettings() {
		$pocketbase
			.collection('settings_private')
			.update($settingsPriv.id, {
				scan_range: scanRange
			})
			.then((res) => {
				settingsPriv.set(res as SettingsPrivate);
				toast.success('Scan range saved');
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function scan() {
		scanRunning = true;
		fetch(`${backendUrl}api/upsnap/scan`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then(async (resp) => {
				if (resp.ok) {
					return resp.json();
				} else {
					return Promise.reject(await resp.json());
				}
			})
			.then((data) => {
				scanResponse = data as ScanResponse;
			})
			.catch((err) => {
				toast.error(err);
			})
			.finally(() => (scanRunning = false));
	}

	async function createDevice(device: ScannedDevice): Promise<Device> {
		if (replaceNetmaskCheckbox) {
			device.netmask = replaceNetmask;
		} else {
			device.netmask = scanResponse.netmask;
		}
		return $pocketbase.collection('devices').create(device);
	}

	async function addSingle(device: ScannedDevice) {
		await createDevice(device)
			.then(() => {
				toast.success(`Added ${device.name}`);
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	async function addAll() {
		let count = 0;
		await Promise.all(
			scanResponse.devices.map(async (dev) => {
				if (!addAllCheckbox && dev.name === 'Unknown') return;
				await createDevice(dev)
					.catch((err) => {
						toast.error(err.message);
					})
					.then(() => {
						count += 1;
					});
			})
		);
		toast.success(`Added ${count} devices`);
		goto('/');
	}
</script>

{#if $settingsPriv}
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Network scan</h2>
			<p class="my-2">
				Automatically scan your network for devices. For this to work, you need to run UpSnap as
				root/admin and have nmap installed and available in your $PATH (For docker users, thats
				already the case and you don't need to do anything). Scanning might take some seconds.
			</p>
			<div class="flex flex-row flex-wrap gap-4 items-end">
				<form on:submit|preventDefault={saveSettings}>
					<label class="label" for="scan-range">
						<span class="label-text">IP range</span>
					</label>
					<div class="join">
						<input
							id="scan-range"
							class="input input-bordered join-item"
							type="text"
							placeholder="192.168.1.0/24"
							bind:value={scanRange}
						/>
						<button class="btn btn-primary join-item" type="submit">Save</button>
					</div>
				</form>
				<div>
					<div>
						{#if !$settingsPriv.scan_range}
							<button class="btn btn-error" disabled>
								<Fa icon={faX} /> No scan range
							</button>
						{:else if scanRange !== $settingsPriv.scan_range}
							<button class="btn btn-error" disabled>
								<Fa icon={faX} /> Unsaved changes
							</button>
						{:else if scanRunning}
							<button class="btn no-animation">
								<span class="loading loading-spinner" />
								Scan running
							</button>
						{:else}
							<button class="btn btn-success" on:click={() => scan()}>
								<Fa icon={faMagnifyingGlass} /> Scan
							</button>
						{/if}
					</div>
				</div>
			</div>
			{#if scanResponse.devices?.length > 0}
				{#each scanResponse.devices.sort( (a, b) => a.ip.localeCompare( b.ip, undefined, { numeric: true } ) ) as device, index}
					<div class="collapse collapse-arrow bg-base-200">
						<input type="radio" name="scanned-devices" checked={index === 0} />
						<div class="collapse-title font-bold">
							{device.name} <span class="badge">{device.ip}</span>
						</div>
						<div class="collapse-content">
							<div class="flex flex-row flex-wrap gap-4">
								<div>
									<strong>IP:</strong><br />
									{device.ip}
								</div>
								<div>
									<strong>Mac:</strong><br />
									{device.mac}
								</div>
								<div>
									<strong>Mac vendor:</strong><br />
									{device.mac_vendor}
								</div>

								<div>
									<strong>Netmask:</strong><br />
									{scanResponse.netmask}
								</div>
								<div class="ms-auto">
									<button
										class="btn btn-sm btn-success"
										on:click={(e) => {
											addSingle(device);
											e.currentTarget.disabled = true;
										}}><Fa icon={faPlus} />Add</button
									>
								</div>
							</div>
						</div>
					</div>
				{/each}
				<h2 class="card-title mt-4">Add all devices</h2>
				<div class="form-control max-w-fit">
					<label class="label cursor-pointer">
						<input type="checkbox" class="checkbox" bind:checked={replaceNetmaskCheckbox} />
						<span class="label-text ms-2">Replace netmask for all devices?</span>
					</label>
				</div>
				{#if replaceNetmaskCheckbox}
					<div class="form-control max-w-fit">
						<label class="label cursor-pointer" for="replaceNetmaskInput">
							<span class="label-text ms-2">New netmask</span>
						</label>
						<input
							id="replaceNetmaskInput"
							class="input input-bordered"
							type="text"
							placeholder="255.255.0.0"
							bind:value={replaceNetmask}
						/>
					</div>
				{/if}
				{#if scanResponse.devices.find((dev) => dev.name === 'Unknown')}
					<div class="form-control max-w-fit">
						<label class="label cursor-pointer">
							<input type="checkbox" class="checkbox" bind:checked={addAllCheckbox} />
							<span class="label-text ms-2">Include devices where name is "Unknown"</span>
						</label>
						{#if addAllCheckbox}
							<button
								class="btn btn-success"
								on:click={() => addAll()}
								disabled={scanResponse.devices.length === 0}
							>
								<Fa icon={faPlus} /> Add all ({scanResponse.devices.length})
							</button>
						{:else}
							<button
								class="btn btn-success"
								on:click={() => addAll()}
								disabled={scanResponse.devices.filter((dev) => dev.name !== 'Unknown').length === 0}
							>
								<Fa icon={faPlus} /> Add all ({scanResponse.devices.filter(
									(dev) => dev.name !== 'Unknown'
								).length})
							</button>
						{/if}
					</div>
				{:else}
					<div class="form-control max-w-fit">
						<button
							class="btn btn-success"
							on:click={() => addAll()}
							disabled={scanResponse.devices.length === 0}
						>
							<Fa icon={faPlus} /> Add all ({scanResponse.devices.length})
						</button>
					</div>
				{/if}
			{/if}
		</div>
	</div>
{:else}
	<PageLoading />
{/if}
