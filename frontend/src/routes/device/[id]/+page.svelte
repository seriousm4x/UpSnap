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
	let deviceClone: Device | undefined;

	let portErrMsg = '';
	let portErrTimeout: number;

	let saveErrMsg = '';
	let saveErrTimeout: number;

	devices.subscribe((data) => {
		if (waitingForDevices === false) return;
		if (data.length === 0) return;

		device = data.find((dev) => dev.id === $page.params.id);
		if (device !== undefined) {
			deviceClone = device;
			if (deviceClone.expand.ports === undefined) {
				deviceClone.expand.ports = [];
			}
		}
		waitingForDevices = false;
	});

	async function saveDevice() {
		if (deviceClone === undefined) return;

		// create/update all ports
		let portIds: string[] = [];
		await Promise.all(
			deviceClone.expand.ports.map(async (port) => {
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

		deviceClone.ports = portIds;
		$pocketbase
			.collection('devices')
			.update(deviceClone.id, deviceClone)
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
		if (deviceClone === undefined) return;
		deviceClone.expand.ports = [...deviceClone.expand.ports, { name: '', number: 1 } as Port];
	}
</script>

{#if waitingForDevices}
	<div class="container max-w-lg mx-auto text-center">
		<span class="loading loading-dots loading-lg" />
	</div>
{:else if device === undefined || deviceClone === undefined}
	<div class="container max-w-lg mx-auto">
		<Alert
			color="error"
			message="No device with id: {$page.params.id}"
			icon={faTriangleExclamation}
		/>
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
							bind:value={deviceClone.name}
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
							bind:value={deviceClone.ip}
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
							bind:value={deviceClone.mac}
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
							bind:value={deviceClone.netmask}
							required
						/>
					</div>
				</div>
				<h2 class="card-title mt-8">Ports</h2>
				<p class="my-2">
					UpSnap can also check if given ports are open. You can define them below.
				</p>
				<div class="form-control w-full">
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
						{#each deviceClone.expand.ports as _, index}
							<DeviceFormPort bind:deviceClone {index} bind:portErrMsg bind:portErrTimeout />
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
				<h2 class="card-title mt-8">Link</h2>
				<p class="my-2">
					Makes your device name a clickable link, perfect for linking a dashboard for example.
				</p>
				<div class="form-control w-full">
					<input
						type="url"
						placeholder="https:// ..."
						class="input w-full max-w-xs"
						bind:value={deviceClone.link}
					/>
				</div>
				<h2 class="card-title mt-8">Wake</h2>
				<p class="my-2">
					You can power this device using a scheduled cron job. If you are not familiar with cron,
					you can read about it <a
						class="link"
						href="https://en.wikipedia.org/wiki/Cron"
						target="_blank">here</a
					>
					or refer to the
					<a class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3" target="_blank"
						>package documentation</a
					>
					for more information. <strong>Tip:</strong> Once you've entered your syntax into the input
					field and click out of it, your syntax will be automatically validated to check for correctness.
				</p>
				<div class="form-control flex flex-row flex-wrap gap-4">
					<div>
						<label class="label" for="wake-cron">
							<span class="label-text">Wake cron</span>
						</label>
						<input
							id="wake-cron"
							type="text"
							placeholder="0 9 * * 1-5"
							class="input w-80"
							bind:value={deviceClone.wake_cron}
							disabled={!deviceClone.wake_cron_enabled}
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
							bind:checked={deviceClone.wake_cron_enabled}
						/>
					</div>
				</div>
				<h2 class="card-title mt-8">Shutdown</h2>
				<p class="my-2">
					This <strong>shell command</strong> will run inside your container (if you use docker) or
					on your host (if you use the binary). To verify that works you can run the command inside
					the container or on your host shell first. Common commands are
					<span class="badge">net rpc</span> for windows, <span class="badge">sshpass</span> for
					linux or
					<span class="badge">curl</span> in general to make web requests.
				</p>
				<p class="my-2 font-bold">Example:</p>
				<div class="alert text-sm bg-base-100 text-start">
					<div>
						<span class="font-bold">Shutdown remote windows machine:</span>
						<pre
							class="break-words whitespace-pre-wrap">net rpc shutdown -I 192.168.1.13 -U "user%password"</pre>
					</div>
				</div>
				<div class="alert text-sm my-2 bg-base-100 text-start">
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
						bind:value={deviceClone.shutdown_cmd}
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
							placeholder="0 22 * * 1-5"
							class="input w-80"
							bind:value={deviceClone.shutdown_cron}
							disabled={!deviceClone.shutdown_cron_enabled}
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
							bind:checked={deviceClone.shutdown_cron_enabled}
						/>
					</div>
				</div>
				<h2 class="card-title mt-8">Password</h2>
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
						bind:value={deviceClone.password}
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
				<div class="card-actions justify-end mt-4">
					<span class="text-error me-4 self-center">* required field</span>
					<button class="btn btn-primary" type="submit">Save</button>
				</div>
			</form>
		</div>
	</div>
{/if}
