<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import DeviceForm from '$lib/components/DeviceForm.svelte';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import { m } from '$lib/paraglide/messages';
	import { permission, pocketbase } from '$lib/stores/pocketbase';
	import type { Device, Group, Port } from '$lib/types/device';
	import { onMount } from 'svelte';
	import toast from 'svelte-french-toast';

	let device: Device | null = $state(null);
	let error: string | null = $state(null);

	$effect(() => {
		if (Object.hasOwn($permission, 'update')) {
			const id = page.url.searchParams.get('id');
			if (!$pocketbase.authStore.isSuperuser && !$permission.update.includes(id || '')) {
				toast(m.toasts_no_permission({ url: page.url.pathname }), {
					icon: '⛔'
				});
				goto(resolve('/'));
			}
		}
	});

	onMount(async () => {
		try {
			const id = page.url.searchParams.get('id');
			if (!id) throw new Error('No device ID provided');

			const resp = await $pocketbase.collection('devices').getOne(id, { expand: 'ports,groups' });
			const dev = resp as Device;

			if (!dev.expand) {
				dev.expand = {} as { ports: Port[]; groups: Group[] };
			}
			if (!dev.expand.ports) {
				dev.expand.ports = [] as Port[];
			}

			device = dev;
		} catch (err) {
			error = String(err);
		}
	});
</script>

{#if error}
	<div class="container mx-auto max-w-lg">
		{error}
	</div>
{:else if device}
	<h1 class="mb-8 text-3xl font-bold sm:break-all">{device.name}</h1>
	<DeviceForm {device} />
{:else}
	<PageLoading />
{/if}
