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
	<fieldset class="fieldset bg-base-100 border-base-300 rounded-box w-xs max-w-full border p-4">
		<legend class="fieldset-legend"># {device.expand.ports.length}</legend>
		<div class="flex flex-row gap-2">
			<fieldset class="fieldset">
				<label class="floating-label">
					<span>{$LL.device.ports_name()} <span class="text-error">*</span></span>
					<input
						type="text"
						placeholder={$LL.device.ports_name()}
						class="input"
						required
						bind:value={device.expand.ports[index].name}
					/>
				</label>
			</fieldset>
			<fieldset class="fieldset">
				<label class="floating-label">
					<span>{$LL.device.ports_number()} <span class="text-error">*</span></span>
					<input
						type="number"
						placeholder={$LL.device.ports_number()}
						class="input"
						min="1"
						max="65535"
						required
						bind:value={device.expand.ports[index].number}
					/>
				</label>
			</fieldset>
		</div>
		<fieldset class="fieldset">
			<label class="floating-label">
				<span>{$LL.device.link()}</span>
				<input
					type="text"
					placeholder={$LL.device.link()}
					class="input"
					bind:value={device.expand.ports[index].link}
				/>
			</label>
		</fieldset>
		<button
			class="btn btn-error btn-xs btn-soft ms-auto"
			on:click={() => deletePort(device.expand.ports[index])}
			type="button"><Fa icon={faTrash} />{$LL.buttons.delete()}</button
		>
	</fieldset>
{/if}
