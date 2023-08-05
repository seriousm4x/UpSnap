<script lang="ts">
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import type { Device, Port } from '$lib/types/device';

	let device: Device = {
		expand: {
			ports: [] as Port[]
		}
	} as Device;
</script>

<h1 class="text-3xl font-bold mb-8">New device</h1>
<DeviceForm {device} />
