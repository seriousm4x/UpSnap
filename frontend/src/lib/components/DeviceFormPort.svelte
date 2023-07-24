<script lang="ts">
	import { pocketbase } from '$lib/stores/pocketbase';
	import { scale } from 'svelte/transition';
	import Fa from 'svelte-fa';
	import { faTrash } from '@fortawesome/free-solid-svg-icons';
	import type { Device, Port } from '$lib/types/device';
	import type { Record } from 'pocketbase';

	export let changes: Device;
	export let index: number;
	export let portErrMsg: string;
	export let portErrTimeout: number;

	function deletePort(port: Port | Record) {
		if (port.id !== undefined) {
			$pocketbase
				.collection('ports')
				.delete(port.id)
				.then(() => {
					changes.ports = changes.ports.filter((id: number) => id !== changes.ports[index]);
				})
				.catch((err) => {
					clearTimeout(portErrTimeout);
					portErrTimeout = setTimeout(() => {
						portErrMsg = '';
					}, 10000);
					portErrMsg = err;
				});
		}
		changes.expand.ports = changes.expand.ports.filter(
			(p) => p.id !== changes.expand.ports[index].id
		);
	}
</script>

{#if changes.expand.ports[index]}
	<div class="card bg-base-100 shadow-xl" transition:scale={{ delay: 0, duration: 200 }}>
		<div class="card-body p-3">
			<div class="flex flex-row gap-2">
				<div>
					<label class="label p-0" for="port-name-{index}">
						<span class="label-text">Name</span>
					</label>
					<input
						id="port-name-{index}"
						type="text"
						placeholder="ssh"
						class="input input-sm input-bordered w-full max-w-xs"
						bind:value={changes.expand.ports[index].name}
					/>
				</div>
				<div>
					<label class="label p-0" for="port-number-{index}">
						<span class="label-text">Number</span>
					</label>
					<input
						id="port-number-{index}"
						type="number"
						placeholder="22"
						min="1"
						max="65535"
						class="input input-sm input-bordered w-24"
						bind:value={changes.expand.ports[index].number}
					/>
				</div>
			</div>
			<div class="card-actions justify-end">
				<button
					class="btn btn-xs btn-error"
					on:click={() => deletePort(changes.expand.ports[index])}
					type="button"><Fa icon={faTrash} />Delete</button
				>
			</div>
		</div>
	</div>
{/if}
