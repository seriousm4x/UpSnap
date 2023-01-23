<script>
	import { onMount } from 'svelte';
	import { pocketbase } from '@stores/pocketbase';

	let pb;
	let files;
	let settings = {};
	let buttons = {
		saved: 'none',
		restore: 'none'
	}

	onMount(async () => {
		pocketbase.subscribe((conn) => {
			pb = conn;
		});

		const result = await pb.collection('settings').getList(1, 1);
		settings = result.items[0];
	});

	async function saveSettings() {
		buttons.saved = 'waiting';
		try {
			const result = await pb.collection('settings').update(
				settings.id,
				{
					interval: settings.interval,
					notifications: settings.notifications
				},
				{ $autoCancel: false }
			);
			setTimeout(() => {
				buttons.saved = 'none';
			}, 2000);
			buttons.saved = 'success';
		} catch (error) {
			setTimeout(() => {
				buttons.saved = 'none';
			}, 2000);
			buttons.saved = 'failed';
		}
	}

	async function restoreV2Backup() {
		buttons.restore = 'waiting'

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
				buttons.restore = 'none';
			}, 2000);
			buttons.restore = 'success';
		} catch (error) {
			setTimeout(() => {
				buttons.restore = 'none';
			}, 2000);
			buttons.restore = 'failed';
		}
	}
</script>

<div class="container">
	<section class="m-0 mt-4 p-4 shadow-sm">
		<h3 class="mb-3">Ping interval</h3>
		<div class="row">
			<div class="col-md-6">
				<p>Sets the interval in which the devices are pinged.</p>
				<div class="input-group mb-3">
					<span class="input-group-text">Cron</span>
					<input
						id="settings-interval"
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
	</section>
	<section class="m-0 mt-4 p-4 shadow-sm">
		<h3 class="mb-3">Notifications</h3>
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
	</section>
	<div class="m-0 mt-4 shadow-sm">
		<button
			type="button"
			class="btn btn-secondary w-100"
			class:btn-success={buttons.saved === 'success' ? true : false}
			class:btn-warning={buttons.saved === 'waiting' ? true : false}
			class:btn-danger={buttons.saved === 'failed' ? true : false}
			on:click={() => saveSettings()}
		>
			{#if buttons.saved === 'none'}
				Save
			{:else if buttons.saved === 'success'}
				Saved
			{:else if buttons.saved === 'waiting'}
				Waiting
			{:else if buttons.saved === 'failed'}
				Failed
			{/if}
		</button>
	</div>
	<section class="m-0 mt-5 p-4 shadow-sm">
		<h3 class="mb-3">Restore</h3>
		<p>If you had UpSnap v2.3.2 (or higher) running before, you can restore your devices here.</p>
		<div class="callout callout-danger fw-bold">This will wipe the existing database!</div>
		<input
			class="form-control mb-3"
			placeholder="Username"
			aria-label="Username"
			aria-describedby="addon-wrapping"
			type="file"
			accept=".json"
			bind:files
		/>
		<button
			type="button"
			class="btn btn-secondary"
			class:btn-success={buttons.restore === 'success' ? true : false}
			class:btn-warning={buttons.restore === 'waiting' ? true : false}
			class:btn-danger={buttons.restore === 'failed' ? true : false}
			disabled={files ? false : true}
			on:click={() => restoreV2Backup()}
		>
			{#if buttons.restore === 'none'}
				Restore
			{:else if buttons.restore === 'success'}
				Restored
			{:else if buttons.restore === 'waiting'}
				Waiting
			{:else if buttons.restore === 'failed'}
				Failed
			{/if}
		</button>
	</section>
</div>
