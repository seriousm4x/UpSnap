<script>
	import { parseISO, formatDistance } from 'date-fns';
	import Fa from 'svelte-fa';
	import { faPowerOff } from '@fortawesome/free-solid-svg-icons';

	export let device;
	export let now;
</script>

<div class="col-xs-12 col-sm-6 col-md-4 col-lg-3 g-4">
	<div class="card border-0 rounded-5 px-2 bg-body-secondary text-dark-emphasis">
		<div class="card-body">
			<p
				class="m-0 fs-4"
				class:text-danger={device.status == 'offline' ? true : false}
				class:text-success={device.status == 'online' ? true : false}
                class:text-warning={device.status == 'pending' ? true : false}
			>
				<Fa icon={faPowerOff} />
			</p>
			<p class="m-0 fw-bold fs-5">{device.name} ({device.status})</p>
			<p class="m-0 text-muted">{device.ip}</p>
		</div>
		<div class="card-footer mb-2 bg-transparent fst-italic">
			Last change: {formatDistance(parseISO(device.updated), now, {
				includeSeconds: true,
				addSuffix: true
			})}
		</div>
	</div>
</div>
