<script lang="ts">
	import { goto } from '$app/navigation';
	import { pocketbase } from '$lib/stores/pocketbase';
	import DeviceFormPort from '$lib/components/DeviceFormPort.svelte';
	import Alert from '$lib/components/Alert.svelte';
	import { faSave, faTrash, faTriangleExclamation } from '@fortawesome/free-solid-svg-icons';
	import type { Device, Port } from '$lib/types/device';
	import Fa from 'svelte-fa';

	export let device: Device;

	let portErrMsg = '';
	let portErrTimeout: number;
	let saveErrMsg = '';
	let saveErrTimeout: number;
	let deleteModal: HTMLDialogElement;

	async function save() {
		// create/update all ports
		let portIds: string[] = [];
		await Promise.all(
			device.expand.ports.map(async (port) => {
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
		device.ports = portIds;
		device.id ? updateDevice(device) : createDevice(device);
	}

	async function updateDevice(device: Device) {
		$pocketbase
			.collection('devices')
			.update(device.id, device)
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

	async function createDevice(device: Device) {
		$pocketbase
			.collection('devices')
			.create(device)
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

	function deleteDevice() {
		$pocketbase
			.collection('devices')
			.delete(device.id)
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

	function createEmptyPort() {
		if (device === undefined) return;
		device.expand.ports = [...device.expand.ports, { name: '', number: 1 } as Port];
	}
</script>

<form on:submit|preventDefault={save}>
	<div class="card w-full bg-base-300 shadow-xl">
		<div class="card-body">
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
						bind:value={device.name}
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
						bind:value={device.ip}
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
						bind:value={device.mac}
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
						bind:value={device.netmask}
						required
					/>
				</div>
				<span class="text-error me-4 self-center">* required field</span>
			</div>
		</div>
	</div>
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Ports</h2>
			<p class="my-2">UpSnap can also check if given ports are open. You can define them below.</p>
			<div class="form-control w-full">
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
					{#each device.expand.ports as _, index}
						<DeviceFormPort bind:device {index} bind:portErrMsg bind:portErrTimeout />
					{/each}
				</div>
				{#if portErrMsg !== ''}
					<Alert
						color="error"
						icon={faTriangleExclamation}
						message={portErrMsg}
						customClasses="mt-4 max-w-fit"
					/>
				{/if}
				<button
					class="btn btn-wide btn-neutral mt-4"
					on:click={() => createEmptyPort()}
					type="button">Add new port</button
				>
			</div>
		</div>
	</div>
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Link</h2>
			<p class="my-2">
				Makes your device name a clickable link, perfect for linking a dashboard for example.
			</p>
			<div class="form-control w-full">
				<input
					type="url"
					placeholder="https:// ..."
					class="input w-full max-w-xs"
					bind:value={device.link}
				/>
			</div>
		</div>
	</div>
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Wake</h2>
			<p class="my-2">
				You can power this device using a scheduled cron job. If you are not familiar with cron, you
				can read about it <a class="link" href="https://en.wikipedia.org/wiki/Cron" target="_blank"
					>here</a
				>
				or refer to the
				<a class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3" target="_blank"
					>package documentation</a
				>
				for more information.
			</p>
			<div class="form-control flex flex-row flex-wrap gap-4">
				<div>
					<label class="label" for="wake-cron">
						<span class="label-text">Wake cron</span>
					</label>
					<input
						id="wake-cron"
						type="text"
						placeholder="M H DoM M DoW"
						class="input w-80"
						bind:value={device.wake_cron}
						disabled={!device.wake_cron_enabled}
					/>
				</div>
				<div class="flex flex-col">
					<label class="label cursor-pointer" for="wake-cron-enable">
						<span class="label-text">Enable wake cron</span>
					</label>
					<input
						id="wake-cron-enable"
						type="checkbox"
						class="toggle toggle-success"
						bind:checked={device.wake_cron_enabled}
					/>
				</div>
			</div>
		</div>
	</div>
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Shutdown</h2>
			<p class="my-2">
				This <strong>shell command</strong> will run inside your container (if you use docker) or on
				your host (if you use the binary). To verify that works you can run the command inside the
				container or on your host shell first. Common commands are
				<span class="badge">net rpc</span> for windows, <span class="badge">sshpass</span> for linux
				or
				<span class="badge">curl</span> in general to make web requests.
			</p>
			<p class="my-2 font-bold">Examples:</p>
			<div class="alert text-sm bg-base-100 text-start w-fit">
				<div>
					<span class="font-bold">Shutdown remote windows machine:</span>
					<pre
						class="break-words whitespace-pre-wrap">net rpc shutdown -I 192.168.1.13 -U "user%password"</pre>
				</div>
			</div>
			<div class="alert text-sm my-2 bg-base-100 text-start w-fit">
				<div>
					<span class="font-bold">Shutdown remote linux machine:</span>
					<pre
						class="break-words whitespace-pre-wrap">sshpass -p your_password ssh -o "StrictHostKeyChecking=no" user@hostname "sudo poweroff"</pre>
				</div>
			</div>
			<div class="form-control w-full">
				<label class="label cursor-pointer" for="shutdown-cmd">
					<span class="label-text">Shutdown command</span>
				</label>
				<input
					id="shutdown-cmd"
					type="text"
					placeholder="$:"
					class="input w-full max-w-xs mb-2"
					bind:value={device.shutdown_cmd}
				/>
			</div>
			<p class="my-2">
				Just like setting a cron to wake the device, you can also schedule a cron job to shut down
				this device.
			</p>
			<div class="form-control flex flex-row flex-wrap gap-4">
				<div>
					<label class="label" for="shutdown-cron">
						<span class="label-text">Shutdown cron</span>
					</label>
					<input
						id="shutdown-cron"
						type="text"
						placeholder="M H DoM M DoW"
						class="input w-80"
						bind:value={device.shutdown_cron}
						disabled={!device.shutdown_cron_enabled}
					/>
				</div>
				<div class="flex flex-col">
					<label class="label cursor-pointer" for="shutdown-cron-enable">
						<span class="label-text">Enable shutdown cron</span>
					</label>
					<input
						id="shutdown-cron-enable"
						type="checkbox"
						class="toggle toggle-success"
						bind:checked={device.shutdown_cron_enabled}
					/>
				</div>
			</div>
		</div>
	</div>
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Password</h2>
			<p class="my-2">
				Some network cards have the option to set a password for magic packets, also called <span
					class="badge">SecureON</span
				>. Password can only be 0, 4 or 6 characters in length.
			</p>
			<div class="form-control w-full">
				<input
					type="password"
					placeholder="123456"
					class="input w-full max-w-xs"
					maxlength="6"
					bind:value={device.password}
				/>
			</div>
			{#if saveErrMsg !== ''}
				<Alert
					color="error"
					message={saveErrMsg}
					icon={faTriangleExclamation}
					customClasses="mt-4 max-w-fit"
				/>
			{/if}
		</div>
	</div>
	<div class="card-actions mt-4 justify-end gap-4">
		{#if device.id}
			<button class="btn btn-error" type="button" on:click={() => deleteModal.showModal()}
				><Fa icon={faTrash} />Delete</button
			>
			<dialog id="my_modal_1" class="modal" bind:this={deleteModal}>
				<form method="dialog" class="modal-box">
					<h3 class="font-bold text-lg">Confirm delete</h3>
					<p class="py-4">Are you sure you want to delete <strong>{device.name}</strong>?</p>
					<div class="modal-action">
						<button class="btn">Cancle</button>
						<button class="btn btn-error" on:click={() => deleteDevice()}>Delete</button>
					</div>
				</form>
			</dialog>
		{/if}
		<button class="btn btn-success" type="submit"><Fa icon={faSave} />Save</button>
	</div>
</form>
