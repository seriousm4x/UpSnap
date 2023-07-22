<script lang="ts">
	import type { Record } from 'pocketbase';

	export let device: Record;
</script>

<div class="card bg-base-300 shadow-md">
	<div class="card-body">
		<h2 class="card-title">{device.name}</h2>
		<p>{device.ip}</p>
		<div class="card-actions">
			<button class="btn btn-sm btn-success">Wake</button>
			<button class="btn btn-sm btn-error">Shutdown</button>
		</div>
	</div>
</div>
