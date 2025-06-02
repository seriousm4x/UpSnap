<script lang="ts">
	import { goto } from '$app/navigation';
	import DeviceFormPort from '$lib/components/DeviceFormPort.svelte';
	import { nextCronDate, parseCron as validateCron } from '$lib/helpers/cron';
	import { m } from '$lib/paraglide/messages';
	import { permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Group, Port } from '$lib/types/device';
	import { faSave, faTrash, faX } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	export let device: Device;

	let deleteModal: HTMLDialogElement;
	let deviceGroups = [] as Group[];
	let newGroup = '';

	onMount(async () => {
		await getGroups();
	});

	async function save() {
		// validate crons
		if (device.wake_cron_enabled && !(await validateCron(device.wake_cron))) {
			toast.error(m.settings_invalid_cron());
			throw new Error('wake_cron not valid');
		}
		if (device.shutdown_cron_enabled && !(await validateCron(device.shutdown_cron))) {
			toast.error(m.settings_invalid_cron());
			throw new Error('shutdown_cron not valid');
		}

		// create/update all ports
		let portIds: string[] = [];
		await Promise.all(
			device.expand.ports.map(async (port: Port) => {
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
		device.ip = device.ip.replaceAll(' ', '');
		device.mac = device.mac.replaceAll(' ', '');
		device.netmask = device.netmask.replaceAll(' ', '');
		device.wake_cmd = device.wake_cmd.trim();
		device.shutdown_cmd = device.shutdown_cmd.trim();
		device.wake_cron = device.wake_cron.trim();
		device.shutdown_cron = device.shutdown_cron.trim();
		device.ping_cmd = device.ping_cmd.trim();
		device.ports = portIds;
		if (device.id) {
			updateDevice(device);
		} else {
			createDevice(device);
		}
	}

	async function updateDevice(device: Device) {
		$pocketbase
			.collection('devices')
			.update(device.id, device)
			.then(() => {
				toast.success(m.toasts_device_updated({ device: device.name }));
				goto('/');
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	async function createDevice(device: Device) {
		device.created_by = $pocketbase.authStore.isSuperuser
			? ''
			: ($pocketbase.authStore.record?.id ?? '');
		$pocketbase
			.collection('devices')
			.create(device)
			.then(() => {
				toast.success(m.toasts_device_created({ device: device.name }));
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
				toast.success(m.toasts_device_deleted({ device: device.name }));
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
				toast.success(m.toasts_group_created({ group: newGroup }));
				newGroup = '';
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
				toast.success(m.toasts_group_deleted({ group: group.name }));
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
	<div class="card bg-base-200 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_general()}</h2>
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_general_name()} <span class="text-error">*</span></span>
						<input
							type="text"
							placeholder={m.device_general_name()}
							class="input"
							bind:value={device.name}
							required
						/>
					</label>
				</fieldset>
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_general_ip()} <span class="text-error">*</span></span>
						<input
							type="text"
							placeholder={m.device_general_ip()}
							class="input"
							bind:value={device.ip}
							required
						/>
					</label>
				</fieldset>
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_general_mac()} <span class="text-error">*</span></span>
						<input
							type="text"
							placeholder={m.device_general_mac()}
							class="input"
							bind:value={device.mac}
							required
						/>
					</label>
				</fieldset>
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_general_netmask()} <span class="text-error">*</span></span>
						<input
							type="text"
							placeholder={m.device_general_netmask()}
							class="input"
							bind:value={device.netmask}
							required
						/>
					</label>
				</fieldset>
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_general_description()}</span>
						<input
							type="text"
							placeholder={m.device_general_description()}
							class="input"
							bind:value={device.description}
						/>
					</label>
				</fieldset>
				<span class="badge text-error self-center">* {m.device_general_required_field()}</span>
			</div>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_ports()}</h2>
			<p>{m.device_ports_desc()}</p>
			<div class="w-full">
				<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
					<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
					{#each device.expand.ports as _, index}
						<DeviceFormPort bind:device {index} />
					{/each}
				</div>
				<button
					class="btn max-w-xs {device.expand.ports.length > 0 ? 'mt-4' : ''}"
					on:click={() => createEmptyPort()}
					type="button">{m.device_ports_add_new()}</button
				>
			</div>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_link()}</h2>
			<p>
				{m.device_link_desc()}
			</p>
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_link()}</span>
						<input
							type="url"
							placeholder={m.device_link()}
							class="input"
							bind:value={device.link}
						/>
					</label>
				</fieldset>
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_link_open()}</span>
						<select class="select" bind:value={device.link_open}>
							<option value="">{m.device_link_open_no()}</option>
							<option value="same_tab">{m.device_link_open_same_tab()}</option>
							<option value="new_tab">{m.device_link_open_new_tab()}</option>
						</select>
					</label>
				</fieldset>
			</div>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_ping()}</h2>
			<p>
				<!-- eslint-disable svelte/no-at-html-tags -->
				{@html m.device_ping_desc()}
			</p>
			<fieldset class="fieldset">
				<label class="input mt-2">
					<span>$:</span>
					<input
						type="text"
						placeholder={m.device_ping_cmd()}
						class="grow"
						bind:value={device.ping_cmd}
					/>
				</label>
			</fieldset>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_wake()}</h2>
			<p>
				{m.device_wake_desc()}
				<!-- eslint-disable svelte/no-at-html-tags -->
				{@html m.settings_ping_interval_desc2()}
			</p>
			<div class="flex flex-row flex-wrap gap-4">
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_wake_timeout()}</span>
						<input
							type="number"
							min="0"
							placeholder={m.device_wake_timeout()}
							class="input"
							bind:value={device.wake_timeout}
						/>
					</label>
				</fieldset>
				<fieldset class="fieldset">
					<label class="input mt-2">
						<span>$:</span>
						<input
							type="text"
							placeholder={m.device_wake_cmd()}
							class="grow"
							bind:value={device.wake_cmd}
						/>
					</label>
				</fieldset>
			</div>
			<div class="flex flex-row flex-wrap gap-4">
				<fieldset
					class="fieldset bg-base-100 border-base-300 rounded-box w-fit min-w-64 border p-4"
				>
					<legend class="fieldset-legend">{m.device_wake_cron_enable()}</legend>
					<label class="fieldset-label">
						<input
							type="checkbox"
							bind:checked={device.wake_cron_enabled}
							class="toggle toggle-success"
						/>
						{m.device_wake_cron_enable()}
					</label>
					<label class="floating-label mt-2">
						<span
							>{m.device_wake_cron()}
							{#if device.wake_cron_enabled}
								<span class="text-error">*</span>
							{/if}
						</span>
						<input
							type="text"
							placeholder={m.device_wake_cron()}
							class="input"
							bind:value={device.wake_cron}
							disabled={!device.wake_cron_enabled}
							required={device.wake_cron_enabled}
						/>
						{#if device.wake_cron_enabled}
							<p class="fieldset-label">
								{#await validateCron(device.wake_cron)}
									<span class="loading loading-spinner loading-xs"></span>
								{:then valid}
									{valid ? '✅ ' + nextCronDate(device.wake_cron) : m.settings_invalid_cron()}
								{/await}
							</p>
						{/if}
					</label>
					{#if device.wake_cron_enabled}
						<pre class="bg-base-100 text-base-content/80 w-fit rounded p-2"><code
								>* * * * * *
| | | | | |
| | | | | day of the week (0–6)
| | | | month (1–12)
| | | day of the month (1–31)
| | hour (0–23)
| minute (0–59)
second (0–59, optional)
</code></pre>
					{/if}
				</fieldset>
				<fieldset class="fieldset bg-base-100 border-base-300 rounded-box w-64 border p-4">
					<legend class="fieldset-legend">{m.device_require_confirmation()}</legend>
					<label class="fieldset-label">
						<input
							type="checkbox"
							bind:checked={device.wake_confirm}
							class="toggle toggle-success"
						/>
						{m.device_require_confirmation()}
					</label>
				</fieldset>
			</div>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_shutdown()}</h2>
			<p>
				<!-- eslint-disable svelte/no-at-html-tags -->
				{@html m.device_shutdown_desc()}
			</p>
			<div class="flex flex-row flex-wrap gap-4">
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.device_shutdown_timeout()}</span>
						<input
							type="number"
							min="0"
							placeholder={m.device_shutdown_timeout()}
							class="input"
							bind:value={device.shutdown_timeout}
						/>
					</label>
				</fieldset>
				<fieldset class="fieldset">
					<label class="input mt-2">
						<span>$:</span>
						<input
							type="text"
							placeholder={m.device_shutdown_cmd()}
							class="grow"
							bind:value={device.shutdown_cmd}
						/>
					</label>
				</fieldset>
			</div>
			<p class="font-bold">{m.device_shutdown_examples()}</p>
			<div class="mockup-code max-w-fit min-w-0 text-sm">
				<pre data-prefix="#" class="italic"><code>{m.device_shutdown_examples_windows()}</code
					></pre>
				<pre data-prefix="$" class="text-warning"><code
						>net rpc shutdown -I 192.168.1.13 -U "user%password"</code
					></pre>
			</div>
			<div class="mockup-code max-w-fit min-w-0 text-sm">
				<pre data-prefix="#" class="italic"><code>{m.device_shutdown_examples_linux()}</code></pre>
				<pre data-prefix="$" class="text-warning"><code
						>sshpass -p password ssh -o "StrictHostKeyChecking=no" user@192.168.1.13 "sudo poweroff"</code
					></pre>
			</div>
			<p>
				{m.device_shutdown_cron_desc()}
			</p>
			<div class="flex flex-row flex-wrap gap-4">
				<fieldset
					class="fieldset bg-base-100 border-base-300 rounded-box w-fit min-w-64 border p-4"
				>
					<legend class="fieldset-legend">{m.device_shutdown_cron_enable()}</legend>
					<label class="fieldset-label">
						<input
							type="checkbox"
							bind:checked={device.shutdown_cron_enabled}
							class="toggle toggle-success"
						/>
						{m.device_shutdown_cron_enable()}
					</label>
					<label class="floating-label mt-2">
						<span
							>{m.device_shutdown_cron()}
							{#if device.shutdown_cron_enabled}
								<span class="text-error">*</span>
							{/if}
						</span>
						<input
							type="text"
							placeholder={m.device_shutdown_cron()}
							class="input validator"
							bind:value={device.shutdown_cron}
							disabled={!device.shutdown_cron_enabled}
							required={device.shutdown_cron_enabled}
						/>
						{#if device.shutdown_cron_enabled}
							<p class="fieldset-label">
								{#await validateCron(device.shutdown_cron)}
									<span class="loading loading-spinner loading-xs"></span>
								{:then valid}
									{valid ? '✅ ' + nextCronDate(device.shutdown_cron) : m.settings_invalid_cron()}
								{/await}
							</p>
						{/if}
					</label>
					{#if device.shutdown_cron_enabled}
						<pre class="bg-base-100 text-base-content/80 w-fit rounded p-2"><code
								>* * * * * *
| | | | | |
| | | | | day of the week (0–6)
| | | | month (1–12)
| | | day of the month (1–31)
| | hour (0–23)
| minute (0–59)
second (0–59, optional)
</code></pre>
					{/if}
				</fieldset>
				<fieldset class="fieldset bg-base-100 border-base-300 rounded-box w-64 border p-4">
					<legend class="fieldset-legend">{m.device_require_confirmation()}</legend>
					<label class="fieldset-label">
						<input
							type="checkbox"
							bind:checked={device.shutdown_confirm}
							class="toggle toggle-success"
						/>
						{m.device_require_confirmation()}
					</label>
				</fieldset>
			</div>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">Sleep-On-LAN</h2>
			<p>
				<!-- eslint-disable svelte/no-at-html-tags -->
				{@html m.device_sol_desc1()}
			</p>
			<p>
				{m.device_sol_desc2()}
			</p>
			<p class="font-bold">
				<!-- eslint-disable svelte/no-at-html-tags -->
				{@html m.device_sol_desc3()}
			</p>
			<div class="flex flex-row flex-wrap gap-4">
				<fieldset class="fieldset bg-base-100 border-base-300 rounded-box w-64 border p-4">
					<legend class="fieldset-legend">{m.device_sol_enable()}</legend>
					<label class="fieldset-label">
						<input
							type="checkbox"
							bind:checked={device.sol_enabled}
							class="toggle toggle-success"
						/>
						{m.device_sol_enable()}
					</label>
					<label class="floating-label mt-2">
						<span
							>{m.device_sol_port()}
							{#if device.sol_enabled}
								<span class="text-error">*</span>
							{/if}
						</span>
						<input
							class="input"
							type="number"
							min="1"
							max="65535"
							placeholder={m.device_sol_port()}
							bind:value={device.sol_port}
							disabled={!device.sol_enabled}
							required={device.sol_enabled}
						/>
					</label>
				</fieldset>
				{#if device.sol_enabled}
					<fieldset class="fieldset bg-base-100 border-base-300 rounded-box w-64 border p-4">
						<legend class="fieldset-legend">{m.device_sol_authorization()}</legend>
						<label class="fieldset-label">
							<input type="checkbox" bind:checked={device.sol_auth} class="toggle toggle-success" />
							{m.device_sol_authorization()}
						</label>
						<label class="floating-label mt-2">
							<span
								>{m.device_sol_user()}
								{#if device.sol_auth}
									<span class="text-error">*</span>
								{/if}
							</span>
							<input
								type="text"
								placeholder={m.device_sol_user()}
								class="input"
								bind:value={device.sol_user}
								disabled={!device.sol_auth}
								required={device.sol_auth}
							/>
						</label>
						<label class="floating-label mt-2">
							<span
								>{m.device_sol_password()}
								{#if device.sol_auth}
									<span class="text-error">*</span>
								{/if}
							</span>
							<input
								type="password"
								placeholder={m.device_sol_password()}
								class="input"
								bind:value={device.sol_password}
								disabled={!device.sol_auth}
								required={device.sol_auth}
							/>
						</label>
					</fieldset>
				{/if}
			</div>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_password()}</h2>
			<p>
				<!-- eslint-disable svelte/no-at-html-tags -->
				{@html m.device_password_desc()}
			</p>
			<fieldset class="fieldset">
				<label class="floating-label mt-2">
					<span>{m.device_password()}</span>
					<input
						type="text"
						placeholder={m.device_password()}
						class="input"
						maxlength="6"
						bind:value={device.password}
					/>
				</label>
			</fieldset>
		</div>
	</div>
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.device_groups()}</h2>
			<p>
				{m.device_groups_desc()}
			</p>
			<div class="flex flex-row flex-wrap gap-2">
				{#each deviceGroups as group}
					<div class="join">
						<div class=" tooltip" data-tip="Delete">
							<button
								class="btn btn-error join-item"
								type="button"
								on:click={() => deleteGroup(group)}><Fa icon={faX} /></button
							>
						</div>
						<div
							class="btn join-item bg-base-100 hover:bg-base-200"
							on:click={() => toggleGroup(group.id)}
							role="none"
						>
							<input
								type="checkbox"
								class="checkbox"
								checked={device.groups.indexOf(group.id) !== -1}
							/>
							{group.name}
						</div>
					</div>
				{/each}
			</div>
			<div class="join">
				<input
					class="input join-item"
					placeholder={m.device_groups_placeholder()}
					type="text"
					on:keydown={(e) => {
						if (e.key === 'Enter') {
							e.preventDefault();
							addGroup();
						}
					}}
					bind:value={newGroup}
				/>
				<button class="btn btn-neutral join-item" type="button" on:click={() => addGroup()}
					>{m.buttons_add()}</button
				>
			</div>
		</div>
	</div>
	<div class="card-actions mt-6 justify-end gap-4">
		{#if $pocketbase.authStore.isSuperuser || $permission.delete?.includes(device.id)}
			<button class="btn btn-error" type="button" on:click={() => deleteModal.showModal()}
				><Fa icon={faTrash} />{m.buttons_delete()}</button
			>
		{/if}
		<button class="btn btn-success" type="submit"><Fa icon={faSave} />{m.buttons_save()}</button>
	</div>
</form>
<dialog class="modal" bind:this={deleteModal}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">{m.users_confirm_delete_title()}</h3>
		<p class="py-4">{m.users_confirm_delete_desc({ username: device.name })}</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">{m.buttons_cancel()}</button>
				<button class="btn btn-error" on:click={() => deleteDevice()}>{m.buttons_delete()}</button>
			</form>
		</div>
	</div>
</dialog>
