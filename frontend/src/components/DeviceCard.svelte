<script>
	import { dev } from '$app/environment';
	import { parseISO, formatDistance } from 'date-fns';
	import { faPowerOff, faEllipsisVertical, faCircle } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';

	export let device;
	export let now;

	function shutdown() {
		fetch(`${dev ? 'http://127.0.0.1:8090' : ''}/api/upsnap/shutdown/${device.id}`);
	}

	function wake() {
		fetch(`${dev ? 'http://127.0.0.1:8090' : ''}/api/upsnap/wake/${device.id}`);
	}

	function sortPorts(a, b) {
		return a.number - b.number;
	}
</script>

<div class="col-xs-12 col-sm-6 col-md-4 col-lg-3 g-4">
	<div
		class="card border-0 rounded-5 px-3 py-2 shadow-sm"
		class:offline={device.status == 'offline' ? true : false}
		class:online={device.status == 'online' ? true : false}
		class:pending={device.status == 'pending' ? true : false}
	>
		<div class="card-body">
			<div class="row">
				<div class="col-auto me-auto">
					{#if device.status === 'offline'}
						<span class="text-danger" on:click={() => wake()} on:keydown={() => wake()}>
							<div role="button">
								<Fa icon={faPowerOff} class="fs-4 power-hover" />
							</div>
						</span>
					{:else if device.status === 'online'}
						{#if device.shutdown_cmd !== ''}
							<span class="text-success" on:click={() => shutdown()} on:keydown={() => shutdown()}>
								<div role="button">
									<Fa icon={faPowerOff} class="fs-4 power-hover" />
								</div>
							</span>
						{:else}
							<span class="text-success">
								<div>
									<Fa icon={faPowerOff} class="fs-4" />
								</div>
							</span>
						{/if}
					{:else if device.status === 'pending'}
						<div
							class="spinner-border text-warning"
							style="width: 23px;height: 23px;"
							role="status"
						>
							<span class="visually-hidden">Loading...</span>
						</div>
					{/if}
				</div>
				<div class="col-auto fs-5">
					<Fa icon={faEllipsisVertical} />
				</div>
			</div>
			{#if device.link}
				<p class="m-0 fw-bold fs-5">
					<a class="text-reset" target="_blank" rel="noreferrer" href={device.link}>{device.name}</a
					>
				</p>
			{:else}
				<p class="m-0 fw-bold fs-5">{device.name}</p>
			{/if}
			<p class="text-muted mb-2">{device.ip}</p>
			{#if device?.expand?.ports}
				<div class="mb-2">
					{#each device.expand.ports.sort(sortPorts) as port}
						<p class="m-0">
							<Fa icon={faCircle} class="me-1 {port.status ? 'port-up' : 'port-down'}" />{port.name}
							<span class="text-muted">({port.number})</span>
						</p>
					{/each}
				</div>
			{/if}
			<p class="text-muted m-0">
				{formatDistance(parseISO(device.updated), now, {
					includeSeconds: true,
					addSuffix: true
				})}
			</p>
		</div>
	</div>
</div>
