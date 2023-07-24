<script lang="ts">
	import { page } from '$app/stores';
	import { pocketbase, devices } from '$lib/stores/pocketbase';
	import { faTriangleExclamation } from '@fortawesome/free-solid-svg-icons';
	import Alert from '$lib/components/Alert.svelte';
	import DeviceFormPort from '$lib/components/DeviceFormPort.svelte';
	import type { Device, Port } from '$lib/types/device';
	import { goto } from '$app/navigation';

	let waitingForDevices = true;
	let device: Device | undefined;
	let changes: Device | undefined;

	let portErrMsg = '';
	let portErrTimeout: number;

	let saveErrMsg = '';
	let saveErrTimeout: number;

	devices.subscribe((data) => {
		if (waitingForDevices === false) return;
		if (data.length === 0) return;

		device = data.find((dev) => dev.id === $page.params.id);
		if (device !== undefined) {
			changes = device;
			if (changes.expand.ports === undefined) {
				changes.expand.ports = [];
			}
		}
		waitingForDevices = false;
	});

	async function saveDevice() {
		if (device === undefined || changes === undefined) return;

		// create/update all ports
		let portIds: string[] = [];
		await Promise.all(
			changes.expand.ports.map(async (port) => {
				if (port.id === undefined) {
					const data = await createPort(port as Port);
					if (data === undefined) return;
					portIds = [...portIds, data.id];
				} else {
					const data = await updatePort(port as Port);
					if (data === undefined) return;
					portIds = [...portIds, data.id];
				}
			})
		);

		changes.ports = portIds;
		$pocketbase
			.collection('devices')
			.update(device.id, changes)
			.then(() => {
				goto('/');
			})
			.catch((err) => {
				clearTimeout(saveErrTimeout);
				saveErrTimeout = setTimeout(() => {
					saveErrMsg = '';
				}, 10000);
				saveErrMsg = err;
			});
	}

	async function createPort(port: Port): Promise<Port | undefined> {
		let newport: Port | undefined;
		await $pocketbase
			.collection('ports')
			.create(port)
			.then((data) => {
				newport = data as Port;
			})
			.catch((err) => {
				clearTimeout(saveErrTimeout);
				saveErrTimeout = setTimeout(() => {
					saveErrMsg = '';
				}, 10000);
				saveErrMsg = err;
				newport = undefined;
			});
		return newport;
	}

	async function updatePort(port: Port): Promise<Port | undefined> {
		let newport: Port | undefined;
		await $pocketbase
			.collection('ports')
			.update(port.id, port)
			.then((data) => {
				newport = data as Port;
			})
			.catch((err) => {
				clearTimeout(saveErrTimeout);
				saveErrTimeout = setTimeout(() => {
					saveErrMsg = '';
				}, 10000);
				saveErrMsg = err;
				newport = undefined;
			});
		return newport;
	}

	function createEmptyPort() {
		if (changes === undefined) return;
		changes.expand.ports = [...changes.expand.ports, { name: '', number: 1 } as Port];
	}
</script>

{#if waitingForDevices}
	<div class="container max-w-lg mx-auto text-center">
		<span class="loading loading-dots loading-lg" />
	</div>
{:else if device === undefined || changes === undefined}
	<div class="container max-w-lg mx-auto">
		<div class="alert alert-error">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="stroke-current shrink-0 h-6 w-6"
				fill="none"
				viewBox="0 0 24 24"
				><path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
				/></svg
			>
			<span><strong>No device with id:</strong> {$page.params.id}</span>
			<div>
				<a class="btn btn-sm" href="/">Ok</a>
			</div>
		</div>
	</div>
{:else}
	<h1 class="text-3xl font-bold mb-8">{device.name}</h1>
	<div class="card w-full bg-base-300 shadow-xl">
		<div class="card-body">
			<form on:submit|preventDefault={saveDevice}>
				<h2 class="card-title">General</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
					<div class="form-control w-full max-w-xs">
						<label class="label" for="device-name">
							<div class="label-text">
								<span>Name</span>
								<span class="text-error">*</span>
							</div>
						</label>
						<input
							id="device-name"
							type="text"
							placeholder="Device name"
							class="input w-full max-w-xs"
							bind:value={changes.name}
							required
						/>
					</div>
					<div class="form-control w-full max-w-xs">
						<label class="label" for="device-ip">
							<div class="label-text">
								<span>IP</span>
								<span class="text-error">*</span>
							</div>
						</label>
						<input
							id="device-ip"
							type="text"
							placeholder="192.168.0.5"
							class="input w-full max-w-xs"
							bind:value={changes.ip}
							required
						/>
					</div>
					<div class="form-control w-full max-w-xs">
						<label class="label" for="device-mac">
							<div class="label-text">
								<span>Mac</span>
								<span class="text-error">*</span>
							</div>
						</label>
						<input
							id="device-mac"
							type="text"
							placeholder="aa:bb:cc:dd:ee:ff"
							class="input w-full max-w-xs"
							bind:value={changes.mac}
							required
						/>
					</div>
					<div class="form-control w-full max-w-xs">
						<label class="label" for="device-netmask">
							<div class="label-text">
								<span>Netmask</span>
								<span class="text-error">*</span>
							</div>
						</label>
						<input
							id="device-netmask"
							type="text"
							placeholder="255.255.255.0"
							class="input w-full max-w-xs"
							bind:value={changes.netmask}
							required
						/>
					</div>
				</div>
				<h2 class="card-title mt-8">Ports</h2>
				<div class="form-control w-full">
					<p class="my-2">
						UpSnap can also check if given ports are open. You can define them below.
					</p>
					<div class="flex flex-wrap gap-4">
						{#each changes.expand.ports as _, index}
							<DeviceFormPort bind:changes {index} bind:portErrMsg bind:portErrTimeout />
						{/each}
					</div>
					{#if portErrMsg !== ''}
						<Alert
							color="error"
							icon={faTriangleExclamation}
							message={portErrMsg}
							customClasses="mt-4 max-w-fit"
							duration={10000}
						/>
					{/if}
					<button
						class="btn btn-wide btn-neutral mt-4"
						on:click={() => createEmptyPort()}
						type="button">Add new port</button
					>
				</div>
				<h2 class="card-title mt-8">Link</h2>
				<div class="form-control w-full">
					<p class="my-2">
						Makes your device name a clickable link, perfect for linking a dashboard for example.
					</p>
					<input
						type="url"
						placeholder="https:// ..."
						class="input w-full max-w-xs"
						bind:value={changes.link}
					/>
				</div>
				<div class="card-actions justify-end mt-4">
					<span class="text-error me-4 self-center">* required field</span>
					<button class="btn btn-primary" type="submit">Save</button>
				</div>
			</form>
		</div>
	</div>
{/if}
