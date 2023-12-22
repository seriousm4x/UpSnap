<script lang="ts">
	import LL from '$lib/i18n/i18n-svelte';
	import { pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Port } from '$lib/types/device';
	import { faTrash } from '@fortawesome/free-solid-svg-icons';
	import type { RecordModel } from 'pocketbase';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';
	import { scale } from 'svelte/transition';

	export let device: Device;
	export let index: number;

	function deletePort(port: Port | RecordModel) {
		if (port.id !== undefined) {
			$pocketbase
				.collection('ports')
				.delete(port.id)
				.then(() => {
					device.ports = device.ports.filter((id) => id !== device.ports[index]);
				})
				.catch((err) => {
					toast.error(err.message);
				});
		}
		device.expand.ports = device.expand.ports.filter((p) => p !== device.expand.ports[index]);
	}
</script>

{#if device.expand.ports[index]}
	<div class="card bg-base-100 shadow-sm" transition:scale={{ delay: 0, duration: 200 }}>
		<div class="card-body p-3">
			<div class="flex flex-row gap-2">
				<div>
					<label class="label p-0" for="port-name-{index}">
						<span class="label-text">{$LL.device.ports_name()}</span>
					</label>
					<input
						id="port-name-{index}"
						type="text"
						placeholder="ssh"
						class="input input-sm input-bordered w-full"
						required
						bind:value={device.expand.ports[index].name}
					/>
				</div>
				<div>
					<label class="label p-0" for="port-number-{index}">
						<span class="label-text">{$LL.device.ports_number()}</span>
					</label>
					<input
						id="port-number-{index}"
						type="number"
						placeholder="22"
						min="1"
						max="65535"
						class="input input-sm input-bordered w-full"
						required
						bind:value={device.expand.ports[index].number}
					/>
				</div>
			</div>
			<div>
				<label class="label p-0" for="port-link-{index}">
					<span class="label-text">{$LL.device.link()}</span>
				</label>
				<input
					id="port-link-{index}"
					type="url"
					placeholder="https:// ..."
					class="input input-sm input-bordered w-full"
					bind:value={device.expand.ports[index].link}
				/>
			</div>
			<div class="card-actions justify-end">
				<button
					class="btn btn-xs btn-outline btn-error"
					on:click={() => deletePort(device.expand.ports[index])}
					type="button"><Fa icon={faTrash} />{$LL.buttons.delete()}</button
				>
			</div>
		</div>
	</div>
{/if}
