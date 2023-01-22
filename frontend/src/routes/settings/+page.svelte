<script>
	import { onMount } from 'svelte';
	import { pocketbase } from '@stores/pocketbase';

	let pb;
	let files;

	onMount(() => {
		pocketbase.subscribe(conn => {
			pb = conn;
		});
	});

	async function restoreV2Backup(e) {
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

		let file = e.target.files[0];
		let reader = new FileReader();
		reader.readAsText(file);
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
	}
</script>

<div class="container">
	<h1 class="h1">Settings</h1>
	<section>
		<h3>Restore</h3>
		<div class="callout callout-info">
			If you were running UpSnap v2.3.2 (or higher), you can restore your devices here. This will
			wipe the existing database.
		</div>
		<input
			class="form-control w-50"
			placeholder="Username"
			aria-label="Username"
			aria-describedby="addon-wrapping"
			type="file"
			accept=".json"
			bind:files
			on:change={(e) => restoreV2Backup(e)}
		/>
	</section>
</div>
