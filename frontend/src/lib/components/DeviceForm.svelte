<script lang="ts">
	import { goto } from '$app/navigation';
	import { pocketbase, permission } from '$lib/stores/pocketbase';
	import DeviceFormPort from '$lib/components/DeviceFormPort.svelte';
	import toast from 'svelte-french-toast';
	import Fa from 'svelte-fa';
	import { faSave, faTrash, faX } from '@fortawesome/free-solid-svg-icons';
	import type { Device, Port, Group } from '$lib/types/device';
	import { onMount } from 'svelte';

	export let device: Device;

	let deleteModal: HTMLDialogElement;
	let deviceGroups = [] as Group[];
	let newGroup = '';

	onMount(async () => {
		await getGroups();
	});

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
				toast.success(`Updated ${device.name}`);
				goto('/');
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	async function createDevice(device: Device) {
		device.created_by = $pocketbase.authStore.isAdmin ? '' : $pocketbase.authStore.model?.id ?? '';
		$pocketbase
			.collection('devices')
			.create(device)
			.then(() => {
				toast.success(`Created ${device.name}`);
				goto('/');
			})
			.catch((err) => {
				toast.error(err.message);
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
				toast.error(err.message);

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
				toast.error(err.message);

				newport = undefined;
			});
		return newport;
	}

	function deleteDevice() {
		$pocketbase
			.collection('devices')
			.delete(device.id)
			.then(() => {
				toast.success(`Deleted ${device.name}`);
				goto('/');
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function createEmptyPort() {
		if (device === undefined) return;
		device.expand.ports = [...device.expand.ports, { name: '', number: 1 } as Port];
	}

	async function getGroups() {
		$pocketbase
			.collection('device_groups')
			.getFullList()
			.then((res) => {
				deviceGroups = res as Group[];
			});
	}

	function addGroup() {
		if (!newGroup) return;
		$pocketbase
			.collection('device_groups')
			.create({
				name: newGroup
			})
			.then((res) => {
				deviceGroups = [...deviceGroups, res as Group];
				toast.success(`Created group ${newGroup}`);
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function deleteGroup(group: Group) {
		const i = device.groups.indexOf(group.id);
		if (i !== -1) {
			device.groups.splice(i, 1);
			device.groups = device.groups;
		}
		$pocketbase
			.collection('device_groups')
			.delete(group.id)
			.then(async () => {
				await getGroups();
				toast.success(`Deleted group ${group.name}`);
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function toggleGroup(id: string) {
		const i = device.groups.indexOf(id);
		if (i !== -1) {
			device.groups.splice(i, 1);
			device.groups = device.groups;
		} else {
			device.groups = [...device.groups, id];
		}
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
				<span class="badge text-error self-center">* required field</span>
			</div>
		</div>
	</div>
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Ports</h2>
			<p class="my-2">UpSnap can also check if given ports are open. You can define them below.</p>
			<div class="form-control w-full">
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
					<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
					{#each device.expand.ports as _, index}
						<DeviceFormPort bind:device {index} />
					{/each}
				</div>
				<button
					class="btn btn-wide btn-primary mt-4"
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
						<span class="label-text"
							>Wake cron
							{#if device.wake_cron_enabled}
								<span class="text-error">*</span>
							{/if}
						</span>
					</label>
					<input
						id="wake-cron"
						type="text"
						placeholder="M H DoM M DoW"
						class="input w-80"
						bind:value={device.wake_cron}
						disabled={!device.wake_cron_enabled}
						required={device.wake_cron_enabled}
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
			<h2 class="card-title">Sleep-On-LAN</h2>
			<p class="mt-2">
				You can put computers to sleep using the <a
					class="link"
					href="https://github.com/SR-G/sleep-on-lan"
					target="_blank">Sleep-On-LAN</a
				>
				tool. Sleep-On-LAN (SOL) is an external tool/daemon that operates on the PCs you want to put
				to sleep, providing a REST endpoint. For instructions on setting up Sleep-On-LAN, please refer
				to the
				<a href="https://github.com/SR-G/sleep-on-lan#usage" class="link" target="_blank">Usage</a> section.
			</p>
			<p>
				SOL is confiugred to send requests over HTTP instead of UDP to enable authorization and make
				requests more reliable.
			</p>
			<p class="font-bold">
				Therefore, please ensure that you include <span class="badge">HTTP:&lt;YOURPORT&gt;</span>
				in the <span class="badge">Listeners</span> section of the
				<a href="https://github.com/SR-G/sleep-on-lan#configuration" class="link" target="_blank"
					>SOL configuration</a
				>.
			</p>
			<div class="flex flex-row flex-wrap gap-4 items-end mt-4">
				<div>
					<div class="form-control flex flex-row flex-wrap gap-4">
						<div class="flex flex-row gap-2 items-center">
							<input
								id="sol-enable"
								type="checkbox"
								class="toggle toggle-success"
								bind:checked={device.sol_enabled}
							/>
							<label class="label cursor-pointer" for="sol-enable">
								<span class="label-text">Enable Sleep-On-LAN</span>
							</label>
						</div>
					</div>
					<div class="form-control flex flex-col">
						<label class="label" for="sol-port">
							<span class="label-text"
								>SOL Port
								{#if device.sol_enabled}
									<span class="text-error">*</span>
								{/if}
							</span>
						</label>
						<input
							id="sol-port"
							type="number"
							min="1"
							max="65535"
							placeholder="Default: 8009"
							class="input w-80"
							bind:value={device.sol_port}
							disabled={!device.sol_enabled}
							required={device.sol_enabled}
						/>
					</div>
				</div>
				{#if device.sol_enabled}
					<div>
						<div class="form-control flex flex-row flex-wrap gap-4">
							<div class="flex flex-row gap-2 items-center">
								<input
									id="sol-auth"
									type="checkbox"
									class="toggle toggle-success"
									bind:checked={device.sol_auth}
								/>
								<label class="label cursor-pointer" for="sol-auth">
									<span class="label-text">Authorization</span>
								</label>
							</div>
						</div>
						<div class="form-control flex flex-col">
							<label class="label" for="sol-user">
								<span class="label-text"
									>SOL User
									{#if device.sol_auth}
										<span class="text-error">*</span>
									{/if}
								</span>
							</label>
							<input
								id="sol-user"
								type="text"
								placeholder="username"
								class="input w-80"
								bind:value={device.sol_user}
								disabled={!device.sol_auth}
								required={device.sol_auth}
							/>
						</div>
					</div>
					<div class="form-control flex flex-col">
						<label class="label" for="sol-password">
							<span class="label-text"
								>SOL Password
								{#if device.sol_auth}
									<span class="text-error">*</span>
								{/if}
							</span>
						</label>
						<input
							id="sol-password"
							type="password"
							placeholder="password"
							class="input w-80"
							bind:value={device.sol_password}
							disabled={!device.sol_auth}
							required={device.sol_auth}
						/>
					</div>
				{/if}
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
			<div class="mockup-code text-sm max-w-fit">
				<pre data-prefix="#"><code>Shutdown remote windows machine:</code></pre>
				<pre data-prefix="$" class="text-warning"><code
						>net rpc shutdown -I 192.168.1.13 -U "user%password"</code
					></pre>
			</div>
			<div class="mockup-code text-sm max-w-fit">
				<pre data-prefix="#"><code>Shutdown remote linux machine:</code></pre>
				<pre data-prefix="$" class="text-warning"><code
						>sshpass -p password ssh -o "StrictHostKeyChecking=no" user@192.168.1.13 "sudo poweroff"</code
					></pre>
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
						<span class="label-text"
							>Shutdown cron
							{#if device.shutdown_cron_enabled}
								<span class="text-error">*</span>
							{/if}
						</span>
					</label>
					<input
						id="shutdown-cron"
						type="text"
						placeholder="M H DoM M DoW"
						class="input w-80"
						bind:value={device.shutdown_cron}
						disabled={!device.shutdown_cron_enabled}
						required={device.shutdown_cron_enabled}
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
		</div>
	</div>
	<div class="card w-full bg-base-300 shadow-xl mt-6">
		<div class="card-body">
			<h2 class="card-title">Groups</h2>
			<p class="my-2">
				You can add devices to a group to have them sorted by group on the dashboard.
			</p>
			<div class="flex flex-row flex-wrap gap-2">
				{#each deviceGroups as group}
					<div class="join">
						<div class=" tooltip" data-tip="Delete">
							<button
								class="join-item btn btn-error"
								type="button"
								on:click={() => deleteGroup(group)}><Fa icon={faX} /></button
							>
						</div>
						<div
							class="btn bg-base-100 hover:bg-base-200 join-item"
							on:click={() => toggleGroup(group.id)}
							role="none"
						>
							<input
								type="checkbox"
								class="checkbox checked:checkbox-primary"
								checked={device.groups.indexOf(group.id) !== -1}
							/>
							{group.name}
						</div>
					</div>
				{/each}
			</div>
			<div class="join">
				<input
					class="input input-bordered join-item"
					placeholder="e.g. 'Basement' or 'Office'"
					type="text"
					bind:value={newGroup}
				/>
				<button class="btn btn-primary join-item" type="button" on:click={() => addGroup()}
					>Create</button
				>
			</div>
		</div>
	</div>
	<div class="card-actions mt-6 justify-end gap-4">
		{#if $pocketbase.authStore.isAdmin || $permission.delete?.includes(device.id)}
			<button class="btn btn-error" type="button" on:click={() => deleteModal.showModal()}
				><Fa icon={faTrash} />Delete</button
			>
			<dialog class="modal" bind:this={deleteModal}>
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
