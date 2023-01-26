<script>
	import { onMount } from 'svelte';
	import { pocketbase } from '@stores/pocketbase';
	import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';

	let pb;
	let files;
	let settings = {};
	let buttons = {
		add_device: {
			state: 'none',
			error: ''
		},
		saved: {
			state: 'none',
			error: ''
		},
		restore: {
			state: 'none',
			error: ''
		}
	};
	let newDevice = {
		name: '',
		ip: '',
		mac: '',
		netmask: '255.255.255.0',
		link: '',
		wake_cron: '',
		shutdown_cron: '',
		shutdown_cmd: '',
		password: ''
	};

	let passwordShow = false;
	$: passwordType = passwordShow ? 'text' : 'password';

	onMount(async () => {
		pocketbase.subscribe((conn) => {
			pb = conn;
		});

		const result = await pb.collection('settings').getList(1, 1);
		settings = result.items[0];
	});

	async function addDevice() {
		buttons.add_device.state = 'waiting';
		try {
			await pb.collection('devices').create(newDevice, { $autoCancel: false });
			setTimeout(() => {
				buttons.add_device.state = 'none';
			}, 3000);
			buttons.add_device.state = 'success';
		} catch (error) {
			setTimeout(() => {
				buttons.add_device.error = '';
				buttons.add_device.state = 'none';
			}, 3000);
			buttons.add_device.error = error;
			buttons.add_device.state = 'failed';
		}
	}

	async function saveSettings() {
		buttons.saved.state = 'waiting';
		try {
			await pb.collection('settings').update(
				settings.id,
				{
					interval: settings.interval,
					notifications: settings.notifications
				},
				{ $autoCancel: false }
			);
			setTimeout(() => {
				buttons.saved.state = 'none';
			}, 3000);
			buttons.saved.state = 'success';
		} catch (error) {
			setTimeout(() => {
				buttons.saved.error = '';
				buttons.saved.state = 'none';
			}, 3000);
			buttons.saved.error = error;
			buttons.saved.state = 'failed';
		}
	}

	async function restoreV2Backup() {
		buttons.restore.state = 'waiting';

		try {
			// delete all devices in pocketbase
			let result = await pb.collection('devices').getFullList();
			for (let index = 0; index < result.length; index++) {
				pb.collection('devices').delete(result[index].id);
			}
			// delete all ports in pocketbase
			result = await pb.collection('ports').getFullList();
			for (let index = 0; index < result.length; index++) {
				await pb.collection('ports').delete(result[index].id);
			}

			let reader = new FileReader();
			reader.readAsText(files[0]);
			reader.onload = async (e) => {
				// parse uploaded file
				let data = JSON.parse(e.target.result);
				if (!Array.isArray(data)) {
					return;
				}

				// loop devices
				for (let index = 0; index < data.length; index++) {
					const device = data[index];

					// create ports
					let thisDevicePorts = [];
					for (let index = 0; index < device.ports.length; index++) {
						const port = device.ports[index];
						const record = await pb.collection('ports').create(
							{
								name: port.name,
								number: port.number
							},
							{ $autoCancel: false }
						);
						thisDevicePorts.push(record.id);
					}

					// create device
					await pb.collection('devices').create(
						{
							name: device.name,
							ip: device.ip,
							mac: device.mac,
							netmask: device.netmask,
							ports: thisDevicePorts,
							link: device.link,
							wake: device.wake,
							shutdown: device.shutdown
						},
						{ $autoCancel: false }
					);
				}
			};
			setTimeout(() => {
				buttons.restore.state = 'none';
			}, 3000);
			buttons.restore.state = 'success';
		} catch (error) {
			setTimeout(() => {
				buttons.restore.error = '';
				buttons.restore.state = 'none';
			}, 3000);
			buttons.restore.error = error;
			buttons.restore.state = 'failed';
		}
	}

	function onPasswordInput(event) {
		newDevice.password = event.target.value;
	}
</script>

<div class="container">
	<section class="m-0 mt-4 p-4 shadow-sm">
		<h3 class="mb-3">Add new device</h3>
		<div class="row">
			<div class="col-md-6">
				<form on:submit|preventDefault={addDevice}>
					<h5>Required:</h5>
					<div class="input-group mb-3">
						<span class="input-group-text">Name</span>
						<input
							class="form-control"
							placeholder="Office Pc"
							aria-label="Name"
							aria-describedby="addon-wrapping"
							type="text"
							required
							bind:value={newDevice.name}
						/>
					</div>
					<div class="input-group mb-3">
						<span class="input-group-text">IP</span>
						<input
							class="form-control"
							placeholder="192.168.1.34"
							aria-label="IP"
							aria-describedby="addon-wrapping"
							type="text"
							required
							bind:value={newDevice.ip}
						/>
					</div>
					<div class="input-group mb-3">
						<span class="input-group-text">MAC</span>
						<input
							class="form-control"
							placeholder="aa:bb:cc:dd:ee:ff"
							aria-label="MAC"
							aria-describedby="addon-wrapping"
							type="text"
							required
							bind:value={newDevice.mac}
						/>
					</div>
					<div class="input-group mb-3">
						<span class="input-group-text">Netmask</span>
						<input
							class="form-control"
							placeholder="Most likely 255.255.255.0 or 255.255.0.0"
							aria-label="Netmask"
							aria-describedby="addon-wrapping"
							type="text"
							required
							bind:value={newDevice.netmask}
						/>
					</div>
					<h5 class="mt-4">Optional:</h5>
					<div class="input-group mb-3">
						<span class="input-group-text">Link</span>
						<input
							class="form-control"
							placeholder="Clickable link on device card"
							aria-label="Link"
							aria-describedby="addon-wrapping"
							type="url"
							bind:value={newDevice.link}
						/>
					</div>
					<div class="input-group mb-3">
						<span class="input-group-text">Wake Cron<sup>(1)</sup></span>
						<input
							class="form-control"
							placeholder="Automatically wake device with cron"
							aria-label="Wake Cron"
							aria-describedby="addon-wrapping"
							type="text"
							bind:value={newDevice.wake_cron}
						/>
					</div>
					<div class="input-group mb-3">
						<span class="input-group-text">Shutdown Cron<sup>(1)</sup></span>
						<input
							class="form-control"
							placeholder="Automatically shutdown device with cron"
							aria-label="Shutdown Cron"
							aria-describedby="addon-wrapping"
							type="text"
							bind:value={newDevice.shutdown_cron}
						/>
					</div>
					<div class="input-group mb-3">
						<span class="input-group-text">Shutdown Cmd<sup>(2)</sup></span>
						<input
							class="form-control"
							placeholder="Command to shutdown device"
							aria-label="Shutdown Cmd"
							aria-describedby="addon-wrapping"
							type="text"
							bind:value={newDevice.shutdown_cmd}
						/>
					</div>
					<div class="input-group mb-3">
						<span class="input-group-text">Password<sup>(3)</sup></span>
						<input
							class="form-control"
							placeholder="BIOS password for wol"
							aria-label="IP"
							aria-describedby="addon-wrapping"
							type={passwordType}
							value={newDevice.password}
							on:input={onPasswordInput}
						/>
						<button
							class="btn btn-outline-secondary"
							type="button"
							on:click={() => (passwordShow = !passwordShow)}
						>
							<Fa icon={passwordShow ? faEyeSlash : faEye} />
						</button>
					</div>
					<button
						type="submit"
						class="btn btn-secondary"
						class:btn-success={buttons.add_device.state === 'success' ? true : false}
						class:btn-warning={buttons.add_device.state === 'waiting' ? true : false}
						class:btn-danger={buttons.add_device.state === 'failed' ? true : false}
					>
						{#if buttons.add_device.state === 'none'}
							Add device
						{:else if buttons.add_device.state === 'success'}
							Added successfully
						{:else if buttons.add_device.state === 'waiting'}
							Waiting
						{:else if buttons.add_device.state === 'failed'}
							Failed: {buttons.add_device.error}
						{/if}
					</button>
				</form>
			</div>
			<div class="col-md-6">
				<div class="callout callout-info mb-0">
					<h5>Optional:</h5>
					<p class="m-0">(1) Same cron syntax as for ping interval.</p>
					<p class="m-0">(2) Shell command to be executed. e.g.:</p>
					<ul>
						<li>Windows: "net rpc shutdown -I 192.168.1.13 -U test%test"</li>
						<li>
							Linux: "sshpass -p your_password ssh -o 'StrictHostKeyChecking=no' user@hostname 'sudo
							shutdown'"
						</li>
					</ul>
					<p class="m-0">
						(3) Some network cards have the option to set a password for magic packets.
					</p>
				</div>
			</div>
		</div>
	</section>
	<section class="m-0 mt-4 p-4 shadow-sm">
		<h3 class="mb-3">Ping interval</h3>
		<div class="row">
			<div class="col-md-6">
				<p>Sets the interval in which the devices are pinged.</p>
				<div class="input-group mb-3">
					<span class="input-group-text">Cron</span>
					<input
						class="form-control"
						placeholder="e.g. '@every 5s' or '@every 1m'"
						aria-label="Interval"
						aria-describedby="addon-wrapping"
						type="text"
						bind:value={settings.interval}
					/>
				</div>
			</div>
			<div class="col-md-6">
				<div class="callout callout-info m-0">
					<p class="m-0">
						Leave blank to use default value of <span class="fw-bold">"@every 3s"</span>.
					</p>
					<p class="m-0">
						Read more about valid cron syntax at <a
							target="_blank"
							rel="noreferrer"
							href="https://pkg.go.dev/github.com/robfig/cron">pkg.go.dev/github.com/robfig/cron</a
						>
					</p>
				</div>
			</div>
		</div>
		<h3 class="my-3">Notifications</h3>
		<p>Show toast notifications in the bottom right corner.</p>
		<div class="form-check form-switch">
			<label class="form-check-label" for="settings-notifications">Enable notifications</label>
			<input
				class="form-check-input"
				type="checkbox"
				id="settings-notifications"
				role="switch"
				bind:checked={settings.notifications}
			/>
		</div>
		<button
			type="button"
			class="btn btn-secondary mt-3"
			class:btn-success={buttons.saved.state === 'success' ? true : false}
			class:btn-warning={buttons.saved.state === 'waiting' ? true : false}
			class:btn-danger={buttons.saved.state === 'failed' ? true : false}
			on:click={() => saveSettings()}
		>
			{#if buttons.saved.state === 'none'}
				Save
			{:else if buttons.saved.state === 'success'}
				Saved
			{:else if buttons.saved.state === 'waiting'}
				Waiting
			{:else if buttons.saved.state === 'failed'}
				Failed: {buttons.add_device.error}
			{/if}
		</button>
	</section>
	<section class="m-0 my-4 p-4 shadow-sm">
		<h3 class="mb-3">Restore</h3>
		<p>If you had UpSnap v2.3.2 (or higher) running before, you can restore your devices here.</p>
		<div class="callout callout-danger fw-bold">This will wipe the existing database!</div>
		<input
			class="form-control"
			placeholder="Username"
			aria-label="Username"
			aria-describedby="addon-wrapping"
			type="file"
			accept=".json"
			bind:files
		/>
		<button
			type="button"
			class="btn btn-secondary mt-3"
			class:btn-success={buttons.restore.state === 'success' ? true : false}
			class:btn-warning={buttons.restore.state === 'waiting' ? true : false}
			class:btn-danger={buttons.restore.state === 'failed' ? true : false}
			disabled={files ? false : true}
			on:click={() => restoreV2Backup()}
		>
			{#if buttons.restore.state === 'none'}
				Restore
			{:else if buttons.restore.state === 'success'}
				Restored
			{:else if buttons.restore.state === 'waiting'}
				Waiting
			{:else if buttons.restore.state === 'failed'}
				Failed: {buttons.restore.error}
			{/if}
		</button>
	</section>
</div>
