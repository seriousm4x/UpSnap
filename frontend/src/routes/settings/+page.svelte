<script>
	import PocketBase from 'pocketbase';
	import { onMount } from 'svelte';

	let pb;
	let files;

	onMount(() => {
		pb = new PocketBase('http://127.0.0.1:8090');
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

			let createdPorts = {};
			for (let index = 0; index < data.length; index++) {
				const device = data[index];

				// create ports
				let thisDevicePorts = [];
				for (let index = 0; index < device.ports.length; index++) {
					// keep track of already created ports
					const port = device.ports[index];
					if (createdPorts[port.number] === undefined) {
						const record = await pb.collection('ports').create(
							{
								name: port.name,
								number: port.number
							},
							{ $autoCancel: false }
						);
						createdPorts[port.number] = record;
						thisDevicePorts.push(record.id);
					} else {
						thisDevicePorts.push(createdPorts[port.number].id);
					}
				}

				// create device
				await pb.collection('devices').create(
					{
						name: device.name,
						ip: device.ip,
						mac: device.mac,
						netmask: device.netmask,
						port: thisDevicePorts,
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
    <h1>Settings</h1>
	<div class="input-group flex-nowrap">
		<input
			class="form-control"
			placeholder="Username"
			aria-label="Username"
			aria-describedby="addon-wrapping"
			type="file"
			accept=".json"
			bind:files
			on:change={(e) => restoreV2Backup(e)}
		/>
	</div>
</div>
