<script>
	import { dev } from '$app/environment';
	import { parseISO, formatDistance } from 'date-fns';
	import { faPowerOff, faEllipsisVertical } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';

	export let device;
	export let now;

	function shutdown() {
		fetch(`${dev ? 'http://127.0.0.1:8090' : ''}/api/upsnap/shutdown/${device.id}`);
	}

	function wake() {
		fetch(`${dev ? 'http://127.0.0.1:8090' : ''}/api/upsnap/wake/${device.id}`);
	}
</script>

<div class="col-xs-12 col-sm-6 col-md-4 col-lg-3 g-4">
	<div
		class="card border-0 rounded-5 px-2"
		class:offline={device.status == 'offline' ? true : false}
		class:online={device.status == 'online' ? true : false}
		class:pending={device.status == 'pending' ? true : false}
	>
		<div class="card-body">
			<div class="row fs-4">
				<div class="col-auto me-auto">
					<span
						class:text-danger={device.status == 'offline' ? true : false}
						class:text-success={device.status == 'online' ? true : false}
						class:text-warning={device.status == 'pending' ? true : false}
						on:click={device.status == 'offline'
							? () => wake()
							: device.status == 'online'
							? () => shutdown()
							: ''}
						on:keydown={device.status == 'offline'
							? () => wake()
							: device.status == 'online'
							? () => shutdown()
							: ''}
					>
						<Fa icon={faPowerOff} />
					</span>
				</div>
				<div class="col-auto">
					<Fa icon={faEllipsisVertical} />
				</div>
			</div>
			<p class="m-0 fw-bold fs-5">{device.name}</p>
			<p class="text-muted">{device.ip}</p>
			<p class="m-0 text-muted">
				{formatDistance(parseISO(device.updated), now, {
					includeSeconds: true,
					addSuffix: true
				})}
			</p>
		</div>
	</div>
</div>
