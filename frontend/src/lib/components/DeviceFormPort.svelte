<script lang="ts">
	import LL from '$lib/i18n/i18n-svelte';
	import { pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Port } from '$lib/types/device';
	import { faTrash } from '@fortawesome/free-solid-svg-icons';
	import type { RecordModel } from 'pocketbase';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

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
	<fieldset class="fieldset bg-base-100 border-base-300 rounded-box w-xs border p-4">
		<legend class="fieldset-legend"># {device.expand.ports.length}</legend>
		<div class="flex flex-row gap-2">
			<div>
				<label class="fieldset-label" for="port-name-{index}">{$LL.device.ports_name()}</label>
				<input
					id="port-name-{index}"
					type="text"
					class="input input-sm"
					placeholder="ssh"
					bind:value={device.expand.ports[index].name}
				/>
			</div>
			<div>
				<label class="fieldset-label" for="port-number-{index}">{$LL.device.ports_number()}</label>
				<input
					type="number"
					class="input input-sm"
					min="1"
					max="65535"
					placeholder="22"
					required
					bind:value={device.expand.ports[index].number}
				/>
			</div>
		</div>
		<div>
			<label class="fieldset-label" for="port-link-{index}">{$LL.device.link()}</label>
			<input
				id="port-link-{index}"
				type="url"
				class="input input-sm"
				placeholder="https:// ..."
				bind:value={device.expand.ports[index].link}
			/>
		</div>
		<button
			class="btn btn-error btn-xs btn-soft ms-auto"
			on:click={() => deletePort(device.expand.ports[index])}
			type="button"><Fa icon={faTrash} />{$LL.buttons.delete()}</button
		>
	</fieldset>
{/if}
