<script lang="ts">
	import { goto } from '$app/navigation';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import { m } from '$lib/paraglide/messages';
	import { localeStore } from '$lib/stores/locale';
	import { backendUrl, pocketbase } from '$lib/stores/pocketbase';
	import { settingsPriv } from '$lib/stores/settings';
	import type { Device } from '$lib/types/device';
	import type { ScanResponse, ScannedDevice } from '$lib/types/scan';
	import type { SettingsPrivate } from '$lib/types/settings';
	import { faMagnifyingGlass, faPlus, faX } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

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
					const settings = res as SettingsPrivate;
					settingsPriv.set(settings);
					scanRange = settings.scan_range;
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
				toast.success(m.device_network_scan_range_saved());
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
				toast.success(m.toasts_device_created({ device: device.name }));
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
		toast.success(m.toasts_devices_created_multiple({ count: count }));
		goto('/');
	}
</script>

{#if $settingsPriv}
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m['device_tabs.1']()}</h2>
			<p class="my-2">
				{m.device_network_scan_desc()}
			</p>
			<div class="flex flex-row flex-wrap items-end gap-4">
				<form
					onsubmit={(e) => {
						e.preventDefault();
						saveSettings();
					}}
				>
					<fieldset class="fieldset p-0">
						<label class="floating-label mt-2">
							<span>{m.device_network_scan_ip_range()}</span>
							<div class="join max-w-xs">
								<input
									id="scan-range"
									class="input join-item w-full"
									type="text"
									placeholder="192.168.1.0/24"
									bind:value={scanRange}
								/>
								<button class="btn btn-neutral join-item" type="submit">{m.buttons_save()}</button>
							</div>
						</label>
					</fieldset>
				</form>
				<div>
					<div>
						{#if !$settingsPriv.scan_range}
							<button class="btn btn-error" disabled>
								<Fa icon={faX} />
								{m.device_network_scan_no_range()}
							</button>
						{:else if scanRange !== $settingsPriv.scan_range}
							<button class="btn btn-error" disabled>
								<Fa icon={faX} />
								{m.device_network_scan_unsaved_changes()}
							</button>
						{:else if scanRunning}
							<button class="btn no-animation">
								<span class="loading loading-spinner"></span>
								{m.device_network_scan_running()}
							</button>
						{:else}
							<button class="btn btn-success" onclick={() => scan()}>
								<Fa icon={faMagnifyingGlass} />
								{m.device_network_scan()}
							</button>
						{/if}
					</div>
				</div>
			</div>
			{#if scanResponse.devices?.length > 0}
				{#each scanResponse.devices.sort( (a, b) => a.ip.localeCompare( b.ip, $localeStore, { numeric: true } ) ) as device, index (device.mac)}
					<div class="collapse-arrow bg-base-100 collapse">
						<input type="radio" name="scanned-devices" checked={index === 0} />
						<div class="collapse-title font-bold">
							{device.name} <span class="badge">{device.ip}</span>
						</div>
						<div class="collapse-content">
							<div class="flex flex-row flex-wrap gap-4">
								<div>
									<strong>{m.device_network_scan_ip()}</strong><br />
									{device.ip}
								</div>
								<div>
									<strong>{m.device_network_scan_mac()}</strong><br />
									{device.mac}
								</div>
								<div>
									<strong>{m.device_network_scan_mac_vendor()}</strong><br />
									{device.mac_vendor}
								</div>

								<div>
									<strong>{m.device_network_scan_netmask()}</strong><br />
									{scanResponse.netmask}
								</div>
								<div class="ms-auto">
									<button
										class="btn btn-success btn-sm"
										onclick={(e) => {
											addSingle(device);
											e.currentTarget.disabled = true;
										}}><Fa icon={faPlus} />{m.buttons_add()}</button
									>
								</div>
							</div>
						</div>
					</div>
				{/each}
				<h2 class="card-title mt-4">{m.device_network_scan_add_all()}</h2>
				<div class="max-w-fit">
					<label class="label cursor-pointer">
						<input type="checkbox" class="checkbox" bind:checked={replaceNetmaskCheckbox} />
						<span class="ms-2 text-wrap break-words">{m.device_network_scan_replace_netmask()}</span
						>
					</label>
				</div>
				{#if replaceNetmaskCheckbox}
					<div class="max-w-fit">
						<label class="label cursor-pointer" for="replaceNetmaskInput">
							<span class="ms-2">{m.device_network_scan_new_netmask()}</span>
						</label>
						<input
							id="replaceNetmaskInput"
							class="input"
							type="text"
							placeholder="255.255.255.0"
							bind:value={replaceNetmask}
						/>
					</div>
				{/if}
				{#if scanResponse.devices.find((dev) => dev.name === 'Unknown')}
					<div class="max-w-fit">
						<label class="label cursor-pointer">
							<input type="checkbox" class="checkbox" bind:checked={addAllCheckbox} />
							<span class="ms-2 text-wrap break-words"
								>{m.device_network_scan_include_unknown()}</span
							>
						</label>
						{#if addAllCheckbox}
							<button
								class="btn btn-success"
								onclick={() => addAll()}
								disabled={scanResponse.devices.length === 0}
							>
								<Fa icon={faPlus} />
								{m.device_network_scan_add_all()} ({scanResponse.devices.length})
							</button>
						{:else}
							<button
								class="btn btn-success"
								onclick={() => addAll()}
								disabled={scanResponse.devices.filter((dev) => dev.name !== 'Unknown').length === 0}
							>
								<Fa icon={faPlus} />
								{m.device_network_scan_add_all()} ({scanResponse.devices.filter(
									(dev) => dev.name !== 'Unknown'
								).length})
							</button>
						{/if}
					</div>
				{:else}
					<div class="max-w-fit">
						<button
							class="btn btn-success"
							onclick={() => addAll()}
							disabled={scanResponse.devices.length === 0}
						>
							<Fa icon={faPlus} />
							{m.device_network_scan_add_all()} ({scanResponse.devices.length})
						</button>
					</div>
				{/if}
			{/if}
		</div>
	</div>
{:else}
	<PageLoading />
{/if}
