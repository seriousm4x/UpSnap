<script lang="ts">
	import type { Device } from '$lib/types/device';
	import { formatDistanceToNow, parseISO } from 'date-fns';
	import DeviceCardNic from './DeviceCardNic.svelte';

	export let device: Device;
</script>

<div class="card bg-base-300 shadow-md rounded-3xl">
	<div class="card-body p-6">
		<h1 class="card-title">{device.name}</h1>
		<ul class="menu bg-base-200 rounded-box">
			<li class="menu-title">Interfaces</li>
			<!-- TODO: change to nic array once backend supports it -->
			<DeviceCardNic {device} />
		</ul>
		<div class="card-actions">
			<span class="tooltip" data-tip="Last status change">
				{formatDistanceToNow(parseISO(device.updated), {
					includeSeconds: true,
					addSuffix: true
				})}
			</span>
		</div>
	</div>
</div>
