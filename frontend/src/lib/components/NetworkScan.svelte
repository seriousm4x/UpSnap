<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
	import { settingsPriv } from '$lib/stores/settings';
	import Alert from '$lib/components/Alert.svelte';
	import Fa from 'svelte-fa';
	import {
		faMagnifyingGlass,
		faPlus,
		faTriangleExclamation,
		faX
	} from '@fortawesome/free-solid-svg-icons';
	import type { SettingsPrivate } from '$lib/types/settings';
	import type { ScanResponse, ScannedDevice } from '$lib/types/scan';

	let saveErrMsg = '';
	let saveErrTimeout: number;
	let scanErrMsg = '';
	let scanErrTimeout: number;
	let scanRange = '';
	let scanRunning = false;
	let scanResponse: ScanResponse = {
		netmask: '',
		devices: []
	};
	let addAllCheckbox = true;

	onMount(() => {
		if (!$settingsPriv && $pocketbase.authStore.isValid) {
			$pocketbase
				.collection('settings_private')
				.getFirstListItem('')
				.then((res) => {
					settingsPriv.set(res as SettingsPrivate);
					scanRange = $settingsPriv.scan_range;
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
			})
			.catch((err) => {
				clearTimeout(saveErrTimeout);
				saveErrTimeout = setTimeout(() => {
					saveErrMsg = '';
				}, 10000);
				saveErrMsg = err;
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
				clearTimeout(scanErrTimeout);
				scanErrTimeout = setTimeout(() => {
					scanErrMsg = '';
				}, 10000);
				scanErrMsg = err.message;
			})
			.finally(() => (scanRunning = false));
	}

	function createDevice(device: ScannedDevice) {
		device.netmask = scanResponse.netmask;
		$pocketbase
			.collection('devices')
			.create(device)
			.catch((err) => {
				clearTimeout(saveErrTimeout);
				saveErrTimeout = setTimeout(() => {
					saveErrMsg = '';
				}, 10000);
				saveErrMsg = err;
			});
	}

	function addAll() {
		scanResponse.devices.forEach((device) => {
			if (!addAllCheckbox && device.name === 'Unknown') return;
			createDevice(device);
		});
		goto('/');
	}
</script>

{#if $settingsPriv}
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Network Scan</h2>
			<p class="my-2">
				Automatically scan your network for devices. For this to work, you need to run UpSnap as
				root/admin and have nmap installed and available in your $PATH (For docker users, thats
				already the case and you don't need to do anything). Scanning might take some seconds.
			</p>
			{#if saveErrMsg !== ''}
				<Alert
					color="error"
					message={saveErrMsg}
					icon={faTriangleExclamation}
					customClasses="mt-4 max-w-fit"
				/>
			{/if}
			<div class="flex flex-row flex-wrap gap-4 items-end">
				<div>
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
						<button class="btn join-item" on:click={() => saveSettings()}>Save</button>
					</div>
				</div>
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
			{#if scanErrMsg !== ''}
				<Alert
					color="error"
					message={scanErrMsg}
					icon={faTriangleExclamation}
					customClasses="mt-4 max-w-fit"
				/>
			{/if}
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
											createDevice(device);
											e.currentTarget.disabled = true;
										}}><Fa icon={faPlus} />Add</button
									>
								</div>
							</div>
						</div>
					</div>
				{/each}
				<div class="form-control max-w-fit flex flex-row flex-wrap gap-4 mt-4 ms-auto">
					<label class="label cursor-pointer">
						<input type="checkbox" class="checkbox" bind:checked={addAllCheckbox} />
						<span class="label-text ms-2">Include devices where name is "Unknown"</span>
					</label>
					<button
						class="btn btn-success"
						on:click={() => addAll()}
						disabled={!addAllCheckbox &&
							scanResponse.devices.filter((device) => {
								device.name !== 'Unknown';
							}).length === 0}
					>
						<Fa icon={faPlus} /> Add all ({addAllCheckbox
							? scanResponse.devices.length
							: scanResponse.devices.filter((device) => {
									device.name !== 'Unknown';
							  }).length})
					</button>
				</div>
			{/if}
		</div>
	</div>
{:else}
	<div class="container max-w-lg mx-auto text-center">
		<span class="loading loading-dots loading-lg" />
	</div>
{/if}
